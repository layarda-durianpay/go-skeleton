package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/layarda-durianpay/go-skeleton/internal/disburse/domain/disburse"
	"github.com/layarda-durianpay/go-skeleton/pkg/common/decorator"
	"github.com/layarda-durianpay/go-skeleton/pkg/common/errors"
)

type DisburseParam struct {
	Amount float32
}

type DisburseHandler decorator.CommandHandler[*DisburseParam]

type disburseHandler struct {
	disburseRepo disburse.DisburseRepository
}

func (h disburseHandler) Handle(
	ctx context.Context,
	r *DisburseParam,
) error {
	err := h.disburseRepo.CreateDisbursement(ctx, disburse.Disbursement{
		ID:     uuid.New(),
		Amount: float64(r.Amount),
	})
	if err != nil {
		// always do wrap since we need to keep the stack trace error from the source
		return errors.WrapDpayErrTrace(err)
	}

	return nil
}

func NewDisburseHandler(
	disburseRepo disburse.DisburseRepository,
) DisburseHandler {
	return decorator.ApplyCommandDecorators(
		&disburseHandler{
			disburseRepo,
		},
	)
}
