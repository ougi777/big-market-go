package repository

import (
	"context"
	"errors"
	"testing"

	"bm-go/internal/types"
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

type fakeRateTableStore struct {
	values map[string]string
	hashes map[string]map[string]string
}

func (f *fakeRateTableStore) Get(ctx context.Context, key string) (string, error) {
	return f.values[key], nil
}

func (f *fakeRateTableStore) HGet(ctx context.Context, key string, field string) (string, error) {
	return f.hashes[key][field], nil
}

func (f *fakeRateTableStore) Decr(ctx context.Context, key string) (int64, error) {
	return 0, nil
}

func (f *fakeRateTableStore) Set(ctx context.Context, key string, value string) error {
	return nil
}

func (f *fakeRateTableStore) SetNX(ctx context.Context, key string, value string) (bool, error) {
	return true, nil
}
