package service

import (
	"context"
	"errors"
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

func TestAccountServiceQueryUserCreditAccountIllegalParam(t *testing.T) {
	service := NewAccountService(&fakeCreditAccountRepository{})

	_, err := service.QueryUserCreditAccount(context.Background(), " ")
	if err == nil {
		t.Fatal("expected illegal param error")
	}
}

func TestAccountServiceQueryUserCreditAccountNotExists(t *testing.T) {
	service := NewAccountService(&fakeCreditAccountRepository{})

	amount, err := service.QueryUserCreditAccount(context.Background(), "xiaofuge")
	if err != nil {
		t.Fatalf("query user credit account: %v", err)
	}
	if amount != 0 {
		t.Fatalf("expected amount 0, got %.2f", amount)
	}
}

func TestAccountServiceQueryUserCreditAccountRepositoryError(t *testing.T) {
	service := NewAccountService(&fakeCreditAccountRepository{err: errors.New("query failed")})

	_, err := service.QueryUserCreditAccount(context.Background(), "xiaofuge")
	if err == nil {
		t.Fatal("expected repository error")
	}
}

type fakeCreditAccountRepository struct {
	account credit.AccountEntity
	exists  bool
	err     error
}

func (f *fakeCreditAccountRepository) QueryUserCreditAccount(ctx context.Context, userID string) (credit.AccountEntity, bool, error) {
	return f.account, f.exists, f.err
}
