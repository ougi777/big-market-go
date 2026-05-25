package http

import (
	"context"
	"errors"
	stdhttp "net/http"

	"bm-go/internal/domain/strategy/rule/chain"
	"bm-go/internal/types"

	"github.com/gin-gonic/gin"
)

type strategyArmoryService interface {
	AssembleLotteryStrategy(ctx context.Context, strategyID int64) error
}

type raffleStrategyService interface {
	PerformRaffle(ctx context.Context, userID string, strategyID int64) (chain.AwardResult, error)
}

type raffleStrategyController struct {
	armoryService strategyArmoryService
	raffleService raffleStrategyService
}

type strategyArmoryRequest struct {
	StrategyID int64 `form:"strategyId"`
}

type raffleStrategyRequest struct {
	StrategyID int64 `json:"strategyId"`
}

type raffleStrategyResponse struct {
	AwardID    int `json:"awardId"`
	AwardIndex int `json:"awardIndex"`
}

func registerRaffleStrategyRoutes(router *gin.RouterGroup, opts RouterOptions) {
	controller := &raffleStrategyController{
		armoryService: opts.ArmoryService,
		raffleService: opts.RaffleService,
	}
	if controller.armoryService == nil {
		controller.armoryService = nilArmoryService{}
	}
	if controller.raffleService == nil {
		controller.raffleService = nilRaffleService{}
	}

	strategyGroup := router.Group("/raffle/strategy")
	strategyGroup.GET("/strategy_armory", controller.strategyArmory)
	strategyGroup.POST("/random_raffle", controller.randomRaffle)
}

func (c *raffleStrategyController) strategyArmory(ctx *gin.Context) {
	var request strategyArmoryRequest
	if err := ctx.ShouldBindQuery(&request); err != nil || request.StrategyID == 0 {
		ctx.JSON(stdhttp.StatusOK, types.Failure(types.ResponseCodeIllegalParam, false))
		return
	}

	if err := c.armoryService.AssembleLotteryStrategy(ctx.Request.Context(), request.StrategyID); err != nil {
		ctx.JSON(stdhttp.StatusOK, types.Failure(types.ResponseCodeUnknownError, false))
		return
	}

	ctx.JSON(stdhttp.StatusOK, types.Success(true))
}

func (c *raffleStrategyController) randomRaffle(ctx *gin.Context) {
	var request raffleStrategyRequest
	if err := ctx.ShouldBindJSON(&request); err != nil || request.StrategyID == 0 {
		ctx.JSON(stdhttp.StatusOK, types.Failure(types.ResponseCodeIllegalParam, raffleStrategyResponse{}))
		return
	}

	result, err := c.raffleService.PerformRaffle(ctx.Request.Context(), "system", request.StrategyID)
	if err != nil {
		ctx.JSON(stdhttp.StatusOK, types.Failure(types.ResponseCodeUnknownError, raffleStrategyResponse{}))
		return
	}

	ctx.JSON(stdhttp.StatusOK, types.Success(raffleStrategyResponse{
		AwardID:    result.AwardID,
		AwardIndex: 0,
	}))
}

type nilArmoryService struct{}

func (nilArmoryService) AssembleLotteryStrategy(ctx context.Context, strategyID int64) error {
	return errors.New("armory service is not configured")
}

type nilRaffleService struct{}

func (nilRaffleService) PerformRaffle(ctx context.Context, userID string, strategyID int64) (chain.AwardResult, error) {
	return chain.AwardResult{}, errors.New("raffle service is not configured")
}
