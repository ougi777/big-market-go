package repository

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strconv"

	"bm-go/internal/domain/strategy"
	"bm-go/internal/domain/strategy/rule/tree"
	"bm-go/internal/types"

	"github.com/redis/go-redis/v9"
)

type RateTableStore interface {
	Get(ctx context.Context, key string) (string, error)
	HGet(ctx context.Context, key string, field string) (string, error)
	Decr(ctx context.Context, key string) (int64, error)
	Set(ctx context.Context, key string, value string) error
	SetNX(ctx context.Context, key string, value string) (bool, error)
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

func (s *redisRateTableStore) Decr(ctx context.Context, key string) (int64, error) {
	return s.client.Decr(ctx, key).Result()
}

func (s *redisRateTableStore) Set(ctx context.Context, key string, value string) error {
	return s.client.Set(ctx, key, value, 0).Err()
}

func (s *redisRateTableStore) SetNX(ctx context.Context, key string, value string) (bool, error) {
	return s.client.SetNX(ctx, key, value, 0).Result()
}

type StrategyDispatch struct {
	store RateTableStore
}

var _ strategy.Dispatch = (*StrategyDispatch)(nil)
var _ tree.StockDispatch = (*StrategyDispatch)(nil)

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

func (d *StrategyDispatch) SubtractionAwardStock(ctx context.Context, strategyID int64, awardID int) (bool, error) {
	if d.store == nil {
		return false, types.NewAppError(types.ResponseCodeUnassembledStrategy, nil)
	}

	cacheKey := types.RedisKeyStrategyAwardCount + strconv.FormatInt(strategyID, 10) + types.Underline + strconv.Itoa(awardID)
	surplus, err := d.store.Decr(ctx, cacheKey)
	if err != nil {
		return false, err
	}
	if surplus < 0 {
		if err := d.store.Set(ctx, cacheKey, "0"); err != nil {
			return false, err
		}
		return false, nil
	}

	lockKey := cacheKey + types.Underline + strconv.FormatInt(surplus, 10)
	return d.store.SetNX(ctx, lockKey, "lock")
}

func (d *StrategyDispatch) getRandomAwardID(ctx context.Context, key string) (int, error) {
	if d.store == nil {
		return 0, types.NewAppError(types.ResponseCodeUnassembledStrategy, nil)
	}

	rateRangeValue, err := d.store.Get(ctx, types.RedisKeyStrategyRateRange+key)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, types.NewAppError(types.ResponseCodeUnassembledStrategy, err)
		}
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
		if errors.Is(err, redis.Nil) {
			return 0, types.NewAppError(types.ResponseCodeUnassembledStrategy, err)
		}
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
