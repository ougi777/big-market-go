package repository

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestSetTaskState(t *testing.T) {
	db, mock := newMockGormDB(t)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `task` SET `state`=?,`update_time`=? WHERE user_id = ? and message_id = ?")).
		WithArgs("completed", sqlmock.AnyArg(), "xiaofuge", "12345678901").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	if err := setTaskState(context.Background(), db, "xiaofuge", "12345678901", "completed"); err != nil {
		t.Fatalf("set task state: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}
