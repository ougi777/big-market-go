package types

type ResponseCode struct {
	Code string
	Info string
}

var (
	ResponseCodeSuccess                 = ResponseCode{Code: "0000", Info: "success"}
	ResponseCodeUnknownError            = ResponseCode{Code: "0001", Info: "unknown error"}
	ResponseCodeIllegalParam            = ResponseCode{Code: "0002", Info: "illegal parameter"}
	ResponseCodeIndexDup                = ResponseCode{Code: "0003", Info: "duplicate key"}
	ResponseCodeStrategyRuleWeightNull  = ResponseCode{Code: "ERR_BIZ_001", Info: "strategy rule_weight is not configured"}
	ResponseCodeUnassembledStrategy     = ResponseCode{Code: "ERR_BIZ_002", Info: "raffle strategy armory is not assembled"}
	ResponseCodeActivityStateError      = ResponseCode{Code: "ERR_BIZ_003", Info: "activity state error"}
	ResponseCodeActivityDateError       = ResponseCode{Code: "ERR_BIZ_004", Info: "activity date error"}
	ResponseCodeAccountQuotaError       = ResponseCode{Code: "ERR_BIZ_006", Info: "account quota error"}
	ResponseCodeAccountMonthQuotaError  = ResponseCode{Code: "ERR_BIZ_007", Info: "account month quota error"}
	ResponseCodeAccountDayQuotaError    = ResponseCode{Code: "ERR_BIZ_008", Info: "account day quota error"}
	ResponseCodeActivityOrderStateError = ResponseCode{Code: "ERR_BIZ_009", Info: "activity order state error"}
)
