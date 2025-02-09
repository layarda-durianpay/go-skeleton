package disburse

import "github.com/google/uuid"

type Disbursement struct {
	ID     uuid.UUID
	Amount float64
}
