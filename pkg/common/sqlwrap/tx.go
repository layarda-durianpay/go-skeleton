package sqlwrap

import "context"

// TransactionFromContext returns the transaction object from the context. Return db if not exists
func TransactionFromContext(ctx context.Context) Transaction {
	tx, ok := ctx.Value(txKey).(Transaction)
	if !ok {
		return nil
	}

	return tx
}

// ContextWithTx add database transaction to context
func ContextWithTx(parentContext context.Context, tx Transaction) context.Context {
	return context.WithValue(parentContext, txKey, tx)
}
