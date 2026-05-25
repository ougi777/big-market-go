package http

import (
	"errors"
	stdhttp "net/http"

	"bm-go/internal/types"

	"github.com/gin-gonic/gin"
)

func appErrorCode(err error) (types.ResponseCode, bool) {
	var appErr types.AppError
	if errors.As(err, &appErr) {
		return appErr.Code, true
	}
	return types.ResponseCode{}, false
}

func writeAppErrorFailure[T any](ctx *gin.Context, err error, data T) bool {
	code, ok := appErrorCode(err)
	if !ok {
		return false
	}
	ctx.JSON(stdhttp.StatusOK, types.Failure(code, data))
	return true
}
