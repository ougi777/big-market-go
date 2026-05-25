package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"bm-go/internal/domain/strategy/rule/chain"
	strategyservice "bm-go/internal/domain/strategy/service"
	"bm-go/internal/types"
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

func TestStrategyArmoryRouteReturnsAppErrorCode(t *testing.T) {
	router := NewRouter(RouterOptions{
		ArmoryService: &fakeArmoryService{
			err: types.NewAppError(types.ResponseCodeStrategyRuleWeightNull, nil),
		},
	})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/raffle/strategy/strategy_armory?strategyId=100001", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), `"code":"ERR_BIZ_001"`) {
		t.Fatalf("expected app error code, got %s", recorder.Body.String())
	}
}

func TestStrategyArmoryRouteIllegalParam(t *testing.T) {
	router := NewRouter(RouterOptions{
		ArmoryService: &fakeArmoryService{},
	})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/raffle/strategy/strategy_armory?strategyId=0", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), `"code":"0002"`) {
		t.Fatalf("expected illegal param code, got %s", recorder.Body.String())
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

func TestRandomRaffleRouteReturnsAppErrorCode(t *testing.T) {
	router := NewRouter(RouterOptions{
		RaffleService: &fakeRaffleService{err: types.NewAppError(types.ResponseCodeUnassembledStrategy, nil)},
	})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/raffle/strategy/random_raffle", strings.NewReader(`{"strategyId":100001}`))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), `"code":"ERR_BIZ_002"`) {
		t.Fatalf("expected app error code, got %s", recorder.Body.String())
	}
}

func TestRandomRaffleRouteIllegalParam(t *testing.T) {
	router := NewRouter(RouterOptions{
		RaffleService: &fakeRaffleService{awardID: 101},
	})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/raffle/strategy/random_raffle", strings.NewReader(`{"strategyId":0}`))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), `"code":"0002"`) {
		t.Fatalf("expected illegal param code, got %s", recorder.Body.String())
	}
}

func TestQueryRaffleAwardListRoute(t *testing.T) {
	router := NewRouter(RouterOptions{
		QueryService: &fakeQueryService{
			awards: []strategyservice.RaffleAward{
				{
					AwardID:            101,
					AwardTitle:         "积分",
					AwardSubtitle:      "抽奖1次后解锁",
					Sort:               1,
					AwardRuleLockCount: 1,
					HasAwardRuleLock:   true,
					IsAwardUnlock:      true,
				},
			},
		},
	})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/raffle/strategy/query_raffle_award_list", strings.NewReader(`{"userId":"xiaofuge","activityId":100301}`))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), `"awardTitle":"积分"`) {
		t.Fatalf("expected award title, got %s", recorder.Body.String())
	}
}

func TestQueryRaffleAwardListRouteReturnsAppErrorCode(t *testing.T) {
	router := NewRouter(RouterOptions{
		QueryService: &fakeQueryService{err: types.NewAppError(types.ResponseCodeAccountQuotaError, nil)},
	})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/raffle/strategy/query_raffle_award_list", strings.NewReader(`{"userId":"xiaofuge","activityId":100301}`))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), `"code":"ERR_BIZ_006"`) {
		t.Fatalf("expected app error code, got %s", recorder.Body.String())
	}
}

func TestQueryRaffleAwardListRouteIllegalParam(t *testing.T) {
	router := NewRouter(RouterOptions{
		QueryService: &fakeQueryService{},
	})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/raffle/strategy/query_raffle_award_list", strings.NewReader(`{"userId":"","activityId":100301}`))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), `"code":"0002"`) {
		t.Fatalf("expected illegal param code, got %s", recorder.Body.String())
	}
}

func TestQueryRaffleStrategyRuleWeightRoute(t *testing.T) {
	router := NewRouter(RouterOptions{
		QueryService: &fakeQueryService{
			ruleWeights: []strategyservice.RaffleStrategyRuleWeight{
				{
					RuleWeightCount:                  4000,
					UserActivityAccountTotalUseCount: 4500,
					StrategyAwards: []strategyservice.RuleWeightAward{
						{AwardID: 101, AwardTitle: "积分"},
					},
				},
			},
		},
	})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/raffle/strategy/query_raffle_strategy_rule_weight", strings.NewReader(`{"userId":"xiaofuge","activityId":100301}`))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), `"ruleWeightCount":4000`) {
		t.Fatalf("expected rule weight count, got %s", recorder.Body.String())
	}
}

func TestQueryRaffleStrategyRuleWeightRouteReturnsAppErrorCode(t *testing.T) {
	router := NewRouter(RouterOptions{
		QueryService: &fakeQueryService{err: types.NewAppError(types.ResponseCodeStrategyRuleWeightNull, nil)},
	})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/raffle/strategy/query_raffle_strategy_rule_weight", strings.NewReader(`{"userId":"xiaofuge","activityId":100301}`))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), `"code":"ERR_BIZ_001"`) {
		t.Fatalf("expected app error code, got %s", recorder.Body.String())
	}
}

func TestQueryRaffleStrategyRuleWeightRouteIllegalParam(t *testing.T) {
	router := NewRouter(RouterOptions{
		QueryService: &fakeQueryService{},
	})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/raffle/strategy/query_raffle_strategy_rule_weight", strings.NewReader(`{"userId":"xiaofuge","activityId":0}`))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), `"code":"0002"`) {
		t.Fatalf("expected illegal param code, got %s", recorder.Body.String())
	}
}

type fakeArmoryService struct {
	err error
}

func (f *fakeArmoryService) AssembleLotteryStrategy(ctx context.Context, strategyID int64) error {
	if f.err != nil {
		return f.err
	}
	return nil
}

type fakeRaffleService struct {
	awardID int
	err     error
}

func (f *fakeRaffleService) PerformRaffle(ctx context.Context, userID string, strategyID int64) (chain.AwardResult, error) {
	if f.err != nil {
		return chain.AwardResult{}, f.err
	}
	return chain.AwardResult{AwardID: f.awardID}, nil
}

type fakeQueryService struct {
	awards      []strategyservice.RaffleAward
	ruleWeights []strategyservice.RaffleStrategyRuleWeight
	err         error
}

func (f *fakeQueryService) QueryRaffleAwardList(ctx context.Context, activityID int64, userID string) ([]strategyservice.RaffleAward, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.awards, nil
}

func (f *fakeQueryService) QueryRaffleStrategyRuleWeight(ctx context.Context, activityID int64, userID string) ([]strategyservice.RaffleStrategyRuleWeight, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.ruleWeights, nil
}
