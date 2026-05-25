package redis

import (
	"context"
	"strconv"

	"bm-go/internal/domain/strategy"
	"bm-go/internal/types"

	goredis "github.com/redis/go-redis/v9"
)

type StrategyStore struct {
	client *goredis.Client
}

var _ strategy.RateTableStore = (*StrategyStore)(nil)

func NewStrategyStore(client *goredis.Client) *StrategyStore {
	return &StrategyStore{client: client}
}

func (s *StrategyStore) StoreStrategyAwardSearchRateTable(ctx context.Context, key string, rateRange int, table map[int]int) error {
	if err := s.client.Set(ctx, types.RedisKeyStrategyRateRange+key, rateRange, 0).Err(); err != nil {
		return err
	}

	values := make(map[string]interface{}, len(table))
	for rateKey, awardID := range table {
		values[strconv.Itoa(rateKey)] = awardID
	}
	return s.client.HSet(ctx, types.RedisKeyStrategyRateTable+key, values).Err()
}

func (s *StrategyStore) CacheStrategyAwardCount(ctx context.Context, key string, awardCount int) error {
	exists, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return err
	}
	if exists > 0 {
		return nil
	}
	return s.client.Set(ctx, key, awardCount, 0).Err()
}
