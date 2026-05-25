package redis

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"bm-go/internal/domain/strategy"
	"bm-go/internal/types"

	goredis "github.com/redis/go-redis/v9"
)

type StrategyStore struct {
	client *goredis.Client
}

var _ strategy.RateTableStore = (*StrategyStore)(nil)
var _ strategy.StockQueue = (*StrategyStore)(nil)

func NewStrategyStore(client *goredis.Client) *StrategyStore {
	return &StrategyStore{client: client}
}

func (s *StrategyStore) ensureClient() error {
	if s == nil || s.client == nil {
		return ErrClientNotConnected
	}
	return nil
}

func (s *StrategyStore) StoreStrategyAwardSearchRateTable(ctx context.Context, key string, rateRange int, table map[int]int) error {
	if err := s.ensureClient(); err != nil {
		return err
	}

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
	if err := s.ensureClient(); err != nil {
		return err
	}

	exists, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return err
	}
	if exists > 0 {
		return nil
	}
	return s.client.Set(ctx, key, awardCount, 0).Err()
}

func (s *StrategyStore) AwardStockConsumeSendQueue(ctx context.Context, strategyID int64, awardID int) error {
	if err := s.ensureClient(); err != nil {
		return err
	}

	value := fmt.Sprintf("%d:%d", strategyID, awardID)
	return s.client.RPush(ctx, types.RedisKeyStrategyAwardCountQueue, value).Err()
}

func (s *StrategyStore) TakeQueueValue(ctx context.Context) (strategy.AwardStockKey, bool, error) {
	if err := s.ensureClient(); err != nil {
		return strategy.AwardStockKey{}, false, err
	}

	value, err := s.client.LPop(ctx, types.RedisKeyStrategyAwardCountQueue).Result()
	if err == goredis.Nil {
		return strategy.AwardStockKey{}, false, nil
	}
	if err != nil {
		return strategy.AwardStockKey{}, false, err
	}
	key, err := parseAwardStockQueueValue(value)
	if err != nil {
		return strategy.AwardStockKey{}, false, err
	}
	return key, true, nil
}

func parseAwardStockQueueValue(value string) (strategy.AwardStockKey, error) {
	parts := strings.SplitN(value, ":", 2)
	if len(parts) != 2 {
		return strategy.AwardStockKey{}, fmt.Errorf("invalid award stock queue value: %s", value)
	}
	strategyID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return strategy.AwardStockKey{}, err
	}
	awardID, err := strconv.Atoi(parts[1])
	if err != nil {
		return strategy.AwardStockKey{}, err
	}
	return strategy.AwardStockKey{StrategyID: strategyID, AwardID: awardID}, nil
}
