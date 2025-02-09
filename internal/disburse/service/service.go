package service

import (
	"context"
	"log"

	"github.com/durianpay/dpay-common/logger"
	"github.com/durianpay/dpay-common/proto/client"
	"github.com/layarda-durianpay/go-skeleton/internal/config"
	"github.com/layarda-durianpay/go-skeleton/internal/disburse/adapter"
	"github.com/layarda-durianpay/go-skeleton/internal/disburse/app"
	"github.com/layarda-durianpay/go-skeleton/internal/disburse/app/command"
	"github.com/layarda-durianpay/go-skeleton/internal/disburse/domain/disburse"
	"github.com/layarda-durianpay/go-skeleton/pkg/common/sqlwrap"
	"go.uber.org/zap"
)

const serviceName = "disbursement_service"

type closeFn func() error

func NewApplication(zapLogger *zap.SugaredLogger) (app.Application, closeFn) {
	var (
		disbursementConf = config.ProvideDisbursementConfig()
		globalConf       = config.ProvideGlobalConfig()
		db               = sqlwrap.ProvidePostgres(
			sqlwrap.ServiceNameOption(serviceName),
			sqlwrap.WithOTelOption(disbursementConf.GetEnableConfigOpenTelemetry()),
		)
	)

	merchantGRPCClient, err := client.InitMerchantClient(globalConf.GetMerchantServiceGRPCAddr())
	if err != nil {
		logger.Errorw(context.Background(), "error initializing merchant grpc client", "error", err.Error())
		panic(err)
	}

	// repository
	disburseRepo := adapter.NewPostgresDisbursementRepository(db)

	return newApplication(
			db,
			zapLogger,
			merchantGRPCClient,
			disburseRepo,
		), close(
			db,
			zapLogger,
			merchantGRPCClient,
		)
}

func newApplication(
	db sqlwrap.Database,
	logger *zap.SugaredLogger,

	// grpc related
	merchantService client.MerchantServiceClient,

	// repo related
	disburseRepository disburse.DisburseRepository,
) app.Application {
	return app.Application{
		Dependencies: app.Dependencies{
			DB:                 db,
			MerchantGRPCClient: &merchantService,
			Logger:             logger,
		},
		Commands: app.Commands{
			Disburse: command.NewDisburseHandler(disburseRepository),
		},
		Queries: app.Queries{},
	}
}

func close(
	db sqlwrap.Database,
	zapLogger *zap.SugaredLogger,

	// grpc related
	merchantService client.MerchantServiceClient,
) closeFn {
	return func() (err error) {
		var errs = make(map[string]error)

		if zapLogger != nil {
			err = zapLogger.Sync()
			if err != nil {
				errs["failed to flush log"] = err
			}
		}

		err = db.Close()
		if err != nil {
			errs["failed to close db"] = err
		}

		merchantService.Close()

		if zapLogger != nil {
			for msg, err := range errs {
				logger.Errorw(context.TODO(), msg, "error", err)
			}

			return err
		}

		for msg, err := range errs {
			log.Printf("Warn: %s, got error: %s"+msg, err.Error())
		}

		return err
	}
}
