package utils

import (
	"context"

	"github.com/durianpay/dpay-common/constants"
)

func GetFromContext[T any](ctx context.Context, key constants.ContextKey) (res T) {
	val, ok := ctx.Value(key).(T)
	if !ok {
		return res
	}

	return val
}
