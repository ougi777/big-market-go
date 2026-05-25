package redis

import (
	"context"

	"bm-go/internal/domain/activity"

	goredis "github.com/redis/go-redis/v9"
)

type ActivityStore struct {
	client *goredis.Client
}

var _ activity.SkuStockStore = (*ActivityStore)(nil)

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
