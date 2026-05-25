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

func TestRouterDefaultCountWhenInvalid(t *testing.T) {
	router := NewRouterWithDBCount(0, 0)

	if router.DBCount != 1 || router.TableCount != 1 {
		t.Fatalf("expected default count 1/1, got %d/%d", router.DBCount, router.TableCount)
	}
	if table := router.Table("raffle_activity_order", "xiaofuge"); table != "raffle_activity_order" {
		t.Fatalf("expected base table, got %s", table)
	}
}

func TestRouterTableEmptyKey(t *testing.T) {
	router := NewRouterWithDBCount(2, 4)

	if table := router.Table("raffle_activity_order", ""); table != "raffle_activity_order" {
		t.Fatalf("expected base table, got %s", table)
	}
}

func TestRouterDBKeyDefault(t *testing.T) {
	router := NewRouterWithDBCount(1, 4)

	if dbKey := router.DBKey("xiaofuge"); dbKey != "default" {
		t.Fatalf("expected default db, got %s", dbKey)
	}
	if dbKey := NewRouterWithDBCount(2, 4).DBKey(""); dbKey != "default" {
		t.Fatalf("expected default db for empty key, got %s", dbKey)
	}
}
