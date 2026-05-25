package service

import (
	"context"
	"testing"

	"bm-go/internal/domain/credit"
)

func TestAccountServiceQueryUserCreditAccount(t *testing.T) {
	repo := &fakeCreditAccountRepository{
		account: credit.AccountEntity{
			UserID:          "xiaofuge",
			AvailableAmount: 12.35,
		},
		exists: true,
	}
	service := NewAccountService(repo)

	amount, err := service.QueryUserCreditAccount(context.Background(), "xiaofuge")
	if err != nil {
		t.Fatalf("query user credit account: %v", err)
	}
	if amount != 12.35 {
		t.Fatalf("expected amount 12.35, got %.2f", amount)
	}
}

type fakeCreditAccountRepository struct {
	account credit.AccountEntity
	exists  bool
}

func (f *fakeCreditAccountRepository) QueryUserCreditAccount(ctx context.Context, userID string) (credit.AccountEntity, bool, error) {
	return f.account, f.exists, nil
}
