package service

import (
	"context"
	"testing"

	"bm-go/internal/domain/activity"
)

func TestCreditServiceQueryUserCreditAccount(t *testing.T) {
	repo := &fakeCreditAccountRepository{
		account: activity.CreditAccountEntity{
			UserID:          "xiaofuge",
			AvailableAmount: 12.35,
		},
		exists: true,
	}
	service := NewCreditService(repo)

	amount, err := service.QueryUserCreditAccount(context.Background(), "xiaofuge")
	if err != nil {
		t.Fatalf("query user credit account: %v", err)
	}
	if amount != 12.35 {
		t.Fatalf("expected amount 12.35, got %.2f", amount)
	}
}

type fakeCreditAccountRepository struct {
	account activity.CreditAccountEntity
	exists  bool
}

func (f *fakeCreditAccountRepository) QueryUserCreditAccount(ctx context.Context, userID string) (activity.CreditAccountEntity, bool, error) {
	return f.account, f.exists, nil
}
