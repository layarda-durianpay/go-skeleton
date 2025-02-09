package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	commoncfg "github.com/durianpay/dpay-common/config"
	"github.com/durianpay/dpay-common/debugserver"
	"github.com/durianpay/dpay-common/logger"
	"github.com/layarda-durianpay/go-skeleton/internal/config"
	"github.com/layarda-durianpay/go-skeleton/internal/disburse/app"
	"github.com/layarda-durianpay/go-skeleton/internal/disburse/service"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	appObj app.Application

	appObjCleanup = func() error { return nil }

	disbursementCfg config.DisbursementServiceConfig
	globalCfg       config.GlobalConfig
)

func Init() error {
	err := initConfig()
	if err != nil {
		return err
	}

	zapLogger, err := logger.SetupLogger(commoncfg.Env())
	if err != nil {
		panic(err)
	}

	appObj, appObjCleanup = service.NewApplication(zapLogger)
	disbursementCfg = config.ProvideDisbursementConfig()
	globalCfg = config.ProvideGlobalConfig()

	return nil
}

func initConfig() (err error) {
	_ = commoncfg.Load("./", "application")
	// ignore load config since we no need all for now

	return nil
}

func serveHelperServer(ctx context.Context) {
	metricsPort := commoncfg.MetricsPort()
	metricsAddr := fmt.Sprintf(":%s", strconv.Itoa(metricsPort))

	server := http.Server{
		Addr:    metricsAddr,
		Handler: promhttp.Handler(),
	}

	// Serve metrics.
	logger.Infof(context.Background(), "serving metrics at: %s", metricsAddr)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Warnw(context.Background(), "failed to serve metrics", "error", err.Error())
		}
	}()

	go func() {
		<-ctx.Done()

		// ignore error caused by metric server
		_ = server.Shutdown(context.Background())
		logger.Infof(context.Background(), "shutting down Metrics server %s", fmt.Sprintf(":%d", commoncfg.MetricsPort()))
	}()

	// Serve pprof.
	startDebugServer := disbursementCfg.GetStartDebugServer()
	if startDebugServer {
		debugserver.StartDebugServer(ctx, globalCfg.GetDebugPortForServer())
	}
}
