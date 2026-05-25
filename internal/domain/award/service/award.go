package service

import (
	"context"

	"bm-go/internal/domain/award"
)

type AwardService struct {
	repo award.Repository
}

func NewAwardService(repo award.Repository) *AwardService {
	return &AwardService{repo: repo}
}

func (s *AwardService) SaveUserAwardRecord(ctx context.Context, record award.UserAwardRecordEntity) error {
	return s.repo.SaveUserAwardRecord(ctx, record)
}
