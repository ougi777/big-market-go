package po

import "time"

type RaffleActivitySku struct {
	ID                int64     `gorm:"column:id;primaryKey"`
	SKU               int64     `gorm:"column:sku"`
	ActivityID        int64     `gorm:"column:activity_id"`
	ActivityCountID   int64     `gorm:"column:activity_count_id"`
	StockCount        int       `gorm:"column:stock_count"`
	StockCountSurplus int       `gorm:"column:stock_count_surplus"`
	ProductAmount     float64   `gorm:"column:product_amount"`
	CreateTime        time.Time `gorm:"column:create_time"`
	UpdateTime        time.Time `gorm:"column:update_time"`
}

func (RaffleActivitySku) TableName() string {
	return "raffle_activity_sku"
}
