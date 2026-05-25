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

type raffleActivityController struct {
	accountService activityAccountService
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

func registerRaffleActivityRoutes(router *gin.RouterGroup, opts RouterOptions) {
	controller := &raffleActivityController{
		accountService: opts.ActivityAccountService,
	}
	if controller.accountService == nil {
		controller.accountService = nilActivityAccountService{}
	}

	activityGroup := router.Group("/raffle/activity")
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

type nilActivityAccountService struct{}

func (nilActivityAccountService) QueryActivityAccount(ctx context.Context, activityID int64, userID string) (activity.AccountEntity, error) {
	return activity.AccountEntity{}, errors.New("activity account service is not configured")
}
