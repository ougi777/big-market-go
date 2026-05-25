package strategy

import (
	"reflect"
	"testing"
)

func TestStrategyEntityRuleModels(t *testing.T) {
	entity := StrategyEntity{RuleModel: "rule_weight, rule_blacklist;rule_lock"}

	got := entity.RuleModels()
	want := []string{"rule_weight", "rule_blacklist", "rule_lock"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected rule models %v, got %v", want, got)
	}
}

func TestStrategyRuleEntityRuleWeightValues(t *testing.T) {
	entity := StrategyRuleEntity{
		RuleModel: "rule_weight",
		RuleValue: "4000:102,103,104,105 5000:102,103,104,105,106",
	}

	got := entity.RuleWeightValues()
	want := map[string][]int{
		"4000:102,103,104,105":     {102, 103, 104, 105},
		"5000:102,103,104,105,106": {102, 103, 104, 105, 106},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected rule weight values %v, got %v", want, got)
	}
}

func TestStrategyRuleEntityRuleWeightValuesIgnoresInvalidAwardID(t *testing.T) {
	entity := StrategyRuleEntity{
		RuleModel: "rule_weight",
		RuleValue: "4000:102,abc,0,103",
	}

	got := entity.RuleWeightValues()
	want := map[string][]int{
		"4000:102,abc,0,103": {102, 103},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected rule weight values %v, got %v", want, got)
	}
}
