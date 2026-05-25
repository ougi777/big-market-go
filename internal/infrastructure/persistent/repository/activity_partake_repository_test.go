package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"bm-go/internal/domain/activity"
	"bm-go/internal/types"

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

func TestActivityRepositorySaveCreatePartakeOrderQuotaNotEnough(t *testing.T) {
	db, mock := newMockGormDB(t)
	repo := NewActivityRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `raffle_activity_account` SET `total_count_surplus`=total_count_surplus - ? WHERE user_id = ? and activity_id = ? and total_count_surplus > 0")).
		WithArgs(1, "xiaofuge", int64(100301)).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectRollback()

	err := repo.SaveCreatePartakeOrder(context.Background(), activity.CreatePartakeOrderAggregate{
		UserID:     "xiaofuge",
		ActivityID: 100301,
		UserRaffleOrder: activity.UserRaffleOrderEntity{
			UserID:       "xiaofuge",
			ActivityID:   100301,
			ActivityName: "大营销抽奖",
			StrategyID:   100006,
			OrderID:      "order-001",
			OrderTime:    time.Date(2026, 5, 25, 10, 0, 0, 0, time.Local),
			OrderState:   activity.UserRaffleOrderCreate,
		},
	})

	var appErr types.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected app error, got %v", err)
	}
	if appErr.Code != types.ResponseCodeAccountQuotaError {
		t.Fatalf("expected account quota code, got %s", appErr.Code.Code)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestActivityRepositorySaveCreatePartakeOrder(t *testing.T) {
	db, mock := newMockGormDB(t)
	repo := NewActivityRepository(db)
	orderTime := time.Date(2026, 5, 25, 10, 0, 0, 0, time.Local)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `raffle_activity_account` SET `total_count_surplus`=total_count_surplus - ? WHERE user_id = ? and activity_id = ? and total_count_surplus > 0")).
		WithArgs(1, "xiaofuge", int64(100301)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `raffle_activity_account_month` SET `month_count_surplus`=month_count_surplus - ? WHERE user_id = ? and activity_id = ? and month = ? and month_count_surplus > 0")).
		WithArgs(1, "xiaofuge", int64(100301), "2026-05").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `raffle_activity_account` SET `month_count_surplus`=month_count_surplus - ? WHERE user_id = ? and activity_id = ? and month_count_surplus > 0")).
		WithArgs(1, "xiaofuge", int64(100301)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `raffle_activity_account_day` SET `day_count_surplus`=day_count_surplus - ? WHERE user_id = ? and activity_id = ? and day = ? and day_count_surplus > 0")).
		WithArgs(1, "xiaofuge", int64(100301), "2026-05-25").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `raffle_activity_account` SET `day_count_surplus`=day_count_surplus - ? WHERE user_id = ? and activity_id = ? and day_count_surplus > 0")).
		WithArgs(1, "xiaofuge", int64(100301)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `user_raffle_order`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.SaveCreatePartakeOrder(context.Background(), activity.CreatePartakeOrderAggregate{
		UserID:            "xiaofuge",
		ActivityID:        100301,
		ExistAccountMonth: true,
		ActivityAccountMonth: activity.AccountMonthEntity{
			UserID:            "xiaofuge",
			ActivityID:        100301,
			Month:             "2026-05",
			MonthCount:        30,
			MonthCountSurplus: 20,
		},
		ExistAccountDay: true,
		ActivityAccountDay: activity.AccountDayEntity{
			UserID:          "xiaofuge",
			ActivityID:      100301,
			Day:             "2026-05-25",
			DayCount:        10,
			DayCountSurplus: 8,
		},
		UserRaffleOrder: activity.UserRaffleOrderEntity{
			UserID:       "xiaofuge",
			ActivityID:   100301,
			ActivityName: "大营销抽奖",
			StrategyID:   100006,
			OrderID:      "order-001",
			OrderTime:    orderTime,
			OrderState:   activity.UserRaffleOrderCreate,
		},
	})
	if err != nil {
		t.Fatalf("save create partake order: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}
