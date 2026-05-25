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
	if v.GetInt("mysql.max_idle_conns") != 10 || v.GetInt("mysql.max_open_conns") != 50 {
		t.Fatalf("expected default mysql pool 10/50, got %d/%d", v.GetInt("mysql.max_idle_conns"), v.GetInt("mysql.max_open_conns"))
	}
	if v.GetString("rabbitmq.url") != "amqp://admin:admin@localhost:5672/" {
		t.Fatalf("expected default rabbitmq url, got %s", v.GetString("rabbitmq.url"))
	}
	if v.GetString("log.level") != "info" {
		t.Fatalf("expected default log level info, got %s", v.GetString("log.level"))
	}
}

func TestNewLogger(t *testing.T) {
	logger, err := NewLogger(LogConfig{Level: "debug"})
	if err != nil {
		t.Fatalf("new logger: %v", err)
	}

	if !logger.Core().Enabled(-1) {
		t.Fatal("expected debug log enabled")
	}
}

func TestNewLoggerFallback(t *testing.T) {
	logger, err := NewLogger(LogConfig{Level: "bad-level"})
	if err != nil {
		t.Fatalf("new logger: %v", err)
	}

	if logger.Core().Enabled(-1) {
		t.Fatal("expected debug log disabled for fallback info level")
	}
}
