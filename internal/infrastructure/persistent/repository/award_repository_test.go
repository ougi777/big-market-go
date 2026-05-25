package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"bm-go/internal/domain/award"
	"bm-go/internal/types"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
)

func TestAwardRepositoryQueryAwardConfig(t *testing.T) {
	db, mock := newMockGormDB(t)
	repo := NewAwardRepository(db)

	rows := sqlmock.NewRows([]string{"award_config"}).AddRow("10,100")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT `award_config` FROM `award` WHERE award_id = ? LIMIT ?")).
		WithArgs(101, 1).
		WillReturnRows(rows)

	config, err := repo.QueryAwardConfig(context.Background(), 101)
	if err != nil {
		t.Fatalf("query award config: %v", err)
	}
	if config != "10,100" {
		t.Fatalf("unexpected award config: %s", config)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestAwardRepositoryQueryAwardKey(t *testing.T) {
	db, mock := newMockGormDB(t)
	repo := NewAwardRepository(db)

	rows := sqlmock.NewRows([]string{"award_key"}).AddRow("user_credit_random")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT `award_key` FROM `award` WHERE award_id = ? LIMIT ?")).
		WithArgs(101, 1).
		WillReturnRows(rows)

	key, err := repo.QueryAwardKey(context.Background(), 101)
	if err != nil {
		t.Fatalf("query award key: %v", err)
	}
	if key != award.AwardKeyUserCreditRand {
		t.Fatalf("unexpected award key: %s", key)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestAwardRepositoryQueryNoSendMessageTaskList(t *testing.T) {
	db, mock := newMockGormDB(t)
	repo := NewAwardRepository(db)

	rows := sqlmock.NewRows([]string{"user_id", "topic", "message_id", "message", "state"}).
		AddRow("xiaofuge", award.TopicSendAward, "msg-001", `{"orderId":"order-001"}`, award.TaskStateCreate)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT `user_id`,`topic`,`message_id`,`message`,`state` FROM `task` WHERE state = ? or (state = ? and update_time < date_sub(now(), interval 6 second)) LIMIT ?")).
		WithArgs(award.TaskStateFail, award.TaskStateCreate, 10).
		WillReturnRows(rows)

	tasks, err := repo.QueryNoSendMessageTaskList(context.Background(), 0)
	if err != nil {
		t.Fatalf("query no send message task list: %v", err)
	}
	if len(tasks) != 1 || tasks[0].MessageID != "msg-001" || tasks[0].Topic != award.TopicSendAward {
		t.Fatalf("unexpected tasks: %+v", tasks)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestAwardRepositorySaveUserAwardRecordDuplicate(t *testing.T) {
	db, mock := newMockGormDB(t)
	repo := NewAwardRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `user_award_record`")).
		WillReturnError(&mysql.MySQLError{Number: 1062, Message: "duplicate"})
	mock.ExpectRollback()

	err := repo.SaveUserAwardRecord(context.Background(), award.UserAwardRecordEntity{
		UserID:     "xiaofuge",
		ActivityID: 100301,
		StrategyID: 100006,
		OrderID:    "order-001",
		AwardID:    101,
		AwardTitle: "积分",
		AwardTime:  time.Date(2026, 5, 25, 10, 0, 0, 0, time.Local),
		AwardState: award.AwardStateCreate,
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

func TestAwardRepositorySaveGiveOutPrizesAwardStateError(t *testing.T) {
	db, mock := newMockGormDB(t)
	repo := NewAwardRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `user_credit_account` SET `available_amount`=available_amount + ?,`total_amount`=total_amount + ?,`update_time`=? WHERE user_id = ?")).
		WithArgs(10.0, 10.0, sqlmock.AnyArg(), "xiaofuge").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `user_award_record` SET `award_state`=?,`update_time`=? WHERE user_id = ? and order_id = ? and award_state = ?")).
		WithArgs(award.AwardStateComplete, sqlmock.AnyArg(), "xiaofuge", "order-001", award.AwardStateCreate).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectRollback()

	err := repo.SaveGiveOutPrizes(context.Background(), award.GiveOutPrizesAggregate{
		UserID: "xiaofuge",
		UserAwardRecord: award.UserAwardRecordEntity{
			UserID:     "xiaofuge",
			OrderID:    "order-001",
			AwardState: award.AwardStateComplete,
		},
		UserCreditAward: award.UserCreditAwardEntity{
			UserID:       "xiaofuge",
			CreditAmount: 10,
		},
	})

	var appErr types.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected app error, got %v", err)
	}
	if appErr.Code != types.ResponseCodeActivityOrderStateError {
		t.Fatalf("expected order state code, got %s", appErr.Code.Code)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}
