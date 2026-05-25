package types

type AppError struct {
	Code ResponseCode
	Err  error
}

func (e AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Code.Info
}

func NewAppError(code ResponseCode, err error) AppError {
	return AppError{Code: code, Err: err}
}
