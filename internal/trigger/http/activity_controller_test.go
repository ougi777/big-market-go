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
