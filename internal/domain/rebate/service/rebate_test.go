package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"bm-go/internal/domain/rebate"
)

func TestCalendarSignRebate(t *testing.T) {
	repo := &fakeRebateRepository{
		configs: []rebate.DailyBehaviorRebateEntity{
			{
				BehaviorType: rebate.BehaviorTypeSign,
				RebateDesc:   "签到返利 SKU",
				RebateType:   "sku",
				RebateConfig: "9011",
			},
			{
				BehaviorType: rebate.BehaviorTypeSign,
				RebateDesc:   "签到返利积分",
				RebateType:   "integral",
				RebateConfig: "10",
			},
		},
	}
	publisher := &fakeRebatePublisher{}
	service := NewRebateService(repo, publisher)
	service.now = func() time.Time { return time.Date(2026, 5, 25, 9, 30, 0, 0, time.Local) }
	ids := []string{"123456789012", "12345678901", "223456789012", "22345678901"}
	service.newID = func(length int) (string, error) {
		id := ids[0]
		ids = ids[1:]
		return id, nil
	}

	result, err := service.CalendarSignRebate(context.Background(), "xiaofuge")
	if err != nil {
		t.Fatalf("calendar sign rebate: %v", err)
	}

	if !result {
		t.Fatalf("expected sign rebate success")
	}
	if len(repo.saved) != 2 {
		t.Fatalf("expected 2 aggregates, got %d", len(repo.saved))
	}
	if repo.saved[0].Order.OutBusinessNo != "20260525" || repo.saved[0].Order.BizID != "xiaofuge_sku_20260525" {
		t.Fatalf("expected sku order, got %+v", repo.saved[0].Order)
	}
	if repo.saved[1].Order.BizID != "xiaofuge_integral_20260525" {
		t.Fatalf("expected integral order, got %+v", repo.saved[1].Order)
	}
	if len(publisher.messages) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(publisher.messages))
	}
	if !strings.Contains(publisher.messages[0], `"rebateType":"sku"`) {
		t.Fatalf("expected sku message, got %s", publisher.messages[0])
	}
	if repo.completed[0] != "12345678901" || repo.completed[1] != "22345678901" {
		t.Fatalf("expected completed messages, got %+v", repo.completed)
	}
}

func TestIsCalendarSignRebate(t *testing.T) {
	repo := &fakeRebateRepository{
		orders: []rebate.BehaviorRebateOrderEntity{
			{UserID: "xiaofuge", OutBusinessNo: "20260525"},
		},
	}
	service := NewRebateService(repo, &fakeRebatePublisher{})
	service.now = func() time.Time { return time.Date(2026, 5, 25, 9, 30, 0, 0, time.Local) }

	result, err := service.IsCalendarSignRebate(context.Background(), "xiaofuge")
	if err != nil {
		t.Fatalf("is calendar sign rebate: %v", err)
	}
	if !result {
		t.Fatalf("expected signed")
	}
	if repo.queryUserID != "xiaofuge" || repo.queryOutBusinessNo != "20260525" {
		t.Fatalf("expected query today order, got %s/%s", repo.queryUserID, repo.queryOutBusinessNo)
	}
}

func TestCalendarSignRebateMarksTaskFailOnPublishError(t *testing.T) {
	repo := &fakeRebateRepository{
		configs: []rebate.DailyBehaviorRebateEntity{
			{
				BehaviorType: rebate.BehaviorTypeSign,
				RebateDesc:   "签到返利积分",
				RebateType:   "integral",
				RebateConfig: "10",
			},
		},
	}
	publisher := &fakeRebatePublisher{err: errors.New("publish failed")}
	service := NewRebateService(repo, publisher)
	service.now = func() time.Time { return time.Date(2026, 5, 25, 9, 30, 0, 0, time.Local) }
	ids := []string{"123456789012", "12345678901"}
	service.newID = func(length int) (string, error) {
		id := ids[0]
		ids = ids[1:]
		return id, nil
	}

	result, err := service.CalendarSignRebate(context.Background(), "xiaofuge")
	if err != nil {
		t.Fatalf("calendar sign rebate: %v", err)
	}
	if !result {
		t.Fatal("expected sign rebate success")
	}
	if len(repo.failed) != 1 || repo.failed[0] != "12345678901" {
		t.Fatalf("expected failed task, got %+v", repo.failed)
	}
}

func TestCalendarSignRebateIllegalParam(t *testing.T) {
	service := NewRebateService(&fakeRebateRepository{}, &fakeRebatePublisher{})

	_, err := service.CalendarSignRebate(context.Background(), " ")
	if err == nil {
		t.Fatal("expected illegal param error")
	}
}

func TestCalendarSignRebateReturnsFalseWhenConfigEmpty(t *testing.T) {
	service := NewRebateService(&fakeRebateRepository{}, &fakeRebatePublisher{})

	result, err := service.CalendarSignRebate(context.Background(), "xiaofuge")
	if err != nil {
		t.Fatalf("calendar sign rebate: %v", err)
	}
	if result {
		t.Fatal("expected no rebate")
	}
}

func TestCalendarSignRebateConfigError(t *testing.T) {
	repo := &fakeRebateRepository{configErr: errors.New("query config failed")}
	service := NewRebateService(repo, &fakeRebatePublisher{})

	_, err := service.CalendarSignRebate(context.Background(), "xiaofuge")
	if err == nil {
		t.Fatal("expected config error")
	}
}

func TestCalendarSignRebateSaveError(t *testing.T) {
	repo := &fakeRebateRepository{
		configs: []rebate.DailyBehaviorRebateEntity{
			{BehaviorType: rebate.BehaviorTypeSign, RebateType: rebate.RebateTypeIntegral, RebateConfig: "10"},
		},
		saveErr: errors.New("save failed"),
	}
	service := NewRebateService(repo, &fakeRebatePublisher{})
	service.newID = func(length int) (string, error) { return "123456789012", nil }

	_, err := service.CalendarSignRebate(context.Background(), "xiaofuge")
	if err == nil {
		t.Fatal("expected save error")
	}
}

func TestIsCalendarSignRebateIllegalParam(t *testing.T) {
	service := NewRebateService(&fakeRebateRepository{}, &fakeRebatePublisher{})

	_, err := service.IsCalendarSignRebate(context.Background(), " ")
	if err == nil {
		t.Fatal("expected illegal param error")
	}
}

func TestIsCalendarSignRebateRepositoryError(t *testing.T) {
	repo := &fakeRebateRepository{queryErr: errors.New("query order failed")}
	service := NewRebateService(repo, &fakeRebatePublisher{})

	_, err := service.IsCalendarSignRebate(context.Background(), "xiaofuge")
	if err == nil {
		t.Fatal("expected query order error")
	}
}

type fakeRebateRepository struct {
	configs            []rebate.DailyBehaviorRebateEntity
	orders             []rebate.BehaviorRebateOrderEntity
	saved              []rebate.BehaviorRebateAggregate
	completed          []string
	failed             []string
	queryUserID        string
	queryOutBusinessNo string
	configErr          error
	saveErr            error
	queryErr           error
}

func (f *fakeRebateRepository) QueryDailyBehaviorRebateConfig(ctx context.Context, behaviorType string) ([]rebate.DailyBehaviorRebateEntity, error) {
	return f.configs, f.configErr
}

func (f *fakeRebateRepository) SaveUserRebateRecords(ctx context.Context, aggregates []rebate.BehaviorRebateAggregate) error {
	f.saved = aggregates
	return f.saveErr
}

func (f *fakeRebateRepository) QueryOrderByOutBusinessNo(ctx context.Context, userID string, outBusinessNo string) ([]rebate.BehaviorRebateOrderEntity, error) {
	f.queryUserID = userID
	f.queryOutBusinessNo = outBusinessNo
	return f.orders, f.queryErr
}

func (f *fakeRebateRepository) UpdateTaskSendMessageCompleted(ctx context.Context, userID string, messageID string) error {
	f.completed = append(f.completed, messageID)
	return nil
}

func (f *fakeRebateRepository) UpdateTaskSendMessageFail(ctx context.Context, userID string, messageID string) error {
	f.failed = append(f.failed, messageID)
	return nil
}

type fakeRebatePublisher struct {
	messages []string
	err      error
}

func (f *fakeRebatePublisher) Publish(ctx context.Context, topic string, message string) error {
	f.messages = append(f.messages, message)
	if f.err != nil {
		return f.err
	}
	return nil
}
