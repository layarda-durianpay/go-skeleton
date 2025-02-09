package sqlwrap

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// check DB implements Database interface
var _ Database = &DB{}

type DB struct {
	DB     *sqlx.DB
	before []BeforeFunc
	after  []AfterFunc
}

func NewDB(db *sqlx.DB) *DB {
	return &DB{
		DB:     db,
		before: make([]BeforeFunc, 0),
		after:  make([]AfterFunc, 0),
	}
}

func (db *DB) GetRawDB() *sqlx.DB {
	return db.DB
}

func (db *DB) AddBeforeFunc(f BeforeFunc) {
	db.before = append(db.before, f)
}

func (db *DB) AddAfterFunc(f AfterFunc) {
	db.after = append(db.after, f)
}

func (db DB) DriverName() (res string) {
	return db.DB.DriverName()
}
func (db DB) Rebind(query string) string {
	return db.DB.Rebind(query)
}
func (db DB) BindNamed(query string, arg interface{}) (string, []interface{}, error) {
	return db.DB.BindNamed(query, arg)
}

func (db DB) Exec(query string, args ...interface{}) (res sql.Result, err error) {
	ctx := getRepositoryFnCaller(context.Background())
	ctx = db.doBefore(ctx, query, args...)

	defer db.doAfter(ctx, err, query, args...)

	res, err = db.DB.Exec(query, args...)

	return
}

func (db DB) Query(query string, args ...interface{}) (res *sql.Rows, err error) {
	ctx := getRepositoryFnCaller(context.Background())
	ctx = db.doBefore(ctx, query, args...)

	defer db.doAfter(ctx, err, query, args...)

	res, err = db.DB.Query(query, args...)

	return
}

func (db DB) QueryRowx(query string, args ...interface{}) (res *sqlx.Row) {
	ctx := getRepositoryFnCaller(context.Background())
	ctx = db.doBefore(ctx, query, args...)

	defer db.doAfter(ctx, res.Err(), query, args...)

	res = db.DB.QueryRowx(query, args...)

	return
}

func (db DB) Queryx(query string, args ...interface{}) (res *sqlx.Rows, err error) {
	ctx := getRepositoryFnCaller(context.Background())
	ctx = db.doBefore(ctx, query, args...)

	defer db.doAfter(ctx, err, query, args...)

	res, err = db.DB.Queryx(query, args...)

	return
}

func (db DB) ExecContext(ctx context.Context, query string, args ...interface{}) (res sql.Result, err error) {
	ctx = getRepositoryFnCaller(ctx)
	ctx = db.doBefore(ctx, query, args...)

	defer db.doAfter(ctx, err, query, args...)

	res, err = db.DB.ExecContext(ctx, query, args...)

	return
}

func (db DB) QueryContext(ctx context.Context, query string, args ...interface{}) (res *sql.Rows, err error) {
	ctx = getRepositoryFnCaller(ctx)
	ctx = db.doBefore(ctx, query, args...)

	defer db.doAfter(ctx, err, query, args...)

	res, err = db.DB.QueryContext(ctx, query, args...)

	return
}

func (db DB) QueryRowxContext(ctx context.Context, query string, args ...interface{}) (res *sqlx.Row) {
	ctx = getRepositoryFnCaller(ctx)
	ctx = db.doBefore(ctx, query, args...)

	defer db.doAfter(ctx, res.Err(), query, args...)

	res = db.DB.QueryRowxContext(ctx, query, args...)

	return
}

func (db DB) QueryxContext(ctx context.Context, query string, args ...interface{}) (res *sqlx.Rows, err error) {
	ctx = getRepositoryFnCaller(ctx)
	ctx = db.doBefore(ctx, query, args...)

	defer db.doAfter(ctx, err, query, args...)

	res, err = db.DB.QueryxContext(ctx, query, args...)

	return
}

func (db *DB) BeginTx(ctx context.Context, options *sql.TxOptions) (Transaction, error) {
	tx, err := db.DB.BeginTxx(ctx, options)
	if err != nil {
		return nil, err
	}

	return &Tx{
		tx:     tx,
		before: db.before,
		after:  db.after,
	}, nil
}

func (db *DB) Ping() (err error) {
	return db.DB.Ping()
}

func (db *DB) PingContext(ctx context.Context) error {
	return db.DB.PingContext(ctx)
}

func (db *DB) Close() error {
	return db.DB.Close()
}

func (db DB) doBefore(ctx context.Context, query string, args ...interface{}) context.Context {
	for _, f := range db.before {
		ctx = f(ctx, query, args...)
	}

	return ctx
}

func (db DB) doAfter(ctx context.Context, err error, query string, args ...interface{}) {
	for _, f := range db.after {
		f(ctx, err, query, args...)
	}
}
