package service

import (
	"context"
	"testing"
	"time"

	"bm-go/internal/domain/activity"
	"bm-go/internal/domain/rebate"
)

func TestRebateProcessorProcessSKU(t *testing.T) {
	now := time.Date(2026, 5, 25, 10, 0, 0, 0, time.Local)
	repo := &fakeRebateProcessRepository{
		product: activity.SkuProductEntity{
			SKU:        9011,
			ActivityID: 100301,
			ActivityCount: activity.ActivityCountEntity{
				TotalCount: 10,
				DayCount:   1,
				MonthCount: 5,
			},
		},
		productExists: true,
		activity: activity.ActivityEntity{
			ActivityID:   100301,
			ActivityName: "test activity",
			StrategyID:   100006,
			State:        activity.ActivityStateOpen,
		},
		activityExists: true,
	}
	processor := NewRebateProcessor(repo, repo)
	processor.now = func() time.Time { return now }
	processor.orderIDGenerator = func() (string, error) { return "123456789012", nil }

	err := processor.ProcessRebate(context.Background(), rebate.SendRebateMessage{
		UserID:       "xiaofuge",
		RebateType:   rebate.RebateTypeSKU,
		RebateConfig: "9011",
		BizID:        "xiaofuge_sku_20260525",
	})
	if err != nil {
		t.Fatalf("process sku rebate: %v", err)
	}

	if !repo.savedSKU {
		t.Fatalf("expected saved sku rebate order")
	}
	if repo.skuAggregate.ActivityOrder.State != activity.ActivityOrderCompleted {
		t.Fatalf("expected completed order, got %+v", repo.skuAggregate.ActivityOrder)
	}
	if repo.skuAggregate.ActivityOrder.PayAmount != 0 {
		t.Fatalf("expected no pay order, got %.2f", repo.skuAggregate.ActivityOrder.PayAmount)
	}
	if repo.skuAggregate.ActivityOrder.OutBusinessNo != "xiaofuge_sku_20260525" {
		t.Fatalf("expected biz id, got %s", repo.skuAggregate.ActivityOrder.OutBusinessNo)
	}
}

func TestRebateProcessorProcessIntegral(t *testing.T) {
	repo := &fakeRebateProcessRepository{}
	processor := NewRebateProcessor(repo, repo)
	processor.orderIDGenerator = func() (string, error) { return "123456789012", nil }

	err := processor.ProcessRebate(context.Background(), rebate.SendRebateMessage{
		UserID:       "xiaofuge",
		RebateType:   rebate.RebateTypeIntegral,
		RebateConfig: "10",
		BizID:        "xiaofuge_integral_20260525",
	})
	if err != nil {
		t.Fatalf("process integral rebate: %v", err)
	}

	if !repo.savedIntegral {
		t.Fatalf("expected saved integral rebate")
	}
	if repo.integral.TradeAmount != 10 || repo.integral.OutBusinessNo != "xiaofuge_integral_20260525" {
		t.Fatalf("expected integral order, got %+v", repo.integral)
	}
}

type fakeRebateProcessRepository struct {
	product        activity.SkuProductEntity
	productExists  bool
	activity       activity.ActivityEntity
	activityExists bool
	savedSKU       bool
	skuAggregate   activity.CreateRebateSkuOrderAggregate
	savedIntegral  bool
	integral       activity.RebateIntegralEntity
}

func (f *fakeRebateProcessRepository) QuerySkuProductBySKU(ctx context.Context, sku int64) (activity.SkuProductEntity, bool, error) {
	return f.product, f.productExists, nil
}

func (f *fakeRebateProcessRepository) QueryActivityByActivityID(ctx context.Context, activityID int64) (activity.ActivityEntity, bool, error) {
	return f.activity, f.activityExists, nil
}

func (f *fakeRebateProcessRepository) SaveRebateSkuOrder(ctx context.Context, aggregate activity.CreateRebateSkuOrderAggregate) error {
	f.savedSKU = true
	f.skuAggregate = aggregate
	return nil
}

func (f *fakeRebateProcessRepository) SaveRebateIntegralOrder(ctx context.Context, rebateIntegral activity.RebateIntegralEntity) error {
	f.savedIntegral = true
	f.integral = rebateIntegral
	return nil
}
