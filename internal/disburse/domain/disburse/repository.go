package disburse

import "context"

type DisburseRepository interface {
	CreateDisbursement(ctx context.Context, disbursement Disbursement) error
}
