package misc

import (
	"context"
	"time"

	"github.com/cuigh/auxo/errors"
)

const (
	ErrInvalidToken         = 1001
	ErrAccountDisabled      = 1002
	ErrOldPasswordIncorrect = 1003
	ErrExternalStack        = 1004
	ErrSystemInitialized    = 1005
)

func Error(code int32, err error) error {
	return errors.Coded(code, err.Error())
}

func Page(count, pageIndex, pageSize int) (start, end int) {
	start = pageSize * (pageIndex - 1)
	end = pageSize * pageIndex
	if count < start {
		start, end = 0, 0
	} else if count < end {
		end = count
	}
	return
}

func Context(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}
