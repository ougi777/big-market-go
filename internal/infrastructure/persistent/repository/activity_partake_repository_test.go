package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"bm-go/internal/domain/activity"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/gorm"
)

func TestActivityRepositoryQueryNoUsedRaffleOrder(t *testing.T) {
	db, mock := newMockGormDB(t)
	repo := NewActivityRepository(db)
	orderTime := time.Date(2026, 5, 25, 10, 0, 0, 0, time.Local)

	rows := sqlmock.NewRows([]string{
		"user_id",
		"activity_id",
		"activity_name",
		"strategy_id",
		"order_id",
		"order_time",
		"order_state",
	}).AddRow("xiaofuge", 100301, "大营销抽奖", 100006, "order-001", orderTime, activity.UserRaffleOrderCreate)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `user_id`,`activity_id`,`activity_name`,`strategy_id`,`order_id`,`order_time`,`order_state` FROM `user_raffle_order` WHERE user_id = ? and activity_id = ? and order_state = ? ORDER BY `user_raffle_order`.`id` LIMIT ?")).
		WithArgs("xiaofuge", int64(100301), activity.UserRaffleOrderCreate, 1).
		WillReturnRows(rows)

	order, exists, err := repo.QueryNoUsedRaffleOrder(context.Background(), "xiaofuge", 100301)
	if err != nil {
		t.Fatalf("query no used raffle order: %v", err)
	}
	if !exists {
		t.Fatal("expected raffle order exists")
	}
	if order.OrderID != "order-001" || order.StrategyID != 100006 || order.OrderState != activity.UserRaffleOrderCreate {
		t.Fatalf("unexpected raffle order: %+v", order)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestActivityRepositoryQueryNoUsedRaffleOrderNotFound(t *testing.T) {
	db, mock := newMockGormDB(t)
	repo := NewActivityRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `user_id`,`activity_id`,`activity_name`,`strategy_id`,`order_id`,`order_time`,`order_state` FROM `user_raffle_order` WHERE user_id = ? and activity_id = ? and order_state = ? ORDER BY `user_raffle_order`.`id` LIMIT ?")).
		WithArgs("xiaofuge", int64(100301), activity.UserRaffleOrderCreate, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	order, exists, err := repo.QueryNoUsedRaffleOrder(context.Background(), "xiaofuge", 100301)
	if err != nil {
		t.Fatalf("query no used raffle order: %v", err)
	}
	if exists {
		t.Fatalf("expected raffle order missing, got %+v", order)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}
