package sharding

import "testing"

func TestRouterTable(t *testing.T) {
	router := NewRouter(4)

	table := router.Table("raffle_activity_order", "xiaofuge")

	if table != "raffle_activity_order_001" {
		t.Fatalf("expected routed table, got %s", table)
	}
}

func TestRouterTableCompatibleWithJavaMiniDBRouter(t *testing.T) {
	router := NewRouterWithDBCount(2, 4)

	tests := map[string]string{
		"xiaofuge": "raffle_activity_order_001",
		"user001":  "raffle_activity_order_000",
		"user002":  "raffle_activity_order_001",
	}

	for userID, expected := range tests {
		if table := router.Table("raffle_activity_order", userID); table != expected {
			t.Fatalf("expected %s for %s, got %s", expected, userID, table)
		}
	}
}

func TestRouterDBKeyCompatibleWithJavaMiniDBRouter(t *testing.T) {
	router := NewRouterWithDBCount(2, 4)

	tests := map[string]string{
		"xiaofuge": "db01",
		"user001":  "db02",
		"user002":  "db02",
	}

	for userID, expected := range tests {
		if dbKey := router.DBKey(userID); dbKey != expected {
			t.Fatalf("expected %s for %s, got %s", expected, userID, dbKey)
		}
	}
}

func TestRouterTableDisabled(t *testing.T) {
	router := NewRouter(1)

	table := router.Table("raffle_activity_order", "xiaofuge")

	if table != "raffle_activity_order" {
		t.Fatalf("expected base table, got %s", table)
	}
}
