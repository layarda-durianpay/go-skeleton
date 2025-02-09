package adapter

import (
	"context"

	"github.com/layarda-durianpay/go-skeleton/internal/disburse/domain/disburse"
	"github.com/layarda-durianpay/go-skeleton/pkg/common/errors"
	"github.com/layarda-durianpay/go-skeleton/pkg/common/sqlwrap"
)

type postgresAgentRepo struct {
	db sqlwrap.Database
}

func (p *postgresAgentRepo) CreateDisbursement(ctx context.Context, disbursement disburse.Disbursement) error {
	qry, args, err := p.db.BindNamed(createDisbursementQuery, disbursementModel{
		ID:     disbursement.ID,
		Amount: disbursement.Amount,
	})
	if err != nil {
		return errors.NewDatabaseError(
			err,
			"failed to bind named for insert query",
			errors.DpayInternalError,
		)
	}

	_, err = p.db.ExecContext(ctx, qry, args...)
	if err != nil {
		return errors.NewDatabaseError(
			err,
			"failed to insert disbursement data",
			errors.DpayInternalError,
		)
	}

	return nil
}

func NewPostgresDisbursementRepository(
	db sqlwrap.Database,
) disburse.DisburseRepository {
	return &postgresAgentRepo{
		db: db,
	}
}
