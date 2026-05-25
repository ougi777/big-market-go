package http

import (
	"errors"

	"bm-go/internal/types"
)

func appErrorCode(err error) (types.ResponseCode, bool) {
	var appErr types.AppError
	if errors.As(err, &appErr) {
		return appErr.Code, true
	}
	return types.ResponseCode{}, false
}
