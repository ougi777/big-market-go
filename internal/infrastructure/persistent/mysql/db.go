package mysql

import (
	"fmt"

	"bm-go/internal/config"

	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Open(cfg config.MySQLConfig) (*gorm.DB, error) {
	return open(cfg.DSN, cfg.MaxIdleConns, cfg.MaxOpenConns)
}

type Router struct {
	defaultDB *gorm.DB
	shards    map[string]*gorm.DB
}

func OpenRouter(cfg config.MySQLConfig) (*Router, error) {
	defaultDB, err := Open(cfg)
	if err != nil {
		return nil, err
	}

	router := &Router{
		defaultDB: defaultDB,
		shards:    make(map[string]*gorm.DB, len(cfg.Shards)),
	}
	for key, shard := range cfg.Shards {
		db, err := open(shard.DSN, cfg.MaxIdleConns, cfg.MaxOpenConns)
		if err != nil {
			return nil, fmt.Errorf("open mysql shard %s: %w", key, err)
		}
		router.shards[key] = db
	}
	return router, nil
}

func (r *Router) Default() *gorm.DB {
	return r.defaultDB
}

func (r *Router) Shard(key string) *gorm.DB {
	if r == nil {
		return nil
	}
	if db, ok := r.shards[key]; ok {
		return db
	}
	return r.defaultDB
}

func open(dsn string, maxIdleConns int, maxOpenConns int) (*gorm.DB, error) {
	db, err := gorm.Open(gormmysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	return db, nil
}
