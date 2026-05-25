package strategy

import "strings"

type StrategyEntity struct {
	StrategyID    int64
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
