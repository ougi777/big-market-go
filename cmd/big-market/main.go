package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bm-go/internal/app"
	"bm-go/internal/config"

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

	application, err := app.New(cfg, logger)
	if err != nil {
		logger.Fatal("create application failed", zap.Error(err))
	}

	runCtx, stop := context.WithCancel(context.Background())
	defer stop()
	if err := application.Start(runCtx); err != nil {
		logger.Fatal("start application failed", zap.Error(err))
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := application.Shutdown(shutdownCtx); err != nil {
		logger.Fatal("shutdown application failed", zap.Error(err))
	}
}
