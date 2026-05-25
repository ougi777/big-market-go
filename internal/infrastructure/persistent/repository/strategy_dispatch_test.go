package repository

import (
	"context"
	"errors"
	"testing"

	"bm-go/internal/types"

	"github.com/redis/go-redis/v9"
)

func TestStrategyDispatchGetRandomAwardID(t *testing.T) {
	store := &fakeRateTableStore{
		values: map[string]string{
			types.RedisKeyStrategyRateRange + "100001": "1",
		},
		hashes: map[string]map[string]string{
			types.RedisKeyStrategyRateTable + "100001": {
				"0": "101",
			},
		},
	}
	dispatch := NewStrategyDispatchWithStore(store)

	awardID, err := dispatch.GetRandomAwardID(context.Background(), 100001)
	if err != nil {
		t.Fatalf("get random award id: %v", err)
	}
	if awardID != 101 {
		t.Fatalf("expected award 101, got %d", awardID)
	}
}

func TestStrategyDispatchGetWeightedRandomAwardID(t *testing.T) {
	store := &fakeRateTableStore{
		values: map[string]string{
			types.RedisKeyStrategyRateRange + "100001_5000:104,105": "1",
		},
		hashes: map[string]map[string]string{
			types.RedisKeyStrategyRateTable + "100001_5000:104,105": {
				"0": "104",
			},
		},
	}
	dispatch := NewStrategyDispatchWithStore(store)

	awardID, err := dispatch.GetRandomAwardID(context.Background(), 100001, "5000:104,105")
	if err != nil {
		t.Fatalf("get weighted random award id: %v", err)
	}
	if awardID != 104 {
		t.Fatalf("expected award 104, got %d", awardID)
	}
}

func TestStrategyDispatchReturnsUnassembledErrorWhenStoreMissing(t *testing.T) {
	dispatch := NewStrategyDispatchWithStore(nil)

	_, err := dispatch.GetRandomAwardID(context.Background(), 100001)
	var appErr types.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected app error, got %v", err)
	}
	if appErr.Code != types.ResponseCodeUnassembledStrategy {
		t.Fatalf("expected unassembled strategy error, got %s", appErr.Code.Code)
	}
}

func TestStrategyDispatchReturnsUnassembledErrorWhenRateRangeMissing(t *testing.T) {
	store := &fakeRateTableStore{
		getErr: redis.Nil,
	}
	dispatch := NewStrategyDispatchWithStore(store)

	_, err := dispatch.GetRandomAwardID(context.Background(), 100001)
	var appErr types.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected app error, got %v", err)
	}
	if appErr.Code != types.ResponseCodeUnassembledStrategy {
		t.Fatalf("expected unassembled strategy error, got %s", appErr.Code.Code)
	}
}

func TestStrategyDispatchReturnsErrorWhenRateRangeInvalid(t *testing.T) {
	store := &fakeRateTableStore{
		values: map[string]string{
			types.RedisKeyStrategyRateRange + "100001": "bad",
		},
	}
	dispatch := NewStrategyDispatchWithStore(store)

	_, err := dispatch.GetRandomAwardID(context.Background(), 100001)
	if err == nil {
		t.Fatal("expected rate range parse error")
	}
}

func TestStrategyDispatchReturnsErrorWhenRateRangeZero(t *testing.T) {
	store := &fakeRateTableStore{
		values: map[string]string{
			types.RedisKeyStrategyRateRange + "100001": "0",
		},
	}
	dispatch := NewStrategyDispatchWithStore(store)

	_, err := dispatch.GetRandomAwardID(context.Background(), 100001)
	if err == nil {
		t.Fatal("expected positive rate range error")
	}
}

func TestStrategyDispatchReturnsUnassembledErrorWhenRateTableMissing(t *testing.T) {
	store := &fakeRateTableStore{
		values: map[string]string{
			types.RedisKeyStrategyRateRange + "100001": "1",
		},
		hErr: redis.Nil,
	}
	dispatch := NewStrategyDispatchWithStore(store)

	_, err := dispatch.GetRandomAwardID(context.Background(), 100001)
	var appErr types.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected app error, got %v", err)
	}
	if appErr.Code != types.ResponseCodeUnassembledStrategy {
		t.Fatalf("expected unassembled strategy error, got %s", appErr.Code.Code)
	}
}

func TestStrategyDispatchReturnsErrorWhenAwardIDInvalid(t *testing.T) {
	store := &fakeRateTableStore{
		values: map[string]string{
			types.RedisKeyStrategyRateRange + "100001": "1",
		},
		hashes: map[string]map[string]string{
			types.RedisKeyStrategyRateTable + "100001": {
				"0": "bad",
			},
		},
	}
	dispatch := NewStrategyDispatchWithStore(store)

	_, err := dispatch.GetRandomAwardID(context.Background(), 100001)
	if err == nil {
		t.Fatal("expected award id parse error")
	}
}

func TestStrategyDispatchSubtractionAwardStock(t *testing.T) {
	store := &fakeRateTableStore{decrValue: 9, setNXValue: true}
	dispatch := NewStrategyDispatchWithStore(store)

	ok, err := dispatch.SubtractionAwardStock(context.Background(), 100001, 101)
	if err != nil {
		t.Fatalf("subtract award stock: %v", err)
	}
	if !ok {
		t.Fatal("expected stock subtraction ok")
	}
	expectedKey := types.RedisKeyStrategyAwardCount + "100001_101"
	if store.decrKey != expectedKey {
		t.Fatalf("expected decr key %s, got %s", expectedKey, store.decrKey)
	}
	if store.setNXKey != expectedKey+"_9" {
		t.Fatalf("expected lock key %s, got %s", expectedKey+"_9", store.setNXKey)
	}
}

func TestStrategyDispatchSubtractionAwardStockSoldOut(t *testing.T) {
	store := &fakeRateTableStore{decrValue: -1}
	dispatch := NewStrategyDispatchWithStore(store)

	ok, err := dispatch.SubtractionAwardStock(context.Background(), 100001, 101)
	if err != nil {
		t.Fatalf("subtract award stock: %v", err)
	}
	if ok {
		t.Fatal("expected sold out")
	}
	expectedKey := types.RedisKeyStrategyAwardCount + "100001_101"
	if store.setKey != expectedKey || store.setValue != "0" {
		t.Fatalf("expected stock reset to 0, got %s=%s", store.setKey, store.setValue)
	}
}

func TestStrategyDispatchSubtractionAwardStockDecrError(t *testing.T) {
	store := &fakeRateTableStore{decrErr: errors.New("decr failed")}
	dispatch := NewStrategyDispatchWithStore(store)

	_, err := dispatch.SubtractionAwardStock(context.Background(), 100001, 101)
	if err == nil {
		t.Fatal("expected decr error")
	}
}

func TestStrategyDispatchSubtractionAwardStockLockConflict(t *testing.T) {
	store := &fakeRateTableStore{decrValue: 8, setNXValue: false}
	dispatch := NewStrategyDispatchWithStore(store)

	ok, err := dispatch.SubtractionAwardStock(context.Background(), 100001, 101)
	if err != nil {
		t.Fatalf("subtract award stock: %v", err)
	}
	if ok {
		t.Fatal("expected lock conflict")
	}
}

type fakeRateTableStore struct {
	values     map[string]string
	hashes     map[string]map[string]string
	getErr     error
	hErr       error
	decrKey    string
	decrValue  int64
	decrErr    error
	setKey     string
	setValue   string
	setErr     error
	setNXKey   string
	setNXValue bool
	setNXErr   error
}

func (f *fakeRateTableStore) Get(ctx context.Context, key string) (string, error) {
	if f.getErr != nil {
		return "", f.getErr
	}
	return f.values[key], nil
}

func (f *fakeRateTableStore) HGet(ctx context.Context, key string, field string) (string, error) {
	if f.hErr != nil {
		return "", f.hErr
	}
	return f.hashes[key][field], nil
}

func (f *fakeRateTableStore) Decr(ctx context.Context, key string) (int64, error) {
	f.decrKey = key
	return f.decrValue, f.decrErr
}

func (f *fakeRateTableStore) Set(ctx context.Context, key string, value string) error {
	f.setKey = key
	f.setValue = value
	return f.setErr
}

func (f *fakeRateTableStore) SetNX(ctx context.Context, key string, value string) (bool, error) {
	f.setNXKey = key
	return f.setNXValue, f.setNXErr
}
