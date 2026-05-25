package listener

import (
	"errors"

	"bm-go/internal/types"
)

func isIndexDuplicateError(err error) bool {
	var appErr types.AppError
	return errors.As(err, &appErr) && appErr.Code == types.ResponseCodeIndexDup
}
