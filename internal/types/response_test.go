package types

import "testing"

func TestSuccessResponse(t *testing.T) {
	response := Success("pong")

	if response.Code != ResponseCodeSuccess.Code || response.Info != ResponseCodeSuccess.Info {
		t.Fatalf("expected success response, got %+v", response)
	}
	if response.Data != "pong" {
		t.Fatalf("expected data pong, got %s", response.Data)
	}
}

func TestFailureResponse(t *testing.T) {
	response := Failure(ResponseCodeIllegalParam, []string{})

	if response.Code != ResponseCodeIllegalParam.Code || response.Info != ResponseCodeIllegalParam.Info {
		t.Fatalf("expected illegal param response, got %+v", response)
	}
	if response.Data == nil {
		t.Fatal("expected data")
	}
}
