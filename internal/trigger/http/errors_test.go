package http

import (
	"errors"
	"testing"

	"bm-go/internal/types"
)

func TestAppErrorCode(t *testing.T) {
	code, ok := appErrorCode(types.NewAppError(types.ResponseCodeActivityDateError, errors.New("date error")))
	if !ok {
		t.Fatal("expected app error")
	}
	if code != types.ResponseCodeActivityDateError {
		t.Fatalf("expected activity date error, got %s", code.Code)
	}
}

func TestAppErrorCodeNonAppError(t *testing.T) {
	code, ok := appErrorCode(errors.New("plain error"))
	if ok {
		t.Fatalf("expected non app error, got %s", code.Code)
	}
}
