package service

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"bm-go/internal/domain/award"
	"bm-go/internal/types"
)

type AwardService struct {
	repo             award.Repository
	taskRepo         award.TaskRepository
	publisher        award.MessagePublisher
	distributors     map[string]AwardDistributor
	now              func() time.Time
	messageGenerator func() (string, error)
	creditGenerator  func(min float64, max float64) (float64, error)
}

func NewAwardService(repo award.Repository, taskRepo award.TaskRepository, publisher award.MessagePublisher) *AwardService {
	service := &AwardService{
		repo:             repo,
		taskRepo:         taskRepo,
		publisher:        publisher,
		distributors:     make(map[string]AwardDistributor),
		now:              time.Now,
		messageGenerator: func() (string, error) { return randomNumeric(11) },
		creditGenerator:  randomCredit,
	}
	service.distributors[award.AwardKeyUserCreditRand] = &UserCreditRandomDistributor{service: service}
	return service
}

func (s *AwardService) SaveUserAwardRecord(ctx context.Context, record award.UserAwardRecordEntity) error {
	messageID, err := s.messageGenerator()
	if err != nil {
		return err
	}
	message, err := json.Marshal(award.EventMessage[award.SendAwardMessage]{
		ID:        messageID,
		Timestamp: s.now().UnixMilli(),
		Data: award.SendAwardMessage{
			UserID:      record.UserID,
			OrderID:     record.OrderID,
			AwardID:     record.AwardID,
			AwardTitle:  record.AwardTitle,
			AwardConfig: record.AwardConfig,
		},
	})
	if err != nil {
		return err
	}

	record.SendTask = award.TaskEntity{
		UserID:    record.UserID,
		Topic:     award.TopicSendAward,
		MessageID: messageID,
		Message:   string(message),
		State:     award.TaskStateCreate,
	}

	if err := s.repo.SaveUserAwardRecord(ctx, record); err != nil {
		return err
	}
	if s.publisher == nil {
		return nil
	}
	if err := s.publisher.Publish(ctx, record.SendTask.Topic, record.SendTask.Message); err != nil {
		if s.taskRepo != nil {
			_ = s.taskRepo.UpdateTaskSendMessageFail(ctx, record.SendTask.UserID, record.SendTask.MessageID)
		}
		return err
	}
	if s.taskRepo != nil {
		return s.taskRepo.UpdateTaskSendMessageCompleted(ctx, record.SendTask.UserID, record.SendTask.MessageID)
	}
	return nil
}

type AwardDistributor interface {
	GiveOutPrizes(ctx context.Context, distribute award.DistributeAwardEntity) error
}

func (s *AwardService) DistributeAward(ctx context.Context, distribute award.DistributeAwardEntity) error {
	awardKey, err := s.repo.QueryAwardKey(ctx, distribute.AwardID)
	if err != nil {
		return err
	}
	distributor := s.distributors[awardKey]
	if distributor == nil {
		return types.NewAppError(types.ResponseCodeIllegalParam, fmt.Errorf("award distributor missing: %s", awardKey))
	}
	return distributor.GiveOutPrizes(ctx, distribute)
}

type UserCreditRandomDistributor struct {
	service *AwardService
}

func (d *UserCreditRandomDistributor) GiveOutPrizes(ctx context.Context, distribute award.DistributeAwardEntity) error {
	awardConfig := strings.TrimSpace(distribute.AwardConfig)
	if awardConfig == "" {
		config, err := d.service.repo.QueryAwardConfig(ctx, distribute.AwardID)
		if err != nil {
			return err
		}
		awardConfig = strings.TrimSpace(config)
	}

	parts := strings.Split(awardConfig, ",")
	if len(parts) != 2 {
		return types.NewAppError(types.ResponseCodeIllegalParam, fmt.Errorf("invalid award config: %s", awardConfig))
	}
	min, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return types.NewAppError(types.ResponseCodeIllegalParam, err)
	}
	max, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return types.NewAppError(types.ResponseCodeIllegalParam, err)
	}
	creditAmount, err := d.service.creditGenerator(min, max)
	if err != nil {
		return err
	}

	return d.service.repo.SaveGiveOutPrizes(ctx, award.GiveOutPrizesAggregate{
		UserID: distribute.UserID,
		UserAwardRecord: award.UserAwardRecordEntity{
			UserID:     distribute.UserID,
			OrderID:    distribute.OrderID,
			AwardID:    distribute.AwardID,
			AwardState: award.AwardStateComplete,
		},
		UserCreditAward: award.UserCreditAwardEntity{
			UserID:       distribute.UserID,
			CreditAmount: creditAmount,
		},
	})
}

func randomNumeric(length int) (string, error) {
	value := make([]byte, length)
	for i := range value {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", fmt.Errorf("generate random numeric: %w", err)
		}
		value[i] = byte('0' + n.Int64())
	}
	return string(value), nil
}

func randomCredit(min float64, max float64) (float64, error) {
	if min > max {
		return 0, types.NewAppError(types.ResponseCodeIllegalParam, fmt.Errorf("invalid credit range: %.2f,%.2f", min, max))
	}
	if min == max {
		return min, nil
	}
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return 0, fmt.Errorf("generate random credit: %w", err)
	}
	rate := float64(n.Int64()) / 1000000
	return roundCredit(min + rate*(max-min)), nil
}

func roundCredit(value float64) float64 {
	return float64(int64(value*100+0.5)) / 100
}
