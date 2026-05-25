package http

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RouterOptions struct {
	Logger        *zap.Logger
	ArmoryService strategyArmoryService
	RaffleService raffleStrategyService
}

func NewRouter(opts RouterOptions) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(requestLogger(opts.Logger))

	registerHealthRoutes(router)
	registerV1Routes(router, opts)

	return router
}

func registerV1Routes(router *gin.Engine, opts RouterOptions) {
	v1 := router.Group("/api/v1")
	v1.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"code": "0000", "message": "success", "data": "pong"})
	})

	registerRaffleStrategyRoutes(v1, opts)
}
