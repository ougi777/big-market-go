package app

import (
	"context"
	"errors"
	"net/http"
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

type Application struct {
	cfg            *config.Config
	logger         *zap.Logger
	server         *http.Server
	scheduler      *triggerjob.Scheduler
	rabbitmqClient *infrabbitmq.Client
	consumers      []triggerlistener.Consumer
	consumerCancel context.CancelFunc
}

type components struct {
	rabbitmqClient          *infrabbitmq.Client
	armoryService           *strategyservice.ArmoryService
	raffleService           *strategyservice.RaffleService
	queryService            *strategyservice.QueryService
	stockService            *strategyservice.StockService
	activityAccountService  *activityservice.AccountService
	activitySkuProduct      *activityservice.SkuProductService
	activityCreditService   *creditservice.AccountService
	activityArmoryService   *activityservice.ArmoryService
	activityStockService    *activityservice.StockService
	activityExchangeService *activityservice.ExchangeService
	activityRebateProcessor *activityservice.RebateProcessor
	activityDeliveryService *activityservice.DeliveryService
	awardService            *awardservice.AwardService
	taskService             *taskservice.Service
	rebateService           *rebateservice.RebateService
	activityDrawService     *activityservice.DrawService
}

func New(cfg *config.Config, logger *zap.Logger) (*Application, error) {
	if logger == nil {
		logger = zap.NewNop()
	}

	components, err := newComponents(cfg, logger)
	if err != nil {
		return nil, err
	}

	router := triggerhttp.NewRouter(triggerhttp.RouterOptions{
		Logger:                        logger,
		ArmoryService:                 components.armoryService,
		RaffleService:                 components.raffleService,
		QueryService:                  components.queryService,
		ActivityAccountService:        components.activityAccountService,
		ActivitySkuProductService:     components.activitySkuProduct,
		ActivityArmoryService:         components.activityArmoryService,
		ActivityStrategyArmoryService: components.armoryService,
		ActivityDrawService:           components.activityDrawService,
		ActivityExchangeService:       components.activityExchangeService,
		ActivityCreditService:         components.activityCreditService,
		ActivityRebateService:         components.rebateService,
	})

	scheduler, err := newScheduler(cfg.JobSpec(), logger, components.stockService, components.activityStockService, components.taskService)
	if err != nil {
		return nil, err
	}

	server := &http.Server{
		Addr:              cfg.HTTPAddr(),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	consumers := newConsumers(
		components.rabbitmqClient,
		components.awardService,
		components.activityStockService,
		components.activityRebateProcessor,
		components.activityDeliveryService,
		logger,
	)

	return &Application{
		cfg:            cfg,
		logger:         logger,
		server:         server,
		scheduler:      scheduler,
		rabbitmqClient: components.rabbitmqClient,
		consumers:      consumers,
	}, nil
}

func newComponents(cfg *config.Config, logger *zap.Logger) (*components, error) {
	dbRouter, err := mysql.OpenRouter(cfg.MySQL)
	if err != nil {
		return nil, err
	}
	redisClient := infrredis.NewClient(cfg.Redis)
	rabbitmqClient, err := infrabbitmq.Dial(cfg.RabbitMQ)
	if err != nil {
		return nil, err
	}

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

	return &components{
		rabbitmqClient:          rabbitmqClient,
		armoryService:           armoryService,
		raffleService:           raffleService,
		queryService:            queryService,
		stockService:            stockService,
		activityAccountService:  activityAccountService,
		activitySkuProduct:      activitySkuProductService,
		activityCreditService:   activityCreditService,
		activityArmoryService:   activityArmoryService,
		activityStockService:    activityStockService,
		activityExchangeService: activityExchangeService,
		activityRebateProcessor: activityRebateProcessor,
		activityDeliveryService: activityDeliveryService,
		awardService:            awardService,
		taskService:             taskService,
		rebateService:           rebateService,
		activityDrawService:     activityDrawService,
	}, nil
}

func newConsumers(
	rabbitmqClient *infrabbitmq.Client,
	awardService *awardservice.AwardService,
	activityStockService *activityservice.StockService,
	activityRebateProcessor *activityservice.RebateProcessor,
	activityDeliveryService *activityservice.DeliveryService,
	logger *zap.Logger,
) []triggerlistener.Consumer {
	return []triggerlistener.Consumer{
		triggerlistener.NewSendAwardConsumer(rabbitmqClient, awardService, logger),
		triggerlistener.NewActivitySkuStockZeroConsumer(rabbitmqClient, activityStockService, logger),
		triggerlistener.NewSendRebateConsumer(rabbitmqClient, activityRebateProcessor, logger),
		triggerlistener.NewCreditAdjustSuccessConsumer(rabbitmqClient, activityDeliveryService, logger),
	}
}

func newScheduler(
	jobSpec string,
	logger *zap.Logger,
	stockService *strategyservice.StockService,
	activityStockService *activityservice.StockService,
	taskService *taskservice.Service,
) (*triggerjob.Scheduler, error) {
	scheduler := triggerjob.NewScheduler()
	if _, err := scheduler.Add(jobSpec, triggerjob.NewUpdateAwardStockJob(stockService, logger).Exec); err != nil {
		return nil, err
	}
	if _, err := scheduler.Add(jobSpec, triggerjob.NewUpdateActivitySkuStockJob(activityStockService, logger).Exec); err != nil {
		return nil, err
	}
	if _, err := scheduler.Add(jobSpec, triggerjob.NewSendMessageTaskJob(taskService, logger).Exec); err != nil {
		return nil, err
	}
	return scheduler, nil
}

func (a *Application) Start(ctx context.Context) error {
	a.scheduler.Start()

	consumerCtx, cancel := context.WithCancel(ctx)
	a.consumerCancel = cancel
	for _, consumer := range a.consumers {
		if err := consumer.Start(consumerCtx); err != nil {
			cancel()
			a.scheduler.Stop()
			return err
		}
	}

	go func() {
		a.logger.Info("big-market go service started", zap.String("addr", a.cfg.HTTPAddr()))
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Fatal("http server stopped", zap.Error(err))
		}
	}()
	return nil
}

func (a *Application) Shutdown(ctx context.Context) error {
	if a.consumerCancel != nil {
		a.consumerCancel()
	}
	for _, consumer := range a.consumers {
		if err := consumer.Stop(ctx); err != nil {
			a.logger.Error("stop consumer failed", zap.Error(err))
		}
	}
	if a.scheduler != nil {
		a.scheduler.Stop()
	}

	var shutdownErr error
	if a.server != nil {
		shutdownErr = a.server.Shutdown(ctx)
	}
	if err := a.rabbitmqClient.Close(); err != nil && shutdownErr == nil {
		shutdownErr = err
	}
	if shutdownErr == nil {
		a.logger.Info("big-market go service stopped")
	}
	return shutdownErr
}
