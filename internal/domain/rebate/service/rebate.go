package service

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"

	"bm-go/internal/domain/rebate"
	"bm-go/internal/types"
)

type RebateService struct {
	repo      rebate.Repository
	publisher rebate.MessagePublisher
	now       func() time.Time
	newID     func(int) (string, error)
}

func NewRebateService(repo rebate.Repository, publisher rebate.MessagePublisher) *RebateService {
	return &RebateService{
		repo:      repo,
		publisher: publisher,
		now:       time.Now,
		newID:     randomNumeric,
	}
}

func (s *RebateService) CalendarSignRebate(ctx context.Context, userID string) (bool, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return false, types.NewAppError(types.ResponseCodeIllegalParam, nil)
	}

	outBusinessNo := s.now().Format("20060102")
	orders, err := s.createOrder(ctx, userID, rebate.BehaviorTypeSign, outBusinessNo)
	if err != nil {
		return false, err
	}
	return len(orders) > 0, nil
}

func (s *RebateService) IsCalendarSignRebate(ctx context.Context, userID string) (bool, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return false, types.NewAppError(types.ResponseCodeIllegalParam, nil)
	}

	orders, err := s.repo.QueryOrderByOutBusinessNo(ctx, userID, s.now().Format("20060102"))
	if err != nil {
		return false, err
	}
	return len(orders) > 0, nil
}

func (s *RebateService) createOrder(ctx context.Context, userID string, behaviorType string, outBusinessNo string) ([]string, error) {
	configs, err := s.repo.QueryDailyBehaviorRebateConfig(ctx, behaviorType)
	if err != nil {
		return nil, err
	}
	if len(configs) == 0 {
		return []string{}, nil
	}

	aggregates := make([]rebate.BehaviorRebateAggregate, 0, len(configs))
	orderIDs := make([]string, 0, len(configs))
	for _, config := range configs {
		orderID, err := s.newID(12)
		if err != nil {
			return nil, err
		}
		messageID, err := s.newID(11)
		if err != nil {
			return nil, err
		}
		bizID := fmt.Sprintf("%s_%s_%s", userID, config.RebateType, outBusinessNo)
		message, err := json.Marshal(rebate.EventMessage[rebate.SendRebateMessage]{
			ID:        messageID,
			Timestamp: s.now().UnixMilli(),
			Data: rebate.SendRebateMessage{
				UserID:       userID,
				RebateDesc:   config.RebateDesc,
				RebateType:   config.RebateType,
				RebateConfig: config.RebateConfig,
				BizID:        bizID,
			},
		})
		if err != nil {
			return nil, err
		}

		orderIDs = append(orderIDs, orderID)
		aggregates = append(aggregates, rebate.BehaviorRebateAggregate{
			UserID: userID,
			Order: rebate.BehaviorRebateOrderEntity{
				UserID:        userID,
				OrderID:       orderID,
				BehaviorType:  config.BehaviorType,
				RebateDesc:    config.RebateDesc,
				RebateType:    config.RebateType,
				RebateConfig:  config.RebateConfig,
				OutBusinessNo: outBusinessNo,
				BizID:         bizID,
			},
			Task: rebate.TaskEntity{
				UserID:    userID,
				Topic:     rebate.TopicSendRebate,
				MessageID: messageID,
				Message:   string(message),
				State:     rebate.TaskStateCreate,
			},
		})
	}

	if err := s.repo.SaveUserRebateRecords(ctx, aggregates); err != nil {
		return nil, err
	}

	for _, aggregate := range aggregates {
		task := aggregate.Task
		if err := s.publisher.Publish(ctx, task.Topic, task.Message); err != nil {
			_ = s.repo.UpdateTaskSendMessageFail(ctx, task.UserID, task.MessageID)
			continue
		}
		_ = s.repo.UpdateTaskSendMessageCompleted(ctx, task.UserID, task.MessageID)
	}
	return orderIDs, nil
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
