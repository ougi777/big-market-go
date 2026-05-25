package activity

type AccountEntity struct {
	UserID            string
	ActivityID        int64
	TotalCount        int
	TotalCountSurplus int
	DayCount          int
	DayCountSurplus   int
	MonthCount        int
	MonthCountSurplus int
}

type AccountDayEntity struct {
	UserID          string
	ActivityID      int64
	Day             string
	DayCount        int
	DayCountSurplus int
}

type AccountMonthEntity struct {
	UserID            string
	ActivityID        int64
	Month             string
	MonthCount        int
	MonthCountSurplus int
}

type SkuProductEntity struct {
	SKU               int64
	ActivityID        int64
	ActivityCountID   int64
	StockCount        int
	StockCountSurplus int
	ProductAmount     float64
	ActivityCount     ActivityCountEntity
}

type ActivityCountEntity struct {
	ActivityCountID int64
	TotalCount      int
	DayCount        int
	MonthCount      int
}
