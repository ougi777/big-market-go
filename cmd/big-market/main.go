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
	"bm-go/internal/domain/strategy/rule/chain"
	"bm-go/internal/domain/strategy/rule/tree"
	strategyservice "bm-go/internal/domain/strategy/service"
	"bm-go/internal/infrastructure/persistent/mysql"
	"bm-go/internal/infrastructure/persistent/repository"
	infrredis "bm-go/internal/infrastructure/redis"
	triggerhttp "bm-go/internal/trigger/http"
	triggerjob "bm-go/internal/trigger/job"

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

	db, err := mysql.Open(cfg.MySQL)
	if err != nil {
		logger.Fatal("open mysql failed", zap.Error(err))
	}
	redisClient := infrredis.NewClient(cfg.Redis)

	strategyStore := infrredis.NewStrategyStore(redisClient)
	strategyRepository := repository.NewStrategyRepository(db, strategyStore)
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

	router := triggerhttp.NewRouter(triggerhttp.RouterOptions{
		Logger:        logger,
		ArmoryService: armoryService,
		RaffleService: raffleService,
		QueryService:  queryService,
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
	scheduler.Start()

	go func() {
		logger.Info("big-market go service started", zap.String("addr", cfg.HTTPAddr()))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("http server stopped", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	scheduler.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("http server shutdown failed", zap.Error(err))
	}
	logger.Info("big-market go service stopped")
}
