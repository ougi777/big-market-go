package http

import (
	"context"
	"errors"
	stdhttp "net/http"
	"strings"

	"bm-go/internal/domain/strategy/rule/chain"
	strategyservice "bm-go/internal/domain/strategy/service"
	"bm-go/internal/types"

	"github.com/gin-gonic/gin"
)

type strategyArmoryService interface {
	AssembleLotteryStrategy(ctx context.Context, strategyID int64) error
}

type raffleStrategyService interface {
	PerformRaffle(ctx context.Context, userID string, strategyID int64) (chain.AwardResult, error)
}

type strategyQueryService interface {
	QueryRaffleAwardList(ctx context.Context, activityID int64, userID string) ([]strategyservice.RaffleAward, error)
	QueryRaffleStrategyRuleWeight(ctx context.Context, activityID int64, userID string) ([]strategyservice.RaffleStrategyRuleWeight, error)
}

type raffleStrategyController struct {
	armoryService strategyArmoryService
	raffleService raffleStrategyService
	queryService  strategyQueryService
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

type raffleAwardListRequest struct {
	UserID     string `json:"userId"`
	ActivityID int64  `json:"activityId"`
}

type raffleAwardListResponse struct {
	AwardID            int    `json:"awardId"`
	AwardTitle         string `json:"awardTitle"`
	AwardSubtitle      string `json:"awardSubtitle"`
	Sort               int    `json:"sort"`
	AwardRuleLockCount *int   `json:"awardRuleLockCount"`
	IsAwardUnlock      bool   `json:"isAwardUnlock"`
	WaitUnlockCount    int    `json:"waitUnLockCount"`
}

type raffleStrategyRuleWeightRequest struct {
	UserID     string `json:"userId"`
	ActivityID int64  `json:"activityId"`
}

type raffleStrategyRuleWeightResponse struct {
	RuleWeightCount                  int                     `json:"ruleWeightCount"`
	UserActivityAccountTotalUseCount int                     `json:"userActivityAccountTotalUseCount"`
	StrategyAwards                   []strategyAwardResponse `json:"strategyAwards"`
}

type strategyAwardResponse struct {
	AwardID    int    `json:"awardId"`
	AwardTitle string `json:"awardTitle"`
}

func registerRaffleStrategyRoutes(router *gin.RouterGroup, opts RouterOptions) {
	controller := &raffleStrategyController{
		armoryService: opts.ArmoryService,
		raffleService: opts.RaffleService,
		queryService:  opts.QueryService,
	}
	if controller.armoryService == nil {
		controller.armoryService = nilArmoryService{}
	}
	if controller.raffleService == nil {
		controller.raffleService = nilRaffleService{}
	}
	if controller.queryService == nil {
		controller.queryService = nilQueryService{}
	}

	strategyGroup := router.Group("/raffle/strategy")
	strategyGroup.GET("/strategy_armory", controller.strategyArmory)
	strategyGroup.POST("/random_raffle", controller.randomRaffle)
	strategyGroup.POST("/query_raffle_award_list", controller.queryRaffleAwardList)
	strategyGroup.POST("/query_raffle_strategy_rule_weight", controller.queryRaffleStrategyRuleWeight)
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
		AwardIndex: result.AwardIndex,
	}))
}

func (c *raffleStrategyController) queryRaffleAwardList(ctx *gin.Context) {
	var request raffleAwardListRequest
	if err := ctx.ShouldBindJSON(&request); err != nil || strings.TrimSpace(request.UserID) == "" || request.ActivityID == 0 {
		ctx.JSON(stdhttp.StatusOK, types.Failure(types.ResponseCodeIllegalParam, []raffleAwardListResponse{}))
		return
	}

	awards, err := c.queryService.QueryRaffleAwardList(ctx.Request.Context(), request.ActivityID, request.UserID)
	if err != nil {
		ctx.JSON(stdhttp.StatusOK, types.Failure(types.ResponseCodeUnknownError, []raffleAwardListResponse{}))
		return
	}

	response := make([]raffleAwardListResponse, 0, len(awards))
	for _, award := range awards {
		var lockCount *int
		if award.HasAwardRuleLock {
			value := award.AwardRuleLockCount
			lockCount = &value
		}

		response = append(response, raffleAwardListResponse{
			AwardID:            award.AwardID,
			AwardTitle:         award.AwardTitle,
			AwardSubtitle:      award.AwardSubtitle,
			Sort:               award.Sort,
			AwardRuleLockCount: lockCount,
			IsAwardUnlock:      award.IsAwardUnlock,
			WaitUnlockCount:    award.WaitUnlockCount,
		})
	}

	ctx.JSON(stdhttp.StatusOK, types.Success(response))
}

func (c *raffleStrategyController) queryRaffleStrategyRuleWeight(ctx *gin.Context) {
	var request raffleStrategyRuleWeightRequest
	if err := ctx.ShouldBindJSON(&request); err != nil || strings.TrimSpace(request.UserID) == "" || request.ActivityID == 0 {
		ctx.JSON(stdhttp.StatusOK, types.Failure(types.ResponseCodeIllegalParam, []raffleStrategyRuleWeightResponse{}))
		return
	}

	ruleWeights, err := c.queryService.QueryRaffleStrategyRuleWeight(ctx.Request.Context(), request.ActivityID, request.UserID)
	if err != nil {
		ctx.JSON(stdhttp.StatusOK, types.Failure(types.ResponseCodeUnknownError, []raffleStrategyRuleWeightResponse{}))
		return
	}

	response := make([]raffleStrategyRuleWeightResponse, 0, len(ruleWeights))
	for _, ruleWeight := range ruleWeights {
		awards := make([]strategyAwardResponse, 0, len(ruleWeight.StrategyAwards))
		for _, award := range ruleWeight.StrategyAwards {
			awards = append(awards, strategyAwardResponse{
				AwardID:    award.AwardID,
				AwardTitle: award.AwardTitle,
			})
		}

		response = append(response, raffleStrategyRuleWeightResponse{
			RuleWeightCount:                  ruleWeight.RuleWeightCount,
			UserActivityAccountTotalUseCount: ruleWeight.UserActivityAccountTotalUseCount,
			StrategyAwards:                   awards,
		})
	}

	ctx.JSON(stdhttp.StatusOK, types.Success(response))
}

type nilArmoryService struct{}

func (nilArmoryService) AssembleLotteryStrategy(ctx context.Context, strategyID int64) error {
	return errors.New("armory service is not configured")
}

type nilRaffleService struct{}

func (nilRaffleService) PerformRaffle(ctx context.Context, userID string, strategyID int64) (chain.AwardResult, error) {
	return chain.AwardResult{}, errors.New("raffle service is not configured")
}

type nilQueryService struct{}

func (nilQueryService) QueryRaffleAwardList(ctx context.Context, activityID int64, userID string) ([]strategyservice.RaffleAward, error) {
	return nil, errors.New("query service is not configured")
}

func (nilQueryService) QueryRaffleStrategyRuleWeight(ctx context.Context, activityID int64, userID string) ([]strategyservice.RaffleStrategyRuleWeight, error) {
	return nil, errors.New("query service is not configured")
}
