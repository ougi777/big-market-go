package repository

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
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
