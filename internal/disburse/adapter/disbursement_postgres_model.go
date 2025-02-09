package adapter

import "github.com/google/uuid"

type disbursementModel struct {
	ID     uuid.UUID `db:"id"`
	Amount float64   `db:"amount"`
}
