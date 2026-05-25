package repository

import (
	"context"
	"errors"

	"bm-go/internal/domain/activity"
	"bm-go/internal/infrastructure/persistent/po"
	"bm-go/internal/infrastructure/persistent/sharding"

	"gorm.io/gorm"
)

type ActivityRepository struct {
	db      dbRouter
	sharder sharding.Router
}

var _ activity.Repository = (*ActivityRepository)(nil)
var _ activity.AccountRepository = (*ActivityRepository)(nil)
var _ activity.SkuProductRepository = (*ActivityRepository)(nil)
var _ activity.SkuStockRepository = (*ActivityRepository)(nil)
var _ activity.PartakeRepository = (*ActivityRepository)(nil)
var _ activity.RebateRepository = (*ActivityRepository)(nil)
var _ activity.DeliveryRepository = (*ActivityRepository)(nil)

func NewActivityRepository(db *gorm.DB, routers ...sharding.Router) *ActivityRepository {
	return NewActivityRepositoryWithDBRouter(singleDBRouter{db: db}, routers...)
}

func NewActivityRepositoryWithDBRouter(db dbRouter, routers ...sharding.Router) *ActivityRepository {
	router := sharding.NewRouter(1)
	if len(routers) > 0 {
		router = routers[0]
	}
	return &ActivityRepository{db: db, sharder: router}
}

func (r *ActivityRepository) defaultDB(ctx context.Context) *gorm.DB {
	return r.db.Default().WithContext(ctx)
}

func (r *ActivityRepository) shardDB(ctx context.Context, userID string) *gorm.DB {
	return r.db.Shard(r.sharder.DBKey(userID)).WithContext(ctx)
}

func (r *ActivityRepository) QueryActivityByActivityID(ctx context.Context, activityID int64) (activity.ActivityEntity, bool, error) {
	var activityPO po.RaffleActivity
	err := r.defaultDB(ctx).
		Select("activity_id", "activity_name", "activity_desc", "begin_date_time", "end_date_time", "strategy_id", "state").
		Where("activity_id = ?", activityID).
		First(&activityPO).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return activity.ActivityEntity{}, false, nil
	}
	if err != nil {
		return activity.ActivityEntity{}, false, err
	}

	return activity.ActivityEntity{
		ActivityID:    activityPO.ActivityID,
		ActivityName:  activityPO.ActivityName,
		ActivityDesc:  activityPO.ActivityDesc,
		BeginDateTime: activityPO.BeginDateTime,
		EndDateTime:   activityPO.EndDateTime,
		StrategyID:    activityPO.StrategyID,
		State:         activityPO.State,
	}, true, nil
}
