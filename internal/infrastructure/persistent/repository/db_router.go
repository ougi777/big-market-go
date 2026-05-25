package repository

import "gorm.io/gorm"

type dbRouter interface {
	Default() *gorm.DB
	Shard(key string) *gorm.DB
	Connections() []*gorm.DB
}

type singleDBRouter struct {
	db *gorm.DB
}

func (r singleDBRouter) Default() *gorm.DB {
	return r.db
}

func (r singleDBRouter) Shard(string) *gorm.DB {
	return r.db
}

func (r singleDBRouter) Connections() []*gorm.DB {
	return []*gorm.DB{r.db}
}
