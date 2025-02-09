package decorator

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/durianpay/dpay-common/logger"
	"github.com/layarda-durianpay/go-skeleton/pkg/common/errors"
	"github.com/samber/lo"
)

type commandLoggingDecorator[C any] struct {
	base CommandHandler[C]
}

func (d commandLoggingDecorator[C]) Handle(ctx context.Context, cmd C) (err error) {
	fields := []any{
		"command", generateActionName(cmd),
		"command_body", getLogBodyParam(cmd),
	}

	logger.Infow(ctx, "Executing command", fields...)

	defer func() {
		errorLog(ctx, fields, err, "command")
	}()

	return d.base.Handle(ctx, cmd)
}

type queryLoggingDecorator[Q any, R any] struct {
	base QueryHandler[Q, R]
}

func (d queryLoggingDecorator[Q, R]) Handle(ctx context.Context, qry Q) (result R, err error) {
	fields := []any{
		"query", generateActionName(qry),
		"query_body", getLogBodyParam(qry),
	}

	logger.Infow(ctx, "Executing query", fields...)

	defer func() {
		errorLog(ctx, fields, err, "query")
	}()

	return d.base.Handle(ctx, qry)
}

// errorLog logs the error with the given key ("command" or "query")
func errorLog(ctx context.Context, fields []any, err error, key string) {
	if err == nil {
		return
	}

	fields = append(fields, "error", err.Error())

	if dpayErr, ok := lo.ErrorsAs[errors.DpayError](err); ok {
		fields = append(fields, "error_original", errors.GetOriginalErr(dpayErr))
	}

	if errors.IsClientError(err) {
		logger.Warnw(
			ctx,
			fmt.Sprintf("client error: failed to execute: %s", key),
			fields...,
		)
		return
	}

	logger.Errorw(
		ctx,
		fmt.Sprintf("internal error: failed to execute: %s", key),
		fields...,
	)
}

func getLogBodyParam(param any) any {
	var body any = fmt.Sprintf("%#v", param)

	jsonBytes, err := json.Marshal(param)
	if err == nil {
		var param map[string]any
		if err := json.Unmarshal(jsonBytes, &param); err == nil {
			body = param
		}
	}

	return body
}
