package sharding

import "testing"

func TestRouterTable(t *testing.T) {
	router := NewRouter(4)

	table := router.Table("raffle_activity_order", "xiaofuge")

	if table != "raffle_activity_order_003" {
		t.Fatalf("expected routed table, got %s", table)
	}
}

func TestRouterTableDisabled(t *testing.T) {
	router := NewRouter(1)

	table := router.Table("raffle_activity_order", "xiaofuge")

	if table != "raffle_activity_order" {
		t.Fatalf("expected base table, got %s", table)
	}
}
