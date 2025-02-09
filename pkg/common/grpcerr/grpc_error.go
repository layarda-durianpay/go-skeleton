package grpcerr

import (
	"github.com/layarda-durianpay/go-skeleton/pkg/common/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var errorCodeMap = map[errors.ErrorType]codes.Code{
	errors.ErrorTypeAuthorization:       codes.Unauthenticated,
	errors.ErrorTypeIncorrectInput:      codes.InvalidArgument,
	errors.ErrorTypeUnprocessableEntity: codes.InvalidArgument,
	errors.ErrorTypeNotFound:            codes.NotFound,
	errors.ErrorTypeForbidden:           codes.PermissionDenied,
	errors.ErrorTypeContextCancelled:    codes.Canceled,
	errors.ErrorTypeUnknown:             codes.Internal,
}

var codeErrorMap = map[codes.Code]errors.ErrorType{
	codes.Unauthenticated:   errors.ErrorTypeAuthorization,
	codes.InvalidArgument:   errors.ErrorTypeIncorrectInput,
	codes.NotFound:          errors.ErrorTypeNotFound,
	codes.ResourceExhausted: errors.ErrorTypeForbidden,
	codes.Canceled:          errors.ErrorTypeContextCancelled,
	codes.Internal:          errors.ErrorTypeUnknown,
}

// TransformToGRPCErr transform error to GRPC status.Error
func TransformToGRPCErr(err error) error {
	dpayErr, ok := errors.GetDPayError(err)
	if !ok {
		return status.Error(
			codes.Internal, err.Error(),
		)
	}

	code, ok := errorCodeMap[dpayErr.ErrorType()]
	if !ok {
		return status.Error(
			codes.Internal, err.Error(),
		)
	}

	return status.Error(code, err.Error())
}

func GetErrorType(err error) errors.ErrorType {
	grpcStatus, ok := status.FromError(err)
	if !ok {
		return errors.ErrorTypeUnknown
	}

	errType, ok := codeErrorMap[grpcStatus.Code()]
	if !ok {
		return errors.ErrorTypeUnknown
	}

	return errType
}

func GetErrorMessage(err error) string {
	grpcStatus, ok := status.FromError(err)
	if !ok {
		return err.Error()
	}

	return grpcStatus.Message()
}
