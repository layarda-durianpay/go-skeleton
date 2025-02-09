package interceptors

import (
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	loginterceptor "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func GetDefaultGRPCUnaryInterceptor(
	zapLogger *zap.SugaredLogger,
	panicRecoverHandler func(p interface{}) (err error),
	logOptions []loginterceptor.Option,
) []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		ContextPropagationUnaryServerInterceptor(),
		grpc_ctxtags.UnaryServerInterceptor(
			grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor),
		),
		grpc_recovery.UnaryServerInterceptor(
			grpc_recovery.WithRecoveryHandler(panicRecoverHandler),
		),
		loginterceptor.UnaryServerInterceptor(interceptorLogger(zapLogger), logOptions...),
	}
}

func GetDefaultGRPCStreamInterceptor(
	zapLogger *zap.SugaredLogger,
	panicRecoverHandler func(p interface{}) (err error),
	logOptions []loginterceptor.Option,
) []grpc.StreamServerInterceptor {
	return []grpc.StreamServerInterceptor{
		ContextPropagationStreamServerInterceptor(),
		grpc_ctxtags.StreamServerInterceptor(
			grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor),
		),
		grpc_recovery.StreamServerInterceptor(
			grpc_recovery.WithRecoveryHandler(panicRecoverHandler),
		),
		loginterceptor.StreamServerInterceptor(interceptorLogger(zapLogger), logOptions...),
	}
}
