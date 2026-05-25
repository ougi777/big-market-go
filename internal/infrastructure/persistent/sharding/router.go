package sharding

import (
	"fmt"
)

type Router struct {
	DBCount    int
	TableCount int
}

func NewRouter(tableCount int) Router {
	return NewRouterWithDBCount(1, tableCount)
}

func NewRouterWithDBCount(dbCount int, tableCount int) Router {
	if dbCount <= 0 {
		dbCount = 1
	}
	if tableCount <= 0 {
		tableCount = 1
	}
	return Router{DBCount: dbCount, TableCount: tableCount}
}

func (r Router) Table(baseTable string, key string) string {
	if r.TableCount <= 1 || key == "" {
		return baseTable
	}
	return fmt.Sprintf("%s_%03d", baseTable, r.tableIndex(key))
}

func (r Router) tableIndex(key string) int32 {
	idx := r.routeIndex(key)
	tableCount := int32(r.TableCount)
	dbIdx := idx/tableCount + 1
	return idx - tableCount*(dbIdx-1)
}

func (r Router) routeIndex(key string) int32 {
	total := int32(r.DBCount * r.TableCount)
	hash := javaHashCode(key)
	spread := hash ^ int32(uint32(hash)>>16)
	return (total - 1) & spread
}

func javaHashCode(value string) int32 {
	var hash int32
	for _, r := range value {
		hash = 31*hash + r
	}
	return hash
}
