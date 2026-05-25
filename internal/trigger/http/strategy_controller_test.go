package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"bm-go/internal/domain/strategy/rule/chain"
)

func TestStrategyArmoryRoute(t *testing.T) {
	router := NewRouter(RouterOptions{
		ArmoryService: &fakeArmoryService{},
	})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/raffle/strategy/strategy_armory?strategyId=100001", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), `"data":true`) {
		t.Fatalf("expected success data true, got %s", recorder.Body.String())
	}
}

func TestRandomRaffleRoute(t *testing.T) {
	router := NewRouter(RouterOptions{
		RaffleService: &fakeRaffleService{awardID: 101},
	})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/raffle/strategy/random_raffle", strings.NewReader(`{"strategyId":100001}`))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), `"awardId":101`) {
		t.Fatalf("expected award id 101, got %s", recorder.Body.String())
	}
}

type fakeArmoryService struct{}

func (f *fakeArmoryService) AssembleLotteryStrategy(ctx context.Context, strategyID int64) error {
	return nil
}

type fakeRaffleService struct {
	awardID int
}

func (f *fakeRaffleService) PerformRaffle(ctx context.Context, userID string, strategyID int64) (chain.AwardResult, error) {
	return chain.AwardResult{AwardID: f.awardID}, nil
}
