package repository

import (
	"context"
	"errors"

	"bm-go/internal/domain/activity"
	"bm-go/internal/infrastructure/persistent/po"

	"gorm.io/gorm"
)

func (r *ActivityRepository) QuerySkuProductListByActivityID(ctx context.Context, activityID int64) ([]activity.SkuProductEntity, error) {
	var skuPOList []po.RaffleActivitySku
	err := r.defaultDB(ctx).
		Select("sku", "activity_id", "activity_count_id", "stock_count", "stock_count_surplus", "product_amount").
		Where("activity_id = ?", activityID).
		Find(&skuPOList).
		Error
	if err != nil {
		return nil, err
	}

	products := make([]activity.SkuProductEntity, 0, len(skuPOList))
	for _, skuPO := range skuPOList {
		activityCount, err := r.queryActivityCount(ctx, skuPO.ActivityCountID)
		if err != nil {
			return nil, err
		}

		products = append(products, activity.SkuProductEntity{
			SKU:               skuPO.SKU,
			ActivityID:        skuPO.ActivityID,
			ActivityCountID:   skuPO.ActivityCountID,
			StockCount:        skuPO.StockCount,
			StockCountSurplus: skuPO.StockCountSurplus,
			ProductAmount:     skuPO.ProductAmount,
			ActivityCount:     activityCount,
		})
	}
	return products, nil
}

func (r *ActivityRepository) QuerySkuProductBySKU(ctx context.Context, sku int64) (activity.SkuProductEntity, bool, error) {
	var skuPO po.RaffleActivitySku
	err := r.defaultDB(ctx).
		Select("sku", "activity_id", "activity_count_id", "stock_count", "stock_count_surplus", "product_amount").
		Where("sku = ?", sku).
		First(&skuPO).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return activity.SkuProductEntity{}, false, nil
	}
	if err != nil {
		return activity.SkuProductEntity{}, false, err
	}
	activityCount, err := r.queryActivityCount(ctx, skuPO.ActivityCountID)
	if err != nil {
		return activity.SkuProductEntity{}, false, err
	}
	return activity.SkuProductEntity{
		SKU:               skuPO.SKU,
		ActivityID:        skuPO.ActivityID,
		ActivityCountID:   skuPO.ActivityCountID,
		StockCount:        skuPO.StockCount,
		StockCountSurplus: skuPO.StockCountSurplus,
		ProductAmount:     skuPO.ProductAmount,
		ActivityCount:     activityCount,
	}, true, nil
}

func (r *ActivityRepository) queryActivityCount(ctx context.Context, activityCountID int64) (activity.ActivityCountEntity, error) {
	var countPO po.RaffleActivityCount
	err := r.defaultDB(ctx).
		Select("activity_count_id", "total_count", "day_count", "month_count").
		Where("activity_count_id = ?", activityCountID).
		First(&countPO).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return activity.ActivityCountEntity{ActivityCountID: activityCountID}, nil
	}
	if err != nil {
		return activity.ActivityCountEntity{}, err
	}

	return activity.ActivityCountEntity{
		ActivityCountID: countPO.ActivityCountID,
		TotalCount:      countPO.TotalCount,
		DayCount:        countPO.DayCount,
		MonthCount:      countPO.MonthCount,
	}, nil
}
