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

type raffleActivityController struct {
	accountService    activityAccountService
	skuProductService activitySkuProductService
}

type userActivityAccountRequest struct {
	UserID     string `json:"userId"`
	ActivityID int64  `json:"activityId"`
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
		accountService:    opts.ActivityAccountService,
		skuProductService: opts.ActivitySkuProductService,
	}
	if controller.accountService == nil {
		controller.accountService = nilActivityAccountService{}
	}
	if controller.skuProductService == nil {
		controller.skuProductService = nilActivitySkuProductService{}
	}

	activityGroup := router.Group("/raffle/activity")
	activityGroup.GET("/query_sku_product_list_by_activity_id", controller.querySkuProductListByActivityID)
	activityGroup.POST("/query_user_activity_account", controller.queryUserActivityAccount)
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

type nilActivityAccountService struct{}

func (nilActivityAccountService) QueryActivityAccount(ctx context.Context, activityID int64, userID string) (activity.AccountEntity, error) {
	return activity.AccountEntity{}, errors.New("activity account service is not configured")
}

type nilActivitySkuProductService struct{}

func (nilActivitySkuProductService) QuerySkuProductListByActivityID(ctx context.Context, activityID int64) ([]activity.SkuProductEntity, error) {
	return nil, errors.New("activity sku product service is not configured")
}
