package server

import (
	"context"
	"net"
	"runtime/debug"
	"time"

	commoncfg "github.com/durianpay/dpay-common/config"
	"github.com/durianpay/dpay-common/logger"
	"github.com/durianpay/dpay-common/proto/client"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/layarda-durianpay/go-skeleton/internal/config"
	"github.com/layarda-durianpay/go-skeleton/internal/disburse/app"
	grpchandler "github.com/layarda-durianpay/go-skeleton/internal/disburse/handler/grpc"
	"github.com/layarda-durianpay/go-skeleton/pkg/common/interceptors"
	"github.com/layarda-durianpay/go-skeleton/pkg/common/protogen"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
)

type rpcHandler struct{}

func startGRPCServer(server *grpc.Server) (err error) {
	addr := commoncfg.GRPCAddr()
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	logger.Infof(context.Background(), "starting GRPC server on %s", addr)
	server.Serve(lis)
	return
}

func buildGRPCServer(apps *app.Application) *grpc.Server {
	globalCfg := config.ProvideGlobalConfig()

	// If MaxConnAge is set to 0, the server will have infinite conn age
	kasp := keepalive.ServerParameters{
		MaxConnectionAge: time.Duration(globalCfg.GetGRPCMaxConnectionAge()) * time.Minute,
	}

	panicRecoveryHandler := func(p interface{}) (err error) {
		logger.Errorf(context.Background(), "msg: recovered from panic: %v; stack: %s", p, debug.Stack())
		return status.Errorf(codes.Internal, "%s", p)
	}
	zapLogOpts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
		logging.WithCodes(logging.DefaultErrorToCode),
		logging.WithDurationField(logging.DurationToTimeMillisFields),
	}

	server := grpc.NewServer(
		grpc.KeepaliveParams(kasp),
		grpc.ChainUnaryInterceptor(
			append(
				interceptors.GetDefaultGRPCUnaryInterceptor(
					apps.Dependencies.Logger,
					panicRecoveryHandler,
					zapLogOpts,
				),
				otelgrpc.UnaryServerInterceptor(),
				client.AuthServerInterceptor,
			)...,
		),
		grpc.ChainStreamInterceptor(
			append(
				interceptors.GetDefaultGRPCStreamInterceptor(
					apps.Dependencies.Logger,
					panicRecoveryHandler,
					zapLogOpts,
				),
				otelgrpc.StreamServerInterceptor(),
			)...,
		),
	)

	protogen.RegisterDisbursementServiceServer(server, grpchandler.NewGrpcServer(apps))

	return server
}
