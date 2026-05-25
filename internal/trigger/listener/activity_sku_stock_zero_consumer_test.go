package listener

import (
	"context"
	"testing"
)

func TestActivitySkuStockZeroConsumerHandle(t *testing.T) {
	clearer := &fakeActivitySkuStockClearer{}
	consumer := NewActivitySkuStockZeroConsumer(nil, clearer, nil)

	err := consumer.handle(context.Background(), `{"id":"12345678901","timestamp":1779703200000,"data":9011}`)
	if err != nil {
		t.Fatalf("handle activity sku stock zero: %v", err)
	}

	if clearer.sku != 9011 {
		t.Fatalf("expected sku 9011, got %d", clearer.sku)
	}
}

type fakeActivitySkuStockClearer struct {
	sku int64
}

func (f *fakeActivitySkuStockClearer) ClearActivitySkuStock(ctx context.Context, sku int64) error {
	f.sku = sku
	return nil
}
