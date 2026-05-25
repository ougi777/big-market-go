package types

type ResponseCode struct {
	Code string
	Info string
}

var (
	ResponseCodeSuccess      = ResponseCode{Code: "0000", Info: "success"}
	ResponseCodeUnknownError = ResponseCode{Code: "0001", Info: "unknown error"}
	ResponseCodeIllegalParam = ResponseCode{Code: "0002", Info: "illegal parameter"}
)
