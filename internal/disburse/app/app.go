package app

import (
	"github.com/durianpay/dpay-common/proto/client"
	"github.com/layarda-durianpay/go-skeleton/internal/disburse/app/command"
	"github.com/layarda-durianpay/go-skeleton/pkg/common/sqlwrap"
	"go.uber.org/zap"
)

type Application struct {
	Dependencies Dependencies
	Commands     Commands
	Queries      Queries
}

type Dependencies struct {
	DB                 sqlwrap.Database
	Logger             *zap.SugaredLogger
	MerchantGRPCClient *client.MerchantServiceClient
}

type Commands struct {
	Disburse command.DisburseHandler
}

type Queries struct {
}
