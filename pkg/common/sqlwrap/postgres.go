package sqlwrap

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/XSAM/otelsql"
	commoncfg "github.com/durianpay/dpay-common/config"
	"github.com/durianpay/dpay-common/db"
	"github.com/durianpay/dpay-common/logger"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func ProvidePostgres(opts ...Option) Database {
	dbSqlx := initDB()

	if err := dbSqlx.Ping(); err != nil {
		panic(fmt.Sprintf("postgres: %s", err.Error()))
	}

	wrappedDB := NewDB(dbSqlx)
	wrappedDB.AddBeforeFunc(postgresBeforeFunc)
	wrappedDB.AddAfterFunc(postgresAfterFunc)

	DBConf := newConfig(opts...)
	DBConf.setConfig(wrappedDB)

	return wrappedDB
}

func initDB() *sqlx.DB {
	dbCfg := commoncfg.NewDatabaseConfig()

	connStr := dbCfg.ConnectionURL()

	otelDB, err := otelsql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	err = otelsql.RegisterDBStatsMetrics(
		otelDB,
		otelsql.WithAttributes(semconv.DBSystemPostgreSQL),
	)
	if err != nil {
		panic(err)
	}

	dbSqlx := sqlx.NewDb(otelDB, db.Driver)

	dbSqlx.SetConnMaxLifetime(dbCfg.ConnMaxLifeTime())
	dbSqlx.SetConnMaxIdleTime(dbCfg.ConnMaxIdleTime())
	dbSqlx.SetMaxOpenConns(dbCfg.MaxPoolSize())
	dbSqlx.SetMaxIdleConns(dbCfg.MaxIdleConns())

	return dbSqlx
}

func postgresBeforeFunc(ctx context.Context, query string, args ...interface{}) context.Context {
	ctx = context.WithValue(ctx, queryStartKey, time.Now())
	return ctx
}

func postgresAfterFunc(ctx context.Context, err error, query string, args ...interface{}) {
	query = strings.ReplaceAll(query, "\n", " ")
	query = strings.ReplaceAll(query, "\t", " ")

	fields := []any{
		"exec_time_elapsed", getRequestDuration(ctx) / time.Millisecond,
		"query", query,
		"args", args,
		"caller", ctx.Value(SQLWrapperCallerKey),
	}

	logger.Infow(ctx, "after executing query", fields...)

	// TODO: should we ignore the cancelling statement error?
	ignoredErrorMessage := []string{
		"pq: canceling statement due to user request",
		context.Canceled.Error(),
	}

	if err != nil && !lo.Contains(ignoredErrorMessage, err.Error()) {
		logger.Errorw(ctx, "error on executing query", "error", err.Error())
	}
}

func getRequestDuration(ctx context.Context) time.Duration {
	now := time.Now()

	start := ctx.Value(queryStartKey)
	if start == nil {
		return 0
	}

	startTime, ok := start.(time.Time)
	if !ok {
		return 0
	}

	return now.Sub(startTime)
}
