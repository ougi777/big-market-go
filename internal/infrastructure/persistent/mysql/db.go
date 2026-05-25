package mysql

import (
	"bm-go/internal/config"

	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Open(cfg config.MySQLConfig) (*gorm.DB, error) {
	db, err := gorm.Open(gormmysql.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	return db, nil
}
