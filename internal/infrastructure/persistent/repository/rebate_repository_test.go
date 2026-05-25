package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"bm-go/internal/domain/rebate"
	"bm-go/internal/types"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
)

func TestRebateRepositoryQueryDailyBehaviorRebateConfig(t *testing.T) {
	db, mock := newMockGormDB(t)
	repo := NewRebateRepository(db)

	rows := sqlmock.NewRows([]string{"behavior_type", "rebate_desc", "rebate_type", "rebate_config"}).
		AddRow("sign", "签到返利积分", "integral", "10")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `behavior_type`,`rebate_desc`,`rebate_type`,`rebate_config` FROM `daily_behavior_rebate` WHERE behavior_type = ? and state = ?")).
		WithArgs("sign", "open").
		WillReturnRows(rows)

	configs, err := repo.QueryDailyBehaviorRebateConfig(context.Background(), "sign")
	if err != nil {
		t.Fatalf("query daily behavior rebate config: %v", err)
	}
	if len(configs) != 1 || configs[0].RebateType != "integral" || configs[0].RebateConfig != "10" {
		t.Fatalf("unexpected rebate configs: %+v", configs)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestRebateRepositoryQueryOrderByOutBusinessNo(t *testing.T) {
	db, mock := newMockGormDB(t)
	repo := NewRebateRepository(db)

	rows := sqlmock.NewRows([]string{
		"user_id",
		"order_id",
		"behavior_type",
		"rebate_desc",
		"rebate_type",
		"rebate_config",
		"out_business_no",
		"biz_id",
	}).AddRow("xiaofuge", "order-001", "sign", "签到返利积分", "integral", "10", "20260525", "xiaofuge_integral_20260525")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `user_id`,`order_id`,`behavior_type`,`rebate_desc`,`rebate_type`,`rebate_config`,`out_business_no`,`biz_id` FROM `user_behavior_rebate_order` WHERE user_id = ? and out_business_no = ?")).
		WithArgs("xiaofuge", "20260525").
		WillReturnRows(rows)

	orders, err := repo.QueryOrderByOutBusinessNo(context.Background(), "xiaofuge", "20260525")
	if err != nil {
		t.Fatalf("query order by out business no: %v", err)
	}
	if len(orders) != 1 || orders[0].BizID != "xiaofuge_integral_20260525" {
		t.Fatalf("unexpected rebate orders: %+v", orders)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestRebateRepositorySaveUserRebateRecordsDuplicate(t *testing.T) {
	db, mock := newMockGormDB(t)
	repo := NewRebateRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `user_behavior_rebate_order`")).
		WillReturnError(&mysql.MySQLError{Number: 1062, Message: "duplicate"})
	mock.ExpectRollback()

	err := repo.SaveUserRebateRecords(context.Background(), []rebate.BehaviorRebateAggregate{
		{
			UserID: "xiaofuge",
			Order: rebate.BehaviorRebateOrderEntity{
				UserID:        "xiaofuge",
				OrderID:       "rebate-001",
				BehaviorType:  rebate.BehaviorTypeSign,
				RebateDesc:    "签到返利积分",
				RebateType:    rebate.RebateTypeIntegral,
				RebateConfig:  "10",
				OutBusinessNo: "20260525",
				BizID:         "xiaofuge_integral_20260525",
			},
			Task: rebate.TaskEntity{
				UserID:    "xiaofuge",
				Topic:     rebate.TopicSendRebate,
				MessageID: "msg-001",
				Message:   `{"bizId":"xiaofuge_integral_20260525"}`,
				State:     rebate.TaskStateCreate,
			},
		},
	})

	var appErr types.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected app error, got %v", err)
	}
	if appErr.Code != types.ResponseCodeIndexDup {
		t.Fatalf("expected duplicate code, got %s", appErr.Code.Code)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestRebateRepositorySaveUserRebateRecords(t *testing.T) {
	db, mock := newMockGormDB(t)
	repo := NewRebateRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `user_behavior_rebate_order`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `task`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.SaveUserRebateRecords(context.Background(), []rebate.BehaviorRebateAggregate{
		{
			UserID: "xiaofuge",
			Order: rebate.BehaviorRebateOrderEntity{
				UserID:        "xiaofuge",
				OrderID:       "rebate-001",
				BehaviorType:  rebate.BehaviorTypeSign,
				RebateDesc:    "签到返利积分",
				RebateType:    rebate.RebateTypeIntegral,
				RebateConfig:  "10",
				OutBusinessNo: "20260525",
				BizID:         "xiaofuge_integral_20260525",
			},
			Task: rebate.TaskEntity{
				UserID:    "xiaofuge",
				Topic:     rebate.TopicSendRebate,
				MessageID: "msg-001",
				Message:   `{"bizId":"xiaofuge_integral_20260525"}`,
				State:     rebate.TaskStateCreate,
			},
		},
	})
	if err != nil {
		t.Fatalf("save user rebate records: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestRebateRepositorySaveUserRebateRecordsEmpty(t *testing.T) {
	db, _ := newMockGormDB(t)
	repo := NewRebateRepository(db)

	if err := repo.SaveUserRebateRecords(context.Background(), nil); err != nil {
		t.Fatalf("save empty user rebate records: %v", err)
	}
}
