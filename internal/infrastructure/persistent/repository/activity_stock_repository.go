package repository

import (
	"context"
	"time"

	"bm-go/internal/infrastructure/persistent/po"

	"gorm.io/gorm"
)

func (r *ActivityRepository) UpdateActivitySkuStock(ctx context.Context, sku int64) error {
	return r.defaultDB(ctx).
		Model(&po.RaffleActivitySku{}).
		Where("sku = ? and stock_count_surplus > 0", sku).
		Updates(map[string]any{
			"stock_count_surplus": gorm.Expr("stock_count_surplus - ?", 1),
			"update_time":         time.Now(),
		}).
		Error
}

func (r *ActivityRepository) ClearActivitySkuStock(ctx context.Context, sku int64) error {
	return r.defaultDB(ctx).
		Model(&po.RaffleActivitySku{}).
		Where("sku = ?", sku).
		Updates(map[string]any{
			"stock_count_surplus": 0,
			"update_time":         time.Now(),
		}).
		Error
}
