package sqlwrap

import (
	"context"

	"github.com/layarda-durianpay/go-skeleton/pkg/common/utils"
)

type ctxDBType string

const (
	txKey               ctxDBType = "db-tx-key"
	SQLWrapperCallerKey ctxDBType = "sql-wrapper-caller-key"
	queryStartKey       ctxDBType = "query-start-key"
)

func getRepositoryFnCaller(ctx context.Context) context.Context {
	return context.WithValue(
		ctx,
		SQLWrapperCallerKey,
		utils.GetFnCallerName(2, 2),
	)
}
