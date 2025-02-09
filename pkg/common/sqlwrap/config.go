package sqlwrap

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

type configSQL struct {
	useOTel     bool
	serviceName string
}

func (c configSQL) setConfig(db *DB) {
	if c.useOTel {
		c.configurateOTel(db)
	}
}

func (c configSQL) configurateOTel(db *DB) {
	// var commonConfig env.CommonConfig

	// env.MustProcess(&commonConfig, nil)

	db.AddBeforeFunc(
		BeforeFunc(func(ctx context.Context, query string, args ...interface{}) context.Context {
			now := time.Now().UTC()

			spanName := "repository"
			if caller, ok := ctx.Value(SQLWrapperCallerKey).(string); caller != "" &&
				ok {
				spanName = caller
			}

			ctx, span := opentelemetry.StartSpan(
				ctx,
				spanName,
				trace.WithSpanKind(trace.SpanKindInternal),
			)
			span.SetAttributes(
				attribute.Stringer("sql.timestamp", now),
				attribute.String(
					opentelemetry.RequestIDKey,
					utils.GetFromContext[string](ctx, constants.RequestIDKey),
				),
			)

			return ctx
		}),
	)

	db.AddAfterFunc(
		AfterFunc(func(ctx context.Context, err error, query string, args ...interface{}) {
			span := trace.SpanFromContext(ctx)
			if span == nil {
				return
			}

			argAttributes := utils.ToOtelAttributes("sql.query.args", args)
			query = strings.ReplaceAll(query, "\n", " ")
			query = strings.ReplaceAll(query, "\t", " ")
			argAttributes = append(
				argAttributes,
				attribute.String("sql.query", query),
			)

			if strings.TrimSpace(c.serviceName) != "" {
				argAttributes = append(
					argAttributes,
					attribute.String("repository.name", c.serviceName),
				)
			}

			span.SetAttributes(argAttributes...)

			if nil != err {
				span.RecordError(err, trace.WithStackTrace(true))
			}

			span.End()
		}),
	)
}

type Option interface {
	apply(c *configSQL)
}

type optionFunc func(*configSQL)

func (o optionFunc) apply(c *configSQL) {
	o(c)
}

func newConfig(opts ...Option) *configSQL {
	conf := &configSQL{}

	for _, opt := range opts {
		opt.apply(conf)
	}

	return conf
}

func WithOTelOption(enableOpentelemetry bool) Option {
	return optionFunc(func(conf *configSQL) {
		conf.useOTel = enableOpentelemetry
	})
}

func ServiceNameOption(name string) Option {
	return optionFunc(func(conf *configSQL) {
		conf.serviceName = name
	})
}
