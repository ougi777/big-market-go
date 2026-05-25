package types

import (
	"errors"
	"testing"
)

func TestAppErrorWithInnerError(t *testing.T) {
	err := NewAppError(ResponseCodeActivityDateError, errors.New("date invalid"))

	if err.Error() != "date invalid" {
		t.Fatalf("expected inner error message, got %s", err.Error())
	}
}

func TestAppErrorWithCodeInfo(t *testing.T) {
	err := NewAppError(ResponseCodeActivityDateError, nil)

	if err.Error() != ResponseCodeActivityDateError.Info {
		t.Fatalf("expected code info, got %s", err.Error())
	}
}
