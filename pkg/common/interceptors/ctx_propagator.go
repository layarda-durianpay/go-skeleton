package interceptors

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/durianpay/dpay-common/constants"
	"github.com/layarda-durianpay/go-skeleton/pkg/common/utils"
)

// listContextPropagation what context want to be propagate to grpc metadata
var listContextPropagation = map[string]constants.ContextKey{
	string(constants.RequestIDKey): constants.RequestIDKey,
}

func ContextPropagationUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return grpc.UnaryServerInterceptor(
		func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ any, err error) {
			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				return handler(ctx, req)
			}

			for key, val := range md {
				ctxKey, ok := listContextPropagation[key]
				if !ok {
					continue
				}

				if len(val) != 1 {
					continue
				}

				ctx = context.WithValue(ctx, ctxKey, val[0])
			}

			return handler(ctx, req)
		},
	)
}

func ContextPropagationStreamServerInterceptor() grpc.StreamServerInterceptor {
	return grpc.StreamServerInterceptor(
		func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
			ctx := stream.Context()

			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				return handler(srv, stream)
			}

			for key, val := range md {
				ctxKey, ok := listContextPropagation[key]
				if !ok {
					continue
				}

				if len(val) != 1 {
					continue
				}

				ctx = context.WithValue(ctx, ctxKey, val[0])
			}

			wrappedStream := utils.NewServerStreamWrapper(stream, ctx)

			return handler(srv, wrappedStream)
		},
	)
}

func ContextPropagationUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return grpc.UnaryClientInterceptor(
		func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			md := make(metadata.MD)

			for mdKey, ctxKey := range listContextPropagation {
				val, ok := (ctx.Value(ctxKey)).(string)
				if val == "" || !ok {
					continue
				}

				md.Append(mdKey, val)
			}

			if len(md) > 0 {
				ctx = metadata.NewOutgoingContext(ctx, md)
			}

			return invoker(ctx, method, req, reply, cc, opts...)
		},
	)
}

func ContextPropagationStreamClientInterceptor() grpc.StreamClientInterceptor {
	return grpc.StreamClientInterceptor(
		func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
			md := make(metadata.MD)

			for mdKey, ctxKey := range listContextPropagation {
				val, ok := (ctx.Value(ctxKey)).(string)
				if val == "" || !ok {
					continue
				}

				md.Append(mdKey, val)
			}

			if len(md) > 0 {
				ctx = metadata.NewOutgoingContext(ctx, md)
			}

			return streamer(ctx, desc, cc, method, opts...)
		},
	)
}
