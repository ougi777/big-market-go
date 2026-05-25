package strategy

import (
	"strconv"
	"strings"
)

type StrategyEntity struct {
	StrategyID    int64
	StrategyDesc  string
	RuleModel     string
	RuleModelList []string
}

func (s StrategyEntity) RuleModels() []string {
	if len(s.RuleModelList) > 0 {
		return append([]string(nil), s.RuleModelList...)
	}
	if strings.TrimSpace(s.RuleModel) == "" {
		return nil
	}

	items := strings.FieldsFunc(s.RuleModel, func(r rune) bool {
		return r == ',' || r == ' ' || r == ';'
	})
	models := make([]string, 0, len(items))
	for _, item := range items {
		if model := strings.TrimSpace(item); model != "" {
			models = append(models, model)
		}
	}
	return models
}

func (s StrategyEntity) RuleWeight() string {
	for _, ruleModel := range s.RuleModels() {
		if ruleModel == "rule_weight" {
			return ruleModel
		}
	}
	return ""
}

type StrategyAwardEntity struct {
	StrategyID        int64
	AwardID           int
	AwardTitle        string
	AwardSubtitle     string
	AwardCount        int
	AwardCountSurplus int
	AwardRate         float64
	Sort              int
	RuleModels        string
}

type StrategyRuleEntity struct {
	StrategyID int64
	AwardID    *int
	RuleType   int
	RuleModel  string
	RuleValue  string
	RuleDesc   string
}

func (s StrategyRuleEntity) RuleWeightValues() map[string][]int {
	if s.RuleModel != "rule_weight" {
		return nil
	}

	groups := strings.Fields(s.RuleValue)
	result := make(map[string][]int, len(groups))
	for _, group := range groups {
		parts := strings.SplitN(group, ":", 2)
		if len(parts) != 2 {
			return result
		}

		values := strings.Split(parts[1], ",")
		awardIDs := make([]int, 0, len(values))
		for _, value := range values {
			value = strings.TrimSpace(value)
			if value == "" {
				continue
			}

			awardID, err := strconv.Atoi(value)
			if err == nil && awardID > 0 {
				awardIDs = append(awardIDs, awardID)
			}
		}
		result[group] = awardIDs
	}
	return result
}

type RuleWeight struct {
	RuleValue string
	Weight    int
	AwardIDs  []int
	AwardList []RuleWeightAward
}

type RuleWeightAward struct {
	AwardID    int
	AwardTitle string
}
