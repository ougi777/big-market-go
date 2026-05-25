package config

import (
	"testing"

	"github.com/spf13/viper"
)

func TestConfigHTTPAddr(t *testing.T) {
	cfg := Config{Server: ServerConfig{Port: 8091}}

	if cfg.HTTPAddr() != ":8091" {
		t.Fatalf("expected :8091, got %s", cfg.HTTPAddr())
	}
}

func TestSetDefaults(t *testing.T) {
	v := viper.New()
	setDefaults(v)

	if v.GetInt("server.port") != 8091 {
		t.Fatalf("expected default port 8091, got %d", v.GetInt("server.port"))
	}
	if v.GetInt("sharding.db_count") != 1 || v.GetInt("sharding.table_count") != 1 {
		t.Fatalf("expected default sharding 1/1, got %d/%d", v.GetInt("sharding.db_count"), v.GetInt("sharding.table_count"))
	}
	if v.GetString("redis.addr") != "localhost:16379" {
		t.Fatalf("expected default redis addr, got %s", v.GetString("redis.addr"))
	}
}
