package tree

const (
	RuleLock      = "rule_lock"
	RuleStock     = "rule_stock"
	RuleLuckAward = "rule_luck_award"
)

const (
	CheckAllow    = "0000"
	CheckTakeOver = "0001"
)

const LimitEqual = "EQUAL"

type RuleTree struct {
	TreeID   string
	TreeName string
	TreeDesc string
	RootRule string
	NodeMap  map[string]RuleTreeNode
}

type RuleTreeNode struct {
	TreeID    string
	RuleKey   string
	RuleDesc  string
	RuleValue string
	Lines     []RuleTreeNodeLine
}

type RuleTreeNodeLine struct {
	TreeID         string
	RuleNodeFrom   string
	RuleNodeTo     string
	RuleLimitType  string
	RuleLimitValue string
}

type TreeAction struct {
	CheckType CheckType
	Award     *StrategyAward
}

type CheckType struct {
	Code string
	Name string
}

type StrategyAward struct {
	AwardID        int
	AwardRuleValue string
}

func Allow() TreeAction {
	return TreeAction{CheckType: CheckType{Code: CheckAllow, Name: "ALLOW"}}
}

func TakeOver(award *StrategyAward) TreeAction {
	return TreeAction{
		CheckType: CheckType{Code: CheckTakeOver, Name: "TAKE_OVER"},
		Award:     award,
	}
}
