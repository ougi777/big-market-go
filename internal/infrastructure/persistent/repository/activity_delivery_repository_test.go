package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"bm-go/internal/domain/activity"
	"bm-go/internal/types"

	"gorm.io/gorm"
)

func TestActivityRepositoryDeliverActivityOrderNotFound(t *testing.T) {
	db, mock := newMockGormDB(t)
	repo := NewActivityRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT `user_id`,`activity_id`,`total_count`,`day_count`,`month_count`,`state` FROM `raffle_activity_order` WHERE user_id = ? and out_business_no = ? ORDER BY `raffle_activity_order`.`id` LIMIT ?")).
		WithArgs("xiaofuge", "biz-001", 1).
		WillReturnError(gorm.ErrRecordNotFound)
	mock.ExpectRollback()

	err := repo.DeliverActivityOrder(context.Background(), activity.DeliveryOrderEntity{
		UserID:        "xiaofuge",
		OutBusinessNo: "biz-001",
	})

	var appErr types.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected app error, got %v", err)
	}
	if appErr.Code != types.ResponseCodeIllegalParam {
		t.Fatalf("expected illegal param code, got %s", appErr.Code.Code)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}
