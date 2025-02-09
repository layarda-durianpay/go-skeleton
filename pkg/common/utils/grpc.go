package utils

import (
	"context"

	"google.golang.org/grpc"
)

// ServerStreamWrapper is a wrapper for gRPC ServerStream that allows modification of the associated context.
//
// This type was inspired by a discussion on the gRPC Google Group:
// https://groups.google.com/g/grpc-io/c/Q88GQFTPF1o
type ServerStreamWrapper struct {
	grpc.ServerStream
	ctx      context.Context
	response interface{}
}

// NewServerStreamWrapper creates a new instance of ServerStreamWrapper wrapping new context.
func NewServerStreamWrapper(
	ss grpc.ServerStream,
	newCtx context.Context,
) *ServerStreamWrapper {
	return &ServerStreamWrapper{
		ServerStream: ss,
		ctx:          newCtx,
	}
}

// Context returns the modified context associated with the ServerStreamWrapper.
func (w *ServerStreamWrapper) Context() context.Context { return w.ctx }

func (w *ServerStreamWrapper) SendMsg(m any) error {
	w.response = m
	return w.ServerStream.SendMsg(m)
}

func (w *ServerStreamWrapper) GetResponse() interface{} {
	if w.response == nil {
		// No body to decode
		return nil
	}

	return unwrapInterfacePointer(w.response)
}

func unwrapInterfacePointer(data interface{}) interface{} {
	if ptr, ok := data.(**interface{}); ok {
		return *ptr
	}

	return data
}
