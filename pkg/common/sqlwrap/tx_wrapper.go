package sqlwrap

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

var _ Transaction = &Tx{}

type Tx struct {
	tx     *sqlx.Tx
	before []BeforeFunc
	after  []AfterFunc
}

func (t Tx) DriverName() (res string) {
	return t.tx.DriverName()
}

func (t Tx) Rebind(query string) string {
	return t.tx.Rebind(query)
}

func (t Tx) BindNamed(query string, arg interface{}) (string, []interface{}, error) {
	return t.tx.BindNamed(query, arg)
}

func (t Tx) Exec(query string, args ...interface{}) (res sql.Result, err error) {
	ctx := getRepositoryFnCaller(context.Background())
	ctx = t.doBefore(ctx, query, args...)

	defer t.doAfter(ctx, err, query, args...)

	res, err = t.tx.Exec(query, args...)

	return
}

func (t Tx) Query(query string, args ...interface{}) (res *sql.Rows, err error) {
	ctx := getRepositoryFnCaller(context.Background())
	ctx = t.doBefore(ctx, query, args...)

	defer t.doAfter(ctx, err, query, args...)

	res, err = t.tx.Query(query, args...)

	return
}

func (t Tx) QueryRowx(query string, args ...interface{}) (res *sqlx.Row) {
	ctx := getRepositoryFnCaller(context.Background())
	ctx = t.doBefore(ctx, query, args...)

	defer t.doAfter(ctx, res.Err(), query, args...)

	res = t.tx.QueryRowx(query, args...)

	return
}

func (t Tx) Queryx(query string, args ...interface{}) (res *sqlx.Rows, err error) {
	ctx := getRepositoryFnCaller(context.Background())
	ctx = t.doBefore(ctx, query, args...)

	defer t.doAfter(ctx, err, query, args...)

	res, err = t.tx.Queryx(query, args...)

	return
}

func (t Tx) ExecContext(ctx context.Context, query string, args ...interface{}) (res sql.Result, err error) {
	ctx = getRepositoryFnCaller(ctx)
	ctx = t.doBefore(ctx, query, args...)

	defer t.doAfter(ctx, err, query, args...)

	res, err = t.tx.ExecContext(ctx, query, args...)

	return
}

func (t Tx) QueryContext(ctx context.Context, query string, args ...interface{}) (res *sql.Rows, err error) {
	ctx = getRepositoryFnCaller(ctx)
	ctx = t.doBefore(ctx, query, args...)

	defer t.doAfter(ctx, err, query, args...)

	res, err = t.tx.QueryContext(ctx, query, args...)

	return
}

func (t Tx) QueryRowxContext(ctx context.Context, query string, args ...interface{}) (res *sqlx.Row) {
	ctx = getRepositoryFnCaller(ctx)
	ctx = t.doBefore(ctx, query, args...)

	defer t.doAfter(ctx, res.Err(), query, args...)

	res = t.tx.QueryRowxContext(ctx, query, args...)

	return
}

func (t Tx) QueryxContext(ctx context.Context, query string, args ...interface{}) (res *sqlx.Rows, err error) {
	ctx = getRepositoryFnCaller(ctx)
	ctx = t.doBefore(ctx, query, args...)

	defer t.doAfter(ctx, err, query, args...)

	res, err = t.tx.QueryxContext(ctx, query, args...)

	return
}

func (t Tx) Rollback() error {
	return t.tx.Rollback()
}

func (t Tx) Commit() error {
	return t.tx.Commit()
}

func (t Tx) doBefore(ctx context.Context, query string, args ...interface{}) context.Context {
	for _, f := range t.before {
		ctx = f(ctx, query, args...)
	}

	return ctx
}

func (t Tx) doAfter(ctx context.Context, err error, query string, args ...interface{}) {
	for _, f := range t.after {
		f(ctx, err, query, args...)
	}
}
