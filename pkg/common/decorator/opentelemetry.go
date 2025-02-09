package decorator

import (
	"context"
	"strings"
	"time"

	"github.com/durianpay/dpay-common/constants"
	"github.com/layarda-durianpay/go-skeleton/pkg/common/opentelemetry"
	"github.com/layarda-durianpay/go-skeleton/pkg/common/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type commandOTelDecorator[C any] struct {
	base CommandHandler[C]
}

func (d commandOTelDecorator[C]) Handle(ctx context.Context, cmd C) (err error) {
	now := time.Now().UTC()

	actionName := strings.ToLower(generateActionName(cmd))

	ctx, span := opentelemetry.StartSpan(
		ctx,
		actionName,
		trace.WithSpanKind(trace.SpanKindInternal),
	)
	defer span.End()

	span.SetAttributes(
		attribute.String("cqrs.type", "command"),
		attribute.String("cqrs.operation", actionName),
		attribute.Stringer("cqrs.timestamp", now),
		attribute.String(opentelemetry.RequestIDKey, utils.GetFromContext[string](ctx, constants.RequestIDKey)),
	)

	err = d.base.Handle(ctx, cmd)
	if nil != err {
		span.RecordError(err, trace.WithStackTrace(true))
	}

	return
}

type queryOTelDecorator[Q any, R any] struct {
	base QueryHandler[Q, R]
}

func (d queryOTelDecorator[Q, R]) Handle(ctx context.Context, query Q) (result R, err error) {
	now := time.Now().UTC()

	actionName := strings.ToLower(generateActionName(query))

	ctx, span := opentelemetry.StartSpan(
		ctx,
		actionName,
		trace.WithSpanKind(trace.SpanKindInternal),
	)
	defer span.End()

	span.SetAttributes(
		attribute.String("cqrs.type", "query"),
		attribute.String("cqrs.operation", actionName),
		attribute.Stringer("cqrs.timestamp", now),
		attribute.String(opentelemetry.RequestIDKey, utils.GetFromContext[string](ctx, constants.RequestIDKey)),
	)

	result, err = d.base.Handle(ctx, query)
	if nil != err {
		span.RecordError(err, trace.WithStackTrace(true))
	}

	return result, err
}
