package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"bm-go/internal/domain/activity"
	"bm-go/internal/types"
)

func TestQueryUserActivityAccountRoute(t *testing.T) {
	router := NewRouter(RouterOptions{
		ActivityAccountService: &fakeActivityAccountService{
			account: activity.AccountEntity{
				TotalCount:        100,
				TotalCountSurplus: 80,
				DayCount:          5,
				DayCountSurplus:   3,
				MonthCount:        50,
				MonthCountSurplus: 35,
			},
		},
	})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/raffle/activity/query_user_activity_account", strings.NewReader(`{"userId":"xiaofuge","activityId":100301}`))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), `"totalCount":100`) {
		t.Fatalf("expected total count, got %s", recorder.Body.String())
	}
	if !strings.Contains(recorder.Body.String(), `"monthCountSurplus":35`) {
		t.Fatalf("expected month count surplus, got %s", recorder.Body.String())
	}
}

func TestQueryUserActivityAccountRouteReturnsAppErrorCode(t *testing.T) {
	router := NewRouter(RouterOptions{
		ActivityAccountService: &fakeActivityAccountService{
			err: types.NewAppError(types.ResponseCodeAccountQuotaError, nil),
		},
	})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/raffle/activity/query_user_activity_account", strings.NewReader(`{"userId":"xiaofuge","activityId":100301}`))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), `"code":"ERR_BIZ_006"`) {
		t.Fatalf("expected app error code, got %s", recorder.Body.String())
	}
}

func TestQueryUserActivityAccountRouteIllegalParam(t *testing.T) {
	router := NewRouter(RouterOptions{
		ActivityAccountService: &fakeActivityAccountService{},
	})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/raffle/activity/query_user_activity_account", strings.NewReader(`{"userId":"","activityId":100301}`))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), `"code":"0002"`) {
		t.Fatalf("expected illegal param code, got %s", recorder.Body.String())
	}
}

func TestActivityArmoryRoute(t *testing.T) {
	activityArmory := &fakeActivityArmoryService{}
	strategyArmory := &fakeActivityStrategyArmoryService{}
	router := NewRouter(RouterOptions{
		ActivityArmoryService:         activityArmory,
		ActivityStrategyArmoryService: strategyArmory,
	})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/raffle/activity/armory?activityId=100301", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), `"data":true`) {
		t.Fatalf("expected success data true, got %s", recorder.Body.String())
	}
	if activityArmory.activityID != 100301 {
		t.Fatalf("expected activity armory activity id 100301, got %d", activityArmory.activityID)
	}
	if strategyArmory.activityID != 100301 {
		t.Fatalf("expected strategy armory activity id 100301, got %d", strategyArmory.activityID)
	}
}

func TestActivityArmoryRouteReturnsActivityArmoryAppErrorCode(t *testing.T) {
	router := NewRouter(RouterOptions{
		ActivityArmoryService: &fakeActivityArmoryService{
			err: types.NewAppError(types.ResponseCodeActivityStateError, nil),
		},
		ActivityStrategyArmoryService: &fakeActivityStrategyArmoryService{},
	})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/raffle/activity/armory?activityId=100301", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), `"code":"ERR_BIZ_003"`) {
		t.Fatalf("expected app error code, got %s", recorder.Body.String())
	}
}

func TestActivityArmoryRouteReturnsStrategyArmoryAppErrorCode(t *testing.T) {
	router := NewRouter(RouterOptions{
		ActivityArmoryService: &fakeActivityArmoryService{},
		ActivityStrategyArmoryService: &fakeActivityStrategyArmoryService{
			err: types.NewAppError(types.ResponseCodeUnassembledStrategy, nil),
		},
	})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/raffle/activity/armory?activityId=100301", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), `"code":"ERR_BIZ_002"`) {
		t.Fatalf("expected app error code, got %s", recorder.Body.String())
	}
}

func TestActivityDrawRoute(t *testing.T) {
	router := NewRouter(RouterOptions{
		ActivityDrawService: &fakeActivityDrawService{
			result: activity.DrawResult{
				AwardID:    101,
				AwardTitle: "积分",
				AwardIndex: 1,
			},
		},
	})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/raffle/activity/draw", strings.NewReader(`{"userId":"xiaofuge","activityId":100301}`))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), `"awardId":101`) {
		t.Fatalf("expected award id, got %s", recorder.Body.String())
	}
	if !strings.Contains(recorder.Body.String(), `"awardTitle":"积分"`) {
		t.Fatalf("expected award title, got %s", recorder.Body.String())
	}
}

func TestActivityDrawRouteReturnsAppErrorCode(t *testing.T) {
	router := NewRouter(RouterOptions{
		ActivityDrawService: &fakeActivityDrawService{
			err: types.NewAppError(types.ResponseCodeActivityDateError, nil),
		},
	})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/raffle/activity/draw", strings.NewReader(`{"userId":"xiaofuge","activityId":100301}`))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), `"code":"ERR_BIZ_004"`) {
		t.Fatalf("expected app error code, got %s", recorder.Body.String())
	}
}

func TestActivityDrawRouteIllegalParam(t *testing.T) {
	router := NewRouter(RouterOptions{
		ActivityDrawService: &fakeActivityDrawService{},
	})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/raffle/activity/draw", strings.NewReader(`{"userId":"xiaofuge","activityId":0}`))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), `"code":"0002"`) {
		t.Fatalf("expected illegal param code, got %s", recorder.Body.String())
	}
}

func TestQuerySkuProductListByActivityIDRoute(t *testing.T) {
	router := NewRouter(RouterOptions{
		ActivitySkuProductService: &fakeActivitySkuProductService{
			products: []activity.SkuProductEntity{
				{
					SKU:               9011,
					ActivityID:        100301,
					ActivityCountID:   11101,
					StockCount:        100000,
					StockCountSurplus: 99890,
					ProductAmount:     1.68,
					ActivityCount: activity.ActivityCountEntity{
						TotalCount: 100,
						DayCount:   100,
						MonthCount: 100,
					},
				},
			},
		},
	})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/raffle/activity/query_sku_product_list_by_activity_id?activityId=100301", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), `"sku":9011`) {
		t.Fatalf("expected sku, got %s", recorder.Body.String())
	}
	if !strings.Contains(recorder.Body.String(), `"productAmount":1.68`) {
		t.Fatalf("expected product amount, got %s", recorder.Body.String())
	}
	if !strings.Contains(recorder.Body.String(), `"totalCount":100`) {
		t.Fatalf("expected total count, got %s", recorder.Body.String())
	}
}

func TestQuerySkuProductListByActivityIDRouteReturnsAppErrorCode(t *testing.T) {
	router := NewRouter(RouterOptions{
		ActivitySkuProductService: &fakeActivitySkuProductService{
			err: types.NewAppError(types.ResponseCodeActivityStateError, nil),
		},
	})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/raffle/activity/query_sku_product_list_by_activity_id?activityId=100301", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), `"code":"ERR_BIZ_003"`) {
		t.Fatalf("expected app error code, got %s", recorder.Body.String())
	}
}

func TestCreditPayExchangeSkuRoute(t *testing.T) {
	exchange := &fakeActivityExchangeService{result: true}
	router := NewRouter(RouterOptions{
		ActivityExchangeService: exchange,
	})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/raffle/activity/credit_pay_exchange_sku", strings.NewReader(`{"userId":"xiaofuge","sku":9011}`))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), `"data":true`) {
		t.Fatalf("expected success data true, got %s", recorder.Body.String())
	}
	if exchange.userID != "xiaofuge" || exchange.sku != 9011 {
		t.Fatalf("expected exchange request, got %s/%d", exchange.userID, exchange.sku)
	}
}

func TestCreditPayExchangeSkuRouteReturnsAppErrorCode(t *testing.T) {
	router := NewRouter(RouterOptions{
		ActivityExchangeService: &fakeActivityExchangeService{
			err: types.NewAppError(types.ResponseCodeAccountQuotaError, nil),
		},
	})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/raffle/activity/credit_pay_exchange_sku", strings.NewReader(`{"userId":"xiaofuge","sku":9011}`))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), `"code":"ERR_BIZ_006"`) {
		t.Fatalf("expected app error code, got %s", recorder.Body.String())
	}
}

func TestCreditPayExchangeSkuRouteIllegalParam(t *testing.T) {
	router := NewRouter(RouterOptions{
		ActivityExchangeService: &fakeActivityExchangeService{},
	})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/raffle/activity/credit_pay_exchange_sku", strings.NewReader(`{"userId":"xiaofuge","sku":0}`))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), `"code":"0002"`) {
		t.Fatalf("expected illegal param code, got %s", recorder.Body.String())
	}
}

func TestQueryUserCreditAccountRoute(t *testing.T) {
	router := NewRouter(RouterOptions{
		ActivityCreditService: &fakeActivityCreditService{amount: 12.35},
	})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/raffle/activity/query_user_credit_account?userId=xiaofuge", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), `"data":12.35`) {
		t.Fatalf("expected credit amount, got %s", recorder.Body.String())
	}
}

func TestQueryUserCreditAccountRouteReturnsAppErrorCode(t *testing.T) {
	router := NewRouter(RouterOptions{
		ActivityCreditService: &fakeActivityCreditService{
			err: types.NewAppError(types.ResponseCodeAccountQuotaError, nil),
		},
	})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/raffle/activity/query_user_credit_account?userId=xiaofuge", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), `"code":"ERR_BIZ_006"`) {
		t.Fatalf("expected app error code, got %s", recorder.Body.String())
	}
}

func TestCalendarSignRebateRoute(t *testing.T) {
	rebateService := &fakeActivityRebateService{signResult: true}
	router := NewRouter(RouterOptions{
		ActivityRebateService: rebateService,
	})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/raffle/activity/calendar_sign_rebate", strings.NewReader("userId=xiaofuge"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), `"data":true`) {
		t.Fatalf("expected sign result, got %s", recorder.Body.String())
	}
	if rebateService.signUserID != "xiaofuge" {
		t.Fatalf("expected sign user, got %s", rebateService.signUserID)
	}
}

func TestCalendarSignRebateRouteIllegalParam(t *testing.T) {
	router := NewRouter(RouterOptions{
		ActivityRebateService: &fakeActivityRebateService{},
	})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/raffle/activity/calendar_sign_rebate", strings.NewReader("userId="))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), `"code":"0002"`) {
		t.Fatalf("expected illegal param code, got %s", recorder.Body.String())
	}
}

func TestIsCalendarSignRebateRoute(t *testing.T) {
	rebateService := &fakeActivityRebateService{queryResult: true}
	router := NewRouter(RouterOptions{
		ActivityRebateService: rebateService,
	})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/raffle/activity/is_calendar_sign_rebate", strings.NewReader("userId=xiaofuge"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), `"data":true`) {
		t.Fatalf("expected signed result, got %s", recorder.Body.String())
	}
	if rebateService.queryUserID != "xiaofuge" {
		t.Fatalf("expected query user, got %s", rebateService.queryUserID)
	}
}

type fakeActivityAccountService struct {
	account activity.AccountEntity
	err     error
}

func (f *fakeActivityAccountService) QueryActivityAccount(ctx context.Context, activityID int64, userID string) (activity.AccountEntity, error) {
	if f.err != nil {
		return activity.AccountEntity{}, f.err
	}
	return f.account, nil
}

type fakeActivitySkuProductService struct {
	products []activity.SkuProductEntity
	err      error
}

func (f *fakeActivitySkuProductService) QuerySkuProductListByActivityID(ctx context.Context, activityID int64) ([]activity.SkuProductEntity, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.products, nil
}

type fakeActivityArmoryService struct {
	activityID int64
	err        error
}

func (f *fakeActivityArmoryService) AssembleActivitySkuByActivityID(ctx context.Context, activityID int64) error {
	f.activityID = activityID
	if f.err != nil {
		return f.err
	}
	return nil
}

type fakeActivityStrategyArmoryService struct {
	activityID int64
	err        error
}

func (f *fakeActivityStrategyArmoryService) AssembleLotteryStrategyByActivityID(ctx context.Context, activityID int64) error {
	f.activityID = activityID
	if f.err != nil {
		return f.err
	}
	return nil
}

type fakeActivityDrawService struct {
	result activity.DrawResult
	err    error
}

func (f *fakeActivityDrawService) Draw(ctx context.Context, userID string, activityID int64) (activity.DrawResult, error) {
	if f.err != nil {
		return activity.DrawResult{}, f.err
	}
	return f.result, nil
}

type fakeActivityExchangeService struct {
	userID string
	sku    int64
	result bool
	err    error
}

func (f *fakeActivityExchangeService) CreditPayExchangeSku(ctx context.Context, userID string, sku int64) (bool, error) {
	f.userID = userID
	f.sku = sku
	if f.err != nil {
		return false, f.err
	}
	return f.result, nil
}

type fakeActivityCreditService struct {
	amount float64
	err    error
}

func (f *fakeActivityCreditService) QueryUserCreditAccount(ctx context.Context, userID string) (float64, error) {
	if f.err != nil {
		return 0, f.err
	}
	return f.amount, nil
}

type fakeActivityRebateService struct {
	signUserID  string
	queryUserID string
	signResult  bool
	queryResult bool
}

func (f *fakeActivityRebateService) CalendarSignRebate(ctx context.Context, userID string) (bool, error) {
	f.signUserID = userID
	return f.signResult, nil
}

func (f *fakeActivityRebateService) IsCalendarSignRebate(ctx context.Context, userID string) (bool, error) {
	f.queryUserID = userID
	return f.queryResult, nil
}
