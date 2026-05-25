package redis

import (
	"context"
	"encoding/json"

	"bm-go/internal/domain/activity"
	"bm-go/internal/types"

	goredis "github.com/redis/go-redis/v9"
)

type ActivityStore struct {
	client *goredis.Client
}

var _ activity.SkuStockStore = (*ActivityStore)(nil)
var _ activity.SkuStockQueue = (*ActivityStore)(nil)

func NewActivityStore(client *goredis.Client) *ActivityStore {
	return &ActivityStore{client: client}
}

func (s *ActivityStore) ensureClient() error {
	if s == nil || s.client == nil {
		return ErrClientNotConnected
	}
	return nil
}

func (s *ActivityStore) CacheActivitySkuStockCount(ctx context.Context, key string, stockCount int) error {
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
	return s.client.Set(ctx, key, stockCount, 0).Err()
}

func (s *ActivityStore) SubtractActivitySkuStock(ctx context.Context, key string) (int64, error) {
	if err := s.ensureClient(); err != nil {
		return 0, err
	}

	return s.client.Decr(ctx, key).Result()
}

func (s *ActivityStore) SendActivitySkuStockConsumeQueue(ctx context.Context, stockKey activity.ActivitySkuStockKey) error {
	if err := s.ensureClient(); err != nil {
		return err
	}

	value, err := json.Marshal(stockKey)
	if err != nil {
		return err
	}
	return s.client.RPush(ctx, types.RedisKeyActivitySkuStockQueue, string(value)).Err()
}

func (s *ActivityStore) TakeActivitySkuStock(ctx context.Context) (activity.ActivitySkuStockKey, bool, error) {
	if err := s.ensureClient(); err != nil {
		return activity.ActivitySkuStockKey{}, false, err
	}

	value, err := s.client.LPop(ctx, types.RedisKeyActivitySkuStockQueue).Result()
	if err == goredis.Nil {
		return activity.ActivitySkuStockKey{}, false, nil
	}
	if err != nil {
		return activity.ActivitySkuStockKey{}, false, err
	}

	var key activity.ActivitySkuStockKey
	if err := json.Unmarshal([]byte(value), &key); err != nil {
		return activity.ActivitySkuStockKey{}, false, err
	}
	return key, true, nil
}

func (s *ActivityStore) ClearActivitySkuStockQueue(ctx context.Context) error {
	if err := s.ensureClient(); err != nil {
		return err
	}

	return s.client.Del(ctx, types.RedisKeyActivitySkuStockQueue).Err()
}
