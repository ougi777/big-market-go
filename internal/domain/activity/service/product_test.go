package service

import (
	"context"
	"testing"

	"bm-go/internal/domain/activity"
)

func TestSkuProductServiceQuerySkuProductListByActivityID(t *testing.T) {
	repo := &fakeSkuProductRepository{
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
	}
	service := NewSkuProductService(repo)

	products, err := service.QuerySkuProductListByActivityID(context.Background(), 100301)
	if err != nil {
		t.Fatalf("query sku product list: %v", err)
	}

	if repo.activityID != 100301 {
		t.Fatalf("expected activity id 100301, got %d", repo.activityID)
	}
	if len(products) != 1 {
		t.Fatalf("expected one product, got %d", len(products))
	}
	if products[0].SKU != 9011 || products[0].ActivityCount.TotalCount != 100 {
		t.Fatalf("expected sku product, got %+v", products[0])
	}
}

type fakeSkuProductRepository struct {
	activityID int64
	products   []activity.SkuProductEntity
}

func (f *fakeSkuProductRepository) QuerySkuProductListByActivityID(ctx context.Context, activityID int64) ([]activity.SkuProductEntity, error) {
	f.activityID = activityID
	return f.products, nil
}
