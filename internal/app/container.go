package app

import (
	"bm-go/internal/config"

	"go.uber.org/zap"
)

type Container struct {
	Config *config.Config
	Logger *zap.Logger
}

func NewContainer(cfg *config.Config, logger *zap.Logger) *Container {
	return &Container{
		Config: cfg,
		Logger: logger,
	}
}
