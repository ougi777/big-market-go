package http

import (
	"context"
	"errors"
	stdhttp "net/http"
	"strings"

	"bm-go/internal/domain/activity"
	"bm-go/internal/types"

	"github.com/gin-gonic/gin"
)

type activityAccountService interface {
	QueryActivityAccount(ctx context.Context, activityID int64, userID string) (activity.AccountEntity, error)
}

type activitySkuProductService interface {
	QuerySkuProductListByActivityID(ctx context.Context, activityID int64) ([]activity.SkuProductEntity, error)
}

type activityArmoryService interface {
	AssembleActivitySkuByActivityID(ctx context.Context, activityID int64) error
}

type activityStrategyArmoryService interface {
	AssembleLotteryStrategyByActivityID(ctx context.Context, activityID int64) error
}

type activityDrawService interface {
	Draw(ctx context.Context, userID string, activityID int64) (activity.DrawResult, error)
}

type activityExchangeService interface {
	CreditPayExchangeSku(ctx context.Context, userID string, sku int64) (bool, error)
}

type activityCreditService interface {
	QueryUserCreditAccount(ctx context.Context, userID string) (float64, error)
}

type raffleActivityController struct {
	accountService        activityAccountService
	skuProductService     activitySkuProductService
	armoryService         activityArmoryService
	strategyArmoryService activityStrategyArmoryService
	drawService           activityDrawService
	exchangeService       activityExchangeService
	creditService         activityCreditService
}

type activityDrawRequest struct {
	UserID     string `json:"userId"`
	ActivityID int64  `json:"activityId"`
}

type activityDrawResponse struct {
	AwardID    int    `json:"awardId"`
	AwardTitle string `json:"awardTitle"`
	AwardIndex int    `json:"awardIndex"`
}

type userActivityAccountRequest struct {
	UserID     string `json:"userId"`
	ActivityID int64  `json:"activityId"`
}

type skuProductShopCartRequest struct {
	UserID string `json:"userId"`
	SKU    int64  `json:"sku"`
}

type userCreditAccountRequest struct {
	UserID string `form:"userId"`
}

type userActivityAccountResponse struct {
	TotalCount        int `json:"totalCount"`
	TotalCountSurplus int `json:"totalCountSurplus"`
	DayCount          int `json:"dayCount"`
	DayCountSurplus   int `json:"dayCountSurplus"`
	MonthCount        int `json:"monthCount"`
	MonthCountSurplus int `json:"monthCountSurplus"`
}

type skuProductListRequest struct {
	ActivityID int64 `form:"activityId"`
}

type activityArmoryRequest struct {
	ActivityID int64 `form:"activityId"`
}

type skuProductResponse struct {
	SKU               int64                 `json:"sku"`
	ActivityID        int64                 `json:"activityId"`
	ActivityCountID   int64                 `json:"activityCountId"`
	StockCount        int                   `json:"stockCount"`
	StockCountSurplus int                   `json:"stockCountSurplus"`
	ProductAmount     float64               `json:"productAmount"`
	ActivityCount     activityCountResponse `json:"activityCount"`
}

type activityCountResponse struct {
	TotalCount int `json:"totalCount"`
	DayCount   int `json:"dayCount"`
	MonthCount int `json:"monthCount"`
}

func registerRaffleActivityRoutes(router *gin.RouterGroup, opts RouterOptions) {
	controller := &raffleActivityController{
		accountService:        opts.ActivityAccountService,
		skuProductService:     opts.ActivitySkuProductService,
		armoryService:         opts.ActivityArmoryService,
		strategyArmoryService: opts.ActivityStrategyArmoryService,
		drawService:           opts.ActivityDrawService,
		exchangeService:       opts.ActivityExchangeService,
		creditService:         opts.ActivityCreditService,
	}
	if controller.accountService == nil {
		controller.accountService = nilActivityAccountService{}
	}
	if controller.skuProductService == nil {
		controller.skuProductService = nilActivitySkuProductService{}
	}
	if controller.armoryService == nil {
		controller.armoryService = nilActivityArmoryService{}
	}
	if controller.strategyArmoryService == nil {
		controller.strategyArmoryService = nilActivityStrategyArmoryService{}
	}
	if controller.drawService == nil {
		controller.drawService = nilActivityDrawService{}
	}
	if controller.exchangeService == nil {
		controller.exchangeService = nilActivityExchangeService{}
	}
	if controller.creditService == nil {
		controller.creditService = nilActivityCreditService{}
	}

	activityGroup := router.Group("/raffle/activity")
	activityGroup.GET("/armory", controller.armory)
	activityGroup.POST("/draw", controller.draw)
	activityGroup.GET("/query_sku_product_list_by_activity_id", controller.querySkuProductListByActivityID)
	activityGroup.GET("/query_user_credit_account", controller.queryUserCreditAccount)
	activityGroup.POST("/query_user_activity_account", controller.queryUserActivityAccount)
	activityGroup.POST("/credit_pay_exchange_sku", controller.creditPayExchangeSku)
}

func (c *raffleActivityController) armory(ctx *gin.Context) {
	var request activityArmoryRequest
	if err := ctx.ShouldBindQuery(&request); err != nil || request.ActivityID == 0 {
		ctx.JSON(stdhttp.StatusOK, types.Failure(types.ResponseCodeIllegalParam, false))
		return
	}

	if err := c.armoryService.AssembleActivitySkuByActivityID(ctx.Request.Context(), request.ActivityID); err != nil {
		ctx.JSON(stdhttp.StatusOK, types.Failure(types.ResponseCodeUnknownError, false))
		return
	}

	if err := c.strategyArmoryService.AssembleLotteryStrategyByActivityID(ctx.Request.Context(), request.ActivityID); err != nil {
		ctx.JSON(stdhttp.StatusOK, types.Failure(types.ResponseCodeUnknownError, false))
		return
	}

	ctx.JSON(stdhttp.StatusOK, types.Success(true))
}

func (c *raffleActivityController) draw(ctx *gin.Context) {
	var request activityDrawRequest
	if err := ctx.ShouldBindJSON(&request); err != nil || strings.TrimSpace(request.UserID) == "" || request.ActivityID == 0 {
		ctx.JSON(stdhttp.StatusOK, types.Failure(types.ResponseCodeIllegalParam, activityDrawResponse{}))
		return
	}

	result, err := c.drawService.Draw(ctx.Request.Context(), request.UserID, request.ActivityID)
	if err != nil {
		var appErr types.AppError
		if errors.As(err, &appErr) {
			ctx.JSON(stdhttp.StatusOK, types.Failure(appErr.Code, activityDrawResponse{}))
			return
		}
		ctx.JSON(stdhttp.StatusOK, types.Failure(types.ResponseCodeUnknownError, activityDrawResponse{}))
		return
	}

	ctx.JSON(stdhttp.StatusOK, types.Success(activityDrawResponse{
		AwardID:    result.AwardID,
		AwardTitle: result.AwardTitle,
		AwardIndex: result.AwardIndex,
	}))
}

func (c *raffleActivityController) queryUserActivityAccount(ctx *gin.Context) {
	var request userActivityAccountRequest
	if err := ctx.ShouldBindJSON(&request); err != nil || strings.TrimSpace(request.UserID) == "" || request.ActivityID == 0 {
		ctx.JSON(stdhttp.StatusOK, types.Failure(types.ResponseCodeIllegalParam, userActivityAccountResponse{}))
		return
	}

	account, err := c.accountService.QueryActivityAccount(ctx.Request.Context(), request.ActivityID, request.UserID)
	if err != nil {
		ctx.JSON(stdhttp.StatusOK, types.Failure(types.ResponseCodeUnknownError, userActivityAccountResponse{}))
		return
	}

	ctx.JSON(stdhttp.StatusOK, types.Success(userActivityAccountResponse{
		TotalCount:        account.TotalCount,
		TotalCountSurplus: account.TotalCountSurplus,
		DayCount:          account.DayCount,
		DayCountSurplus:   account.DayCountSurplus,
		MonthCount:        account.MonthCount,
		MonthCountSurplus: account.MonthCountSurplus,
	}))
}

func (c *raffleActivityController) querySkuProductListByActivityID(ctx *gin.Context) {
	var request skuProductListRequest
	if err := ctx.ShouldBindQuery(&request); err != nil || request.ActivityID == 0 {
		ctx.JSON(stdhttp.StatusOK, types.Failure(types.ResponseCodeIllegalParam, []skuProductResponse{}))
		return
	}

	products, err := c.skuProductService.QuerySkuProductListByActivityID(ctx.Request.Context(), request.ActivityID)
	if err != nil {
		ctx.JSON(stdhttp.StatusOK, types.Failure(types.ResponseCodeUnknownError, []skuProductResponse{}))
		return
	}

	response := make([]skuProductResponse, 0, len(products))
	for _, product := range products {
		response = append(response, skuProductResponse{
			SKU:               product.SKU,
			ActivityID:        product.ActivityID,
			ActivityCountID:   product.ActivityCountID,
			StockCount:        product.StockCount,
			StockCountSurplus: product.StockCountSurplus,
			ProductAmount:     product.ProductAmount,
			ActivityCount: activityCountResponse{
				TotalCount: product.ActivityCount.TotalCount,
				DayCount:   product.ActivityCount.DayCount,
				MonthCount: product.ActivityCount.MonthCount,
			},
		})
	}

	ctx.JSON(stdhttp.StatusOK, types.Success(response))
}

func (c *raffleActivityController) creditPayExchangeSku(ctx *gin.Context) {
	var request skuProductShopCartRequest
	if err := ctx.ShouldBindJSON(&request); err != nil || strings.TrimSpace(request.UserID) == "" || request.SKU == 0 {
		ctx.JSON(stdhttp.StatusOK, types.Failure(types.ResponseCodeIllegalParam, false))
		return
	}

	result, err := c.exchangeService.CreditPayExchangeSku(ctx.Request.Context(), request.UserID, request.SKU)
	if err != nil {
		var appErr types.AppError
		if errors.As(err, &appErr) {
			ctx.JSON(stdhttp.StatusOK, types.Failure(appErr.Code, false))
			return
		}
		ctx.JSON(stdhttp.StatusOK, types.Failure(types.ResponseCodeUnknownError, false))
		return
	}
	ctx.JSON(stdhttp.StatusOK, types.Success(result))
}

func (c *raffleActivityController) queryUserCreditAccount(ctx *gin.Context) {
	var request userCreditAccountRequest
	if err := ctx.ShouldBindQuery(&request); err != nil || strings.TrimSpace(request.UserID) == "" {
		ctx.JSON(stdhttp.StatusOK, types.Failure(types.ResponseCodeIllegalParam, float64(0)))
		return
	}
	amount, err := c.creditService.QueryUserCreditAccount(ctx.Request.Context(), request.UserID)
	if err != nil {
		var appErr types.AppError
		if errors.As(err, &appErr) {
			ctx.JSON(stdhttp.StatusOK, types.Failure(appErr.Code, float64(0)))
			return
		}
		ctx.JSON(stdhttp.StatusOK, types.Failure(types.ResponseCodeUnknownError, float64(0)))
		return
	}
	ctx.JSON(stdhttp.StatusOK, types.Success(amount))
}

type nilActivityAccountService struct{}

func (nilActivityAccountService) QueryActivityAccount(ctx context.Context, activityID int64, userID string) (activity.AccountEntity, error) {
	return activity.AccountEntity{}, errors.New("activity account service is not configured")
}

type nilActivitySkuProductService struct{}

func (nilActivitySkuProductService) QuerySkuProductListByActivityID(ctx context.Context, activityID int64) ([]activity.SkuProductEntity, error) {
	return nil, errors.New("activity sku product service is not configured")
}

type nilActivityArmoryService struct{}

func (nilActivityArmoryService) AssembleActivitySkuByActivityID(ctx context.Context, activityID int64) error {
	return errors.New("activity armory service is not configured")
}

type nilActivityStrategyArmoryService struct{}

func (nilActivityStrategyArmoryService) AssembleLotteryStrategyByActivityID(ctx context.Context, activityID int64) error {
	return errors.New("activity strategy armory service is not configured")
}

type nilActivityDrawService struct{}

func (nilActivityDrawService) Draw(ctx context.Context, userID string, activityID int64) (activity.DrawResult, error) {
	return activity.DrawResult{}, errors.New("activity draw service is not configured")
}

type nilActivityExchangeService struct{}

func (nilActivityExchangeService) CreditPayExchangeSku(ctx context.Context, userID string, sku int64) (bool, error) {
	return false, errors.New("activity exchange service is not configured")
}

type nilActivityCreditService struct{}

func (nilActivityCreditService) QueryUserCreditAccount(ctx context.Context, userID string) (float64, error) {
	return 0, errors.New("activity credit service is not configured")
}
