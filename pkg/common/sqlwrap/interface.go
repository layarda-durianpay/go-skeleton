package sqlwrap

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type (
	BeforeFunc func(ctx context.Context, query string, args ...interface{}) context.Context
	AfterFunc  func(ctx context.Context, err error, query string, args ...interface{})
)

// Database interface
type Database interface {
	// should implements to use preparer and ext context
	sqlx.Ext
	sqlx.ExtContext

	BeginTx(context.Context, *sql.TxOptions) (Transaction, error)
	Ping() error
	PingContext(ctx context.Context) error
	Close() error
	GetRawDB() *sqlx.DB
}

// Transaction interface
type Transaction interface {
	// should implements to use preparer and ext context
	sqlx.Ext
	sqlx.ExtContext

	Rollback() error
	Commit() error
}
