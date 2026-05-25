package repository

import (
	"context"
	"errors"
	"time"

	"bm-go/internal/domain/credit"
	"bm-go/internal/infrastructure/persistent/po"
	"bm-go/internal/infrastructure/persistent/sharding"
	"bm-go/internal/types"

	"gorm.io/gorm"
)

type CreditRepository struct {
	db      dbRouter
	sharder sharding.Router
}

var _ credit.AccountRepository = (*CreditRepository)(nil)
var _ credit.TradeRepository = (*CreditRepository)(nil)

func NewCreditRepository(db *gorm.DB, routers ...sharding.Router) *CreditRepository {
	return NewCreditRepositoryWithDBRouter(singleDBRouter{db: db}, routers...)
}

func NewCreditRepositoryWithDBRouter(db dbRouter, routers ...sharding.Router) *CreditRepository {
	router := sharding.NewRouter(1)
	if len(routers) > 0 {
		router = routers[0]
	}
	return &CreditRepository{db: db, sharder: router}
}

func (r *CreditRepository) shardDB(ctx context.Context, userID string) *gorm.DB {
	return r.db.Shard(r.sharder.DBKey(userID)).WithContext(ctx)
}

func (r *CreditRepository) QueryUserCreditAccount(ctx context.Context, userID string) (credit.AccountEntity, bool, error) {
	var accountPO po.UserCreditAccount
	err := r.shardDB(ctx, userID).
		Select("user_id", "available_amount").
		Where("user_id = ?", userID).
		First(&accountPO).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return credit.AccountEntity{}, false, nil
	}
	if err != nil {
		return credit.AccountEntity{}, false, err
	}
	return credit.AccountEntity{
		UserID:          accountPO.UserID,
		AvailableAmount: accountPO.AvailableAmount,
	}, true, nil
}

func (r *CreditRepository) CompleteCreditPayOrder(ctx context.Context, aggregate credit.CompleteSkuExchangeAggregate) error {
	now := time.Now()
	return r.shardDB(ctx, aggregate.UserID).Transaction(func(tx *gorm.DB) error {
		if err := adjustUserCreditAccount(tx, aggregate.CreditOrder); err != nil {
			return err
		}
		if err := tx.Table(r.sharder.Table("user_credit_order", aggregate.UserID)).Create(&po.UserCreditOrder{
			UserID:        aggregate.CreditOrder.UserID,
			OrderID:       aggregate.CreditOrder.OrderID,
			TradeName:     aggregate.CreditOrder.TradeName,
			TradeType:     aggregate.CreditOrder.TradeType,
			TradeAmount:   aggregate.CreditOrder.TradeAmount,
			OutBusinessNo: aggregate.CreditOrder.OutBusinessNo,
			CreateTime:    now,
			UpdateTime:    now,
		}).Error; err != nil {
			return types.NewAppError(types.ResponseCodeIndexDup, err)
		}
		if aggregate.SendTask.MessageID != "" {
			if err := tx.Create(&po.Task{
				UserID:     aggregate.SendTask.UserID,
				Topic:      aggregate.SendTask.Topic,
				MessageID:  aggregate.SendTask.MessageID,
				Message:    aggregate.SendTask.Message,
				State:      aggregate.SendTask.State,
				CreateTime: now,
				UpdateTime: now,
			}).Error; err != nil {
				return types.NewAppError(types.ResponseCodeIndexDup, err)
			}
		}

		return nil
	})
}

func (r *CreditRepository) SaveRebateIntegralOrder(ctx context.Context, rebateIntegral credit.RebateIntegralEntity) error {
	now := time.Now()
	return r.shardDB(ctx, rebateIntegral.UserID).Transaction(func(tx *gorm.DB) error {
		creditOrder := credit.OrderEntity{
			UserID:        rebateIntegral.UserID,
			OrderID:       rebateIntegral.OrderID,
			TradeName:     "REBATE",
			TradeType:     "forward",
			TradeAmount:   rebateIntegral.TradeAmount,
			OutBusinessNo: rebateIntegral.OutBusinessNo,
		}
		if err := adjustOrCreateUserCreditAccount(tx, creditOrder, now); err != nil {
			return err
		}
		if err := tx.Table(r.sharder.Table("user_credit_order", rebateIntegral.UserID)).Create(&po.UserCreditOrder{
			UserID:        creditOrder.UserID,
			OrderID:       creditOrder.OrderID,
			TradeName:     creditOrder.TradeName,
			TradeType:     creditOrder.TradeType,
			TradeAmount:   creditOrder.TradeAmount,
			OutBusinessNo: creditOrder.OutBusinessNo,
			CreateTime:    now,
			UpdateTime:    now,
		}).Error; err != nil {
			return types.NewAppError(types.ResponseCodeIndexDup, err)
		}
		return nil
	})
}

func adjustUserCreditAccount(tx *gorm.DB, creditOrder credit.OrderEntity) error {
	now := time.Now()
	result := tx.Model(&po.UserCreditAccount{}).
		Where("user_id = ? and available_amount + ? >= 0", creditOrder.UserID, creditOrder.TradeAmount).
		Updates(map[string]any{
			"total_amount":     gorm.Expr("total_amount + ?", creditOrder.TradeAmount),
			"available_amount": gorm.Expr("available_amount + ?", creditOrder.TradeAmount),
			"update_time":      now,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return types.NewAppError(types.ResponseCodeAccountQuotaError, nil)
	}
	return nil
}

func adjustOrCreateUserCreditAccount(tx *gorm.DB, creditOrder credit.OrderEntity, now time.Time) error {
	result := tx.Model(&po.UserCreditAccount{}).
		Where("user_id = ? and available_amount + ? >= 0", creditOrder.UserID, creditOrder.TradeAmount).
		Updates(map[string]any{
			"total_amount":     gorm.Expr("total_amount + ?", creditOrder.TradeAmount),
			"available_amount": gorm.Expr("available_amount + ?", creditOrder.TradeAmount),
			"update_time":      now,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		if creditOrder.TradeAmount < 0 {
			return types.NewAppError(types.ResponseCodeAccountQuotaError, nil)
		}
		if err := tx.Create(&po.UserCreditAccount{
			UserID:          creditOrder.UserID,
			TotalAmount:     creditOrder.TradeAmount,
			AvailableAmount: creditOrder.TradeAmount,
			AccountStatus:   "open",
			CreateTime:      now,
			UpdateTime:      now,
		}).Error; err != nil {
			return types.NewAppError(types.ResponseCodeIndexDup, err)
		}
	}
	return nil
}
