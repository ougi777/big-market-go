package service

import (
	"context"
	"strings"

	"bm-go/internal/domain/credit"
	"bm-go/internal/types"
)

type AccountService struct {
	repo credit.AccountRepository
}

func NewAccountService(repo credit.AccountRepository) *AccountService {
	return &AccountService{repo: repo}
}

func (s *AccountService) QueryUserCreditAccount(ctx context.Context, userID string) (float64, error) {
	if strings.TrimSpace(userID) == "" {
		return 0, types.NewAppError(types.ResponseCodeIllegalParam, nil)
	}
	account, exists, err := s.repo.QueryUserCreditAccount(ctx, userID)
	if err != nil {
		return 0, err
	}
	if !exists {
		return 0, nil
	}
	return account.AvailableAmount, nil
}
