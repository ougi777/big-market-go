package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bm-go/internal/config"
	activityservice "bm-go/internal/domain/activity/service"
	awardservice "bm-go/internal/domain/award/service"
	creditservice "bm-go/internal/domain/credit/service"
	rebateservice "bm-go/internal/domain/rebate/service"
	"bm-go/internal/domain/strategy/rule/chain"
	"bm-go/internal/domain/strategy/rule/tree"
	strategyservice "bm-go/internal/domain/strategy/service"
	taskservice "bm-go/internal/domain/task/service"
	"bm-go/internal/infrastructure/persistent/mysql"
	"bm-go/internal/infrastructure/persistent/repository"
	"bm-go/internal/infrastructure/persistent/sharding"
	infrabbitmq "bm-go/internal/infrastructure/rabbitmq"
	infrredis "bm-go/internal/infrastructure/redis"
	triggerhttp "bm-go/internal/trigger/http"
	triggerjob "bm-go/internal/trigger/job"
	triggerlistener "bm-go/internal/trigger/listener"

	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer func() { _ = logger.Sync() }()

	dbRouter, err := mysql.OpenRouter(cfg.MySQL)
	if err != nil {
		logger.Fatal("open mysql failed", zap.Error(err))
	}
	redisClient := infrredis.NewClient(cfg.Redis)
	rabbitmqClient, err := infrabbitmq.Dial(cfg.RabbitMQ)
	if err != nil {
		logger.Fatal("open rabbitmq failed", zap.Error(err))
	}
	defer func() { _ = rabbitmqClient.Close() }()

	strategyStore := infrredis.NewStrategyStore(redisClient)
	activityStore := infrredis.NewActivityStore(redisClient)
	tableRouter := sharding.NewRouterWithDBCount(cfg.Sharding.DBCount, cfg.Sharding.TableCount)
	strategyRepository := repository.NewStrategyRepositoryWithDBRouter(dbRouter, tableRouter, strategyStore)
	activityRepository := repository.NewActivityRepositoryWithDBRouter(dbRouter, tableRouter)
	creditRepository := repository.NewCreditRepositoryWithDBRouter(dbRouter, tableRouter)
	awardRepository := repository.NewAwardRepositoryWithDBRouter(dbRouter, tableRouter)
	rebateRepository := repository.NewRebateRepositoryWithDBRouter(dbRouter, tableRouter)
	strategyDispatch := repository.NewStrategyDispatch(redisClient)
	chainFactory := chain.NewFactory(strategyRepository, strategyDispatch)
	treeNodes := map[string]tree.Node{
		tree.RuleLock:      tree.NewLockNode(strategyRepository),
		tree.RuleStock:     tree.NewStockNode(strategyRepository, strategyDispatch),
		tree.RuleLuckAward: tree.NewLuckAwardNode(strategyRepository),
	}
	armoryService := strategyservice.NewArmoryService(strategyRepository, strategyStore)
	raffleService := strategyservice.NewRaffleService(chainFactory, strategyRepository, treeNodes)
	queryService := strategyservice.NewQueryService(strategyRepository)
	stockService := strategyservice.NewStockService(strategyRepository, strategyStore)
	activityAccountService := activityservice.NewAccountService(activityRepository)
	activitySkuProductService := activityservice.NewSkuProductService(activityRepository)
	activityCreditService := creditservice.NewAccountService(creditRepository)
	activityArmoryService := activityservice.NewArmoryService(activityRepository, activityStore)
	activityPartakeService := activityservice.NewPartakeService(activityRepository)
	activityStockService := activityservice.NewStockService(activityRepository, activityStore, activityStore, rabbitmqClient)
	activityExchangeService := activityservice.NewExchangeService(activityRepository, activityStockService, rabbitmqClient, creditRepository)
	activityRebateProcessor := activityservice.NewRebateProcessor(activityRepository, creditRepository)
	activityDeliveryService := activityservice.NewDeliveryService(activityRepository)
	awardService := awardservice.NewAwardService(awardRepository, awardRepository, rabbitmqClient)
	taskService := taskservice.NewService(awardRepository, rabbitmqClient)
	rebateService := rebateservice.NewRebateService(rebateRepository, rabbitmqClient)
	activityDrawService := activityservice.NewDrawService(activityPartakeService, raffleService, awardService)

	router := triggerhttp.NewRouter(triggerhttp.RouterOptions{
		Logger:                        logger,
		ArmoryService:                 armoryService,
		RaffleService:                 raffleService,
		QueryService:                  queryService,
		ActivityAccountService:        activityAccountService,
		ActivitySkuProductService:     activitySkuProductService,
		ActivityArmoryService:         activityArmoryService,
		ActivityStrategyArmoryService: armoryService,
		ActivityDrawService:           activityDrawService,
		ActivityExchangeService:       activityExchangeService,
		ActivityCreditService:         activityCreditService,
		ActivityRebateService:         rebateService,
	})
	server := &http.Server{
		Addr:              cfg.HTTPAddr(),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	scheduler := triggerjob.NewScheduler()
	updateAwardStockJob := triggerjob.NewUpdateAwardStockJob(stockService, logger)
	if _, err := scheduler.Add("*/5 * * * * *", updateAwardStockJob.Exec); err != nil {
		logger.Fatal("register update award stock job failed", zap.Error(err))
	}
	updateActivitySkuStockJob := triggerjob.NewUpdateActivitySkuStockJob(activityStockService, logger)
	if _, err := scheduler.Add("*/5 * * * * *", updateActivitySkuStockJob.Exec); err != nil {
		logger.Fatal("register update activity sku stock job failed", zap.Error(err))
	}
	sendMessageTaskJob := triggerjob.NewSendMessageTaskJob(taskService, logger)
	if _, err := scheduler.Add("*/5 * * * * *", sendMessageTaskJob.Exec); err != nil {
		logger.Fatal("register send message task job failed", zap.Error(err))
	}
	scheduler.Start()

	sendAwardConsumer := triggerlistener.NewSendAwardConsumer(rabbitmqClient, awardService, logger)
	consumerCtx, stopConsumer := context.WithCancel(context.Background())
	defer stopConsumer()
	if err := sendAwardConsumer.Start(consumerCtx); err != nil {
		logger.Fatal("start send award consumer failed", zap.Error(err))
	}
	activitySkuStockZeroConsumer := triggerlistener.NewActivitySkuStockZeroConsumer(rabbitmqClient, activityStockService, logger)
	if err := activitySkuStockZeroConsumer.Start(consumerCtx); err != nil {
		logger.Fatal("start activity sku stock zero consumer failed", zap.Error(err))
	}
	sendRebateConsumer := triggerlistener.NewSendRebateConsumer(rabbitmqClient, activityRebateProcessor, logger)
	if err := sendRebateConsumer.Start(consumerCtx); err != nil {
		logger.Fatal("start send rebate consumer failed", zap.Error(err))
	}
	creditAdjustSuccessConsumer := triggerlistener.NewCreditAdjustSuccessConsumer(rabbitmqClient, activityDeliveryService, logger)
	if err := creditAdjustSuccessConsumer.Start(consumerCtx); err != nil {
		logger.Fatal("start credit adjust success consumer failed", zap.Error(err))
	}

	go func() {
		logger.Info("big-market go service started", zap.String("addr", cfg.HTTPAddr()))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("http server stopped", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	_ = sendAwardConsumer.Stop(context.Background())
	_ = activitySkuStockZeroConsumer.Stop(context.Background())
	_ = sendRebateConsumer.Stop(context.Background())
	_ = creditAdjustSuccessConsumer.Stop(context.Background())
	scheduler.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("http server shutdown failed", zap.Error(err))
	}
	logger.Info("big-market go service stopped")
}
