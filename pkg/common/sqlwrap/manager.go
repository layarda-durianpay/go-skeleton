package sqlwrap

import (
	"context"
	"fmt"
	"sync"

	"github.com/durianpay/dpay-common/logger"
	"github.com/layarda-durianpay/go-skeleton/pkg/common/errors"
	"github.com/ztrue/tracerr"
)

type ManagerInterface interface {
	RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type Manager struct {
	db Database
}

var (
	m    *Manager
	once sync.Once
)

// ProvideManager new manager with singleton
func ProvideManager(db Database) *Manager {
	once.Do(func() {
		m = NewManager(db)
	})

	return m
}

// NewManager new manager
func NewManager(db Database) *Manager {
	return &Manager{
		db: db,
	}
}

// RunInTransaction runs the f with the transaction queryable inside the context
func (m *Manager) RunInTransaction(ctx context.Context, f func(ctx context.Context) error) (err error) {
	tx := TransactionFromContext(ctx)

	if tx == nil {
		tx, err = m.db.BeginTx(ctx, nil)
		if err != nil {
			return errors.NewDatabaseError(
				err,
				fmt.Sprintf("error begin transaction: %s", err.Error()),
				errors.DpayInternalError,
			)
		}

		// only assign ctx with tx if new tx was created
		ctx = ContextWithTx(ctx, tx)
	}

	defer func() {
		if r := recover(); r != nil {
			err = tracerr.Errorf("panic error: %v", r)

			if tx != nil {
				errRollback := tx.Rollback()
				if errRollback != nil {
					err = errors.NewDatabaseError(
						errRollback,
						fmt.Sprintf("error on rollback, original panic: %v", r),
						errors.DpayInternalError,
					)
				}
			}
		}
	}()

	err = f(ctx)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return errors.NewDatabaseError(
				err,
				fmt.Sprintf("error on rollback, original error: %v", err.Error()),
				errors.DpayInternalError,
			)
		}

		return errors.WrapDpayErrTrace(err)
	}

	err = tx.Commit()
	if err != nil {
		logger.Errorw(ctx, "error comitting transaction", "error", err)

		errRollback := tx.Rollback()
		if errRollback != nil {
			return errors.NewDatabaseError(
				tracerr.Wrap(errRollback),
				fmt.Sprintf("error on rollback, original error: %v", err.Error()),
				errors.DpayInternalError,
			)
		}

		return errors.NewDatabaseError(
			tracerr.Wrap(err),
			err.Error(),
			errors.DpayInternalError,
		)
	}

	return nil
}
