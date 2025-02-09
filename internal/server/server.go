package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	commoncfg "github.com/durianpay/dpay-common/config"
	"github.com/durianpay/dpay-common/logger"
	"github.com/layarda-durianpay/go-skeleton/internal/config"
	"github.com/layarda-durianpay/go-skeleton/internal/disburse/app"
	"github.com/layarda-durianpay/go-skeleton/pkg/common/opentelemetry"
	"golang.org/x/sync/errgroup"
)

const servicesCount = 1

func Start() error {
	mainCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	disbursementCfg := config.ProvideDisbursementConfig()

	otelCleanup, err := opentelemetry.InitOTelTrace(mainCtx, disbursementCfg.GetEnableConfigOpenTelemetry())
	defer otelCleanup()

	if err != nil {
		return err
	}

	wg := new(sync.WaitGroup)
	wg.Add(servicesCount)

	go func() {
		startServers(mainCtx, appObj)
		wg.Done()
	}()

	wg.Wait()

	return appObjCleanup()
}

func startServers(ctx context.Context, apps app.Application) (err error) {
	// will serve metrics and pprop server
	serveHelperServer(ctx)

	httpServer := buildHTTPServer(&apps)
	grpcServer := buildGRPCServer(&apps)

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		err := startHTTPServer(httpServer)
		if err != nil {
			logger.Infof(context.Background(), "http server stopped - %s", err.Error())
		}
		// setting this nil as we are already panicking if server does not start
		return nil
	})

	g.Go(func() error {
		err := startGRPCServer(grpcServer)
		if err != nil {
			logger.Infof(context.Background(), "grpc server stopped - %s", err.Error())
		}
		// setting this nil as we are already panicking if server does not start
		return nil
	})

	g.Go(func() (err error) {
		<-gCtx.Done()
		err = httpServer.Shutdown(context.Background())
		logger.Infof(context.Background(), "shutting down API server %s", fmt.Sprintf(":%d", commoncfg.AppPort()))

		logger.Infof(context.Background(), "shutting down GRPC server %s", commoncfg.GRPCAddr())
		grpcServer.GracefulStop()

		return err
	})

	if err = g.Wait(); err != nil {
		logger.Infof(context.Background(), "shutting down servers %s", err.Error())
	}

	return
}
