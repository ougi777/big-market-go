package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"bm-go/internal/domain/activity"
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

type fakeActivityAccountService struct {
	account activity.AccountEntity
}

func (f *fakeActivityAccountService) QueryActivityAccount(ctx context.Context, activityID int64, userID string) (activity.AccountEntity, error) {
	return f.account, nil
}

type fakeActivitySkuProductService struct {
	products []activity.SkuProductEntity
}

func (f *fakeActivitySkuProductService) QuerySkuProductListByActivityID(ctx context.Context, activityID int64) ([]activity.SkuProductEntity, error) {
	return f.products, nil
}

type fakeActivityArmoryService struct {
	activityID int64
}

func (f *fakeActivityArmoryService) AssembleActivitySkuByActivityID(ctx context.Context, activityID int64) error {
	f.activityID = activityID
	return nil
}

type fakeActivityStrategyArmoryService struct {
	activityID int64
}

func (f *fakeActivityStrategyArmoryService) AssembleLotteryStrategyByActivityID(ctx context.Context, activityID int64) error {
	f.activityID = activityID
	return nil
}

type fakeActivityDrawService struct {
	result activity.DrawResult
}

func (f *fakeActivityDrawService) Draw(ctx context.Context, userID string, activityID int64) (activity.DrawResult, error) {
	return f.result, nil
}
