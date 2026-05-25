package service

import (
	"context"
	"strings"
	"testing"
	"time"

	"bm-go/internal/domain/activity"
)

func TestExchangeServiceCreditPayExchangeSku(t *testing.T) {
	now := time.Date(2026, 5, 25, 10, 0, 0, 0, time.Local)
	repo := &fakeExchangeRepository{
		product: activity.SkuProductEntity{
			SKU:           9011,
			ActivityID:    100301,
			ProductAmount: 1.68,
			ActivityCount: activity.ActivityCountEntity{
				TotalCount: 10,
				DayCount:   1,
				MonthCount: 5,
			},
		},
		productExists: true,
		activity: activity.ActivityEntity{
			ActivityID:    100301,
			ActivityName:  "test activity",
			BeginDateTime: now.Add(-time.Hour),
			EndDateTime:   now.Add(time.Hour),
			StrategyID:    100006,
			State:         activity.ActivityStateOpen,
		},
		activityExists: true,
	}
	stock := &fakeExchangeStockService{ok: true}
	publisher := &fakeExchangePublisher{}
	service := NewExchangeService(repo, stock, publisher, repo)
	service.now = func() time.Time { return now }
	service.orderIDGenerator = func() (string, error) { return "123456789012", nil }
	service.messageIDGenerator = func() (string, error) { return "22222222222", nil }
	service.businessNoGenerator = func() (string, error) { return "987654321098", nil }
	service.creditOrderGenerator = func() (string, error) { return "111111111111", nil }

	ok, err := service.CreditPayExchangeSku(context.Background(), "xiaofuge", 9011)
	if err != nil {
		t.Fatalf("credit pay exchange sku: %v", err)
	}
	if !ok {
		t.Fatal("expected exchange ok")
	}
	if !repo.saved || !repo.completed {
		t.Fatalf("expected saved and completed")
	}
	if stock.sku != 9011 || stock.activityID != 100301 {
		t.Fatalf("expected stock subtract, got %d/%d", stock.sku, stock.activityID)
	}
	if repo.createAggregate.ActivityOrder.State != activity.ActivityOrderWaitPay {
		t.Fatalf("expected wait pay order, got %+v", repo.createAggregate.ActivityOrder)
	}
	if repo.completeAggregate.CreditOrder.TradeAmount != -1.68 {
		t.Fatalf("expected credit amount -1.68, got %.2f", repo.completeAggregate.CreditOrder.TradeAmount)
	}
	if repo.completeAggregate.SendTask.Topic != "credit_adjust_success" || repo.completeAggregate.SendTask.MessageID != "22222222222" {
		t.Fatalf("expected send task, got %+v", repo.completeAggregate.SendTask)
	}
	if publisher.topic != "credit_adjust_success" {
		t.Fatalf("expected credit adjust success topic, got %s", publisher.topic)
	}
	if !strings.Contains(publisher.message, `"outBusinessNo":"987654321098"`) {
		t.Fatalf("expected adjust success message, got %s", publisher.message)
	}
	if repo.completedMessageID != "22222222222" {
		t.Fatalf("expected completed task, got %s", repo.completedMessageID)
	}
}

type fakeExchangeRepository struct {
	product            activity.SkuProductEntity
	productExists      bool
	activity           activity.ActivityEntity
	activityExists     bool
	unpaid             activity.SkuExchangeOrderEntity
	unpaidExists       bool
	saved              bool
	completed          bool
	completedMessageID string
	failMessageID      string
	createAggregate    activity.CreateSkuExchangeOrderAggregate
	completeAggregate  activity.CompleteSkuExchangeAggregate
}

func (f *fakeExchangeRepository) QuerySkuProductBySKU(ctx context.Context, sku int64) (activity.SkuProductEntity, bool, error) {
	return f.product, f.productExists, nil
}

func (f *fakeExchangeRepository) QueryActivityByActivityID(ctx context.Context, activityID int64) (activity.ActivityEntity, bool, error) {
	return f.activity, f.activityExists, nil
}

func (f *fakeExchangeRepository) QueryUnpaidActivityOrder(ctx context.Context, userID string, sku int64) (activity.SkuExchangeOrderEntity, bool, error) {
	return f.unpaid, f.unpaidExists, nil
}

func (f *fakeExchangeRepository) SaveCreditPayOrder(ctx context.Context, aggregate activity.CreateSkuExchangeOrderAggregate) error {
	f.saved = true
	f.createAggregate = aggregate
	return nil
}

func (f *fakeExchangeRepository) CompleteCreditPayOrder(ctx context.Context, aggregate activity.CompleteSkuExchangeAggregate) error {
	f.completed = true
	f.completeAggregate = aggregate
	return nil
}

func (f *fakeExchangeRepository) UpdateTaskSendMessageCompleted(ctx context.Context, userID string, messageID string) error {
	f.completedMessageID = messageID
	return nil
}

func (f *fakeExchangeRepository) UpdateTaskSendMessageFail(ctx context.Context, userID string, messageID string) error {
	f.failMessageID = messageID
	return nil
}

type fakeExchangeStockService struct {
	ok         bool
	sku        int64
	activityID int64
}

func (f *fakeExchangeStockService) SubtractActivitySkuStock(ctx context.Context, sku int64, activityID int64) (bool, error) {
	f.sku = sku
	f.activityID = activityID
	return f.ok, nil
}

type fakeExchangePublisher struct {
	topic   string
	message string
}

func (f *fakeExchangePublisher) Publish(ctx context.Context, topic string, message string) error {
	f.topic = topic
	f.message = message
	return nil
}
