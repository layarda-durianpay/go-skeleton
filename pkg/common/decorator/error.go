package decorator

import (
	"context"
	"fmt"

	"github.com/layarda-durianpay/go-skeleton/pkg/common/errors"
	"github.com/samber/lo"
	"google.golang.org/grpc/status"
)

type commandErrorDecorator[C any] struct {
	base CommandHandler[C]
}

func (d commandErrorDecorator[C]) Handle(ctx context.Context, cmd C) error {
	err := d.base.Handle(ctx, cmd)
	if err == nil {
		return nil
	}

	isContextCancelledError, transformedErr := transformContextCancelledError(err)
	if isContextCancelledError {
		return transformedErr
	}

	return err
}

type queryErrorDecorator[Q any, R any] struct {
	base QueryHandler[Q, R]
}

func (d queryErrorDecorator[Q, R]) Handle(ctx context.Context, query Q) (R, error) {
	result, err := d.base.Handle(ctx, query)
	if err == nil {
		return result, nil
	}

	isContextCancelledError, transformedErr := transformContextCancelledError(err)
	if isContextCancelledError {
		return result, transformedErr
	}

	return result, err
}

func transformContextCancelledError(origErr error) (bool, error) {
	if lo.IsNil(origErr) {
		return false, origErr
	}

	checkedErr := getOriginalErrMessageFromGrpcError(
		getUnwrappedError(origErr), // get origError from Wrapped Error interface
	) // get origError from grpc error

	if lo.IsNil(checkedErr) {
		return false, origErr
	}

	validContextCancelledMessages := []string{
		"pq: canceling statement due to user request",
		context.Canceled.Error(),
	}

	if lo.Contains(
		validContextCancelledMessages,
		checkedErr.Error(),
	) {
		return true, errors.NewContextCancelledError(
			checkedErr,
			checkedErr.Error(),
			errors.DpayCancelled,
		)
	}

	return false, origErr
}

func getUnwrappedError(err error) error {
	if lo.IsNil(err) {
		return nil
	}

	checkedErr := errors.GetOriginalErr(err)

	if lo.IsNil(checkedErr) {
		// DpayError might used like this
		// errors.NewDpayError(nil, "some description", "some-error-code")
		// so we need to prevent if slug error wrapped nil error
		return err
	}

	return checkedErr
}

func getOriginalErrMessageFromGrpcError(err error) error {
	if lo.IsNil(err) {
		return nil
	}

	if grpcStatus, ok := status.FromError(err); ok {
		// since grpc.Error wrapped into something like this
		// fmt.Sprintf("rpc error: code = %s desc = %s", s.Code(), s.Message())
		// we need to get the original error description through status.Message()
		return fmt.Errorf("%s", grpcStatus.Message())
	}

	return err
}
