package sharding

import (
	"fmt"
	"hash/fnv"
)

type Router struct {
	TableCount int
}

func NewRouter(tableCount int) Router {
	if tableCount <= 0 {
		tableCount = 1
	}
	return Router{TableCount: tableCount}
}

func (r Router) Table(baseTable string, key string) string {
	if r.TableCount <= 1 || key == "" {
		return baseTable
	}
	return fmt.Sprintf("%s_%03d", baseTable, r.tableIndex(key))
}

func (r Router) tableIndex(key string) uint32 {
	hash := fnv.New32a()
	_, _ = hash.Write([]byte(key))
	return hash.Sum32() % uint32(r.TableCount)
}
