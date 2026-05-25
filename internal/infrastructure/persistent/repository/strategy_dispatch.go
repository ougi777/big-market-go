package repository

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"

	"bm-go/internal/domain/strategy"
	"bm-go/internal/types"

	"github.com/redis/go-redis/v9"
)

type RateTableStore interface {
	Get(ctx context.Context, key string) (string, error)
	HGet(ctx context.Context, key string, field string) (string, error)
}

type redisRateTableStore struct {
	client *redis.Client
}

func (s *redisRateTableStore) Get(ctx context.Context, key string) (string, error) {
	return s.client.Get(ctx, key).Result()
}

func (s *redisRateTableStore) HGet(ctx context.Context, key string, field string) (string, error) {
	return s.client.HGet(ctx, key, field).Result()
}

type StrategyDispatch struct {
	store RateTableStore
}

var _ strategy.Dispatch = (*StrategyDispatch)(nil)

func NewStrategyDispatch(redisClient *redis.Client) *StrategyDispatch {
	return NewStrategyDispatchWithStore(&redisRateTableStore{client: redisClient})
}

func NewStrategyDispatchWithStore(store RateTableStore) *StrategyDispatch {
	return &StrategyDispatch{store: store}
}

func (d *StrategyDispatch) GetRandomAwardID(ctx context.Context, strategyID int64, ruleWeightValue ...string) (int, error) {
	key := strconv.FormatInt(strategyID, 10)
	if len(ruleWeightValue) > 0 && ruleWeightValue[0] != "" {
		key = key + types.Underline + ruleWeightValue[0]
	}
	return d.getRandomAwardID(ctx, key)
}

func (d *StrategyDispatch) getRandomAwardID(ctx context.Context, key string) (int, error) {
	if d.store == nil {
		return 0, errRepositoryNotImplemented
	}

	rateRangeValue, err := d.store.Get(ctx, types.RedisKeyStrategyRateRange+key)
	if err != nil {
		return 0, err
	}
	rateRange, err := strconv.Atoi(rateRangeValue)
	if err != nil {
		return 0, fmt.Errorf("parse strategy rate range %q: %w", rateRangeValue, err)
	}
	if rateRange <= 0 {
		return 0, fmt.Errorf("strategy rate range must be positive: %d", rateRange)
	}

	rateKey, err := secureRandomInt(rateRange)
	if err != nil {
		return 0, err
	}

	awardIDValue, err := d.store.HGet(ctx, types.RedisKeyStrategyRateTable+key, strconv.Itoa(rateKey))
	if err != nil {
		return 0, err
	}
	awardID, err := strconv.Atoi(awardIDValue)
	if err != nil {
		return 0, fmt.Errorf("parse strategy award id %q: %w", awardIDValue, err)
	}
	return awardID, nil
}

func secureRandomInt(max int) (int, error) {
	value, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0, err
	}
	return int(value.Int64()), nil
}
