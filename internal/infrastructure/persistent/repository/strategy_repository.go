package repository

import (
	"context"

	"bm-go/internal/domain/strategy"
	treepkg "bm-go/internal/domain/strategy/rule/tree"
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
