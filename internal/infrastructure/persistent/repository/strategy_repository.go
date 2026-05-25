package repository

import (
	"context"
	"errors"

	"bm-go/internal/domain/strategy"
	treepkg "bm-go/internal/domain/strategy/rule/tree"
	"bm-go/internal/infrastructure/persistent/po"
	"bm-go/internal/infrastructure/persistent/sharding"

	"gorm.io/gorm"
)

type StrategyRepository struct {
	db         dbRouter
	sharder    sharding.Router
	stockQueue AwardStockQueue
}

type AwardStockQueue interface {
	AwardStockConsumeSendQueue(ctx context.Context, strategyID int64, awardID int) error
}

var _ strategy.Repository = (*StrategyRepository)(nil)
var _ strategy.ArmoryRepository = (*StrategyRepository)(nil)
var _ strategy.QueryRepository = (*StrategyRepository)(nil)
var _ strategy.RaffleRepository = (*StrategyRepository)(nil)
var _ strategy.StockRepository = (*StrategyRepository)(nil)
var _ treepkg.Repository = (*StrategyRepository)(nil)

func NewStrategyRepository(db *gorm.DB, stockQueues ...AwardStockQueue) *StrategyRepository {
	return NewStrategyRepositoryWithDBRouter(singleDBRouter{db: db}, sharding.NewRouter(1), stockQueues...)
}

func NewStrategyRepositoryWithDBRouter(db dbRouter, sharder sharding.Router, stockQueues ...AwardStockQueue) *StrategyRepository {
	var stockQueue AwardStockQueue
	if len(stockQueues) > 0 {
		stockQueue = stockQueues[0]
	}
	return &StrategyRepository{db: db, sharder: sharder, stockQueue: stockQueue}
}

func (r *StrategyRepository) defaultDB(ctx context.Context) *gorm.DB {
	return r.db.Default().WithContext(ctx)
}

func (r *StrategyRepository) shardDB(ctx context.Context, userID string) *gorm.DB {
	return r.db.Shard(r.sharder.DBKey(userID)).WithContext(ctx)
}

func (r *StrategyRepository) QueryStrategyEntityByStrategyID(ctx context.Context, strategyID int64) (strategy.StrategyEntity, error) {
	var strategyPO po.Strategy
	err := r.defaultDB(ctx).
		Select("strategy_id", "strategy_desc", "rule_models").
		Where("strategy_id = ?", strategyID).
		First(&strategyPO).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return strategy.StrategyEntity{}, nil
	}
	if err != nil {
		return strategy.StrategyEntity{}, err
	}

	return strategy.StrategyEntity{
		StrategyID:   strategyPO.StrategyID,
		StrategyDesc: strategyPO.StrategyDesc,
		RuleModel:    strategyPO.RuleModels,
	}, nil
}

func (r *StrategyRepository) QueryStrategyIDByActivityID(ctx context.Context, activityID int64) (int64, error) {
	var activityPO po.RaffleActivity
	err := r.defaultDB(ctx).
		Select("strategy_id").
		Where("activity_id = ?", activityID).
		First(&activityPO).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return activityPO.StrategyID, nil
}
