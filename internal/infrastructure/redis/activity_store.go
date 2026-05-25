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

func (s *ActivityStore) CacheActivitySkuStockCount(ctx context.Context, key string, stockCount int) error {
	exists, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return err
	}
	if exists > 0 {
		return nil
	}
	return s.client.Set(ctx, key, stockCount, 0).Err()
}

func (s *ActivityStore) TakeActivitySkuStock(ctx context.Context) (activity.ActivitySkuStockKey, bool, error) {
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
	return s.client.Del(ctx, types.RedisKeyActivitySkuStockQueue).Err()
}
