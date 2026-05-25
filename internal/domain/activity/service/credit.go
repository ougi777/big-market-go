package service

import (
	"context"
	"strings"

	"bm-go/internal/domain/activity"
	"bm-go/internal/types"
)

type CreditService struct {
	repo activity.CreditAccountRepository
}

func NewCreditService(repo activity.CreditAccountRepository) *CreditService {
	return &CreditService{repo: repo}
}

func (s *CreditService) QueryUserCreditAccount(ctx context.Context, userID string) (float64, error) {
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
