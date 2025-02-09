package errors

import (
	"github.com/durianpay/dpay-common/api"
	"github.com/samber/lo"
	"github.com/ztrue/tracerr"
)

func NewCustomDpayError(
	err error,
	message string,
	errCode ErrorCode,
	errType ErrorType,
) DpayError {
	if _, ok := lo.ErrorsAs[tracerr.Error](err); !ok {
		err = tracerr.Wrap(err)
	}

	return DpayError{
		err:       err,
		message:   message,
		errorCode: errCode,
		errorType: errType,
	}
}

func NewDpayError(err error, message string, errCode ErrorCode) DpayError {
	if _, ok := lo.ErrorsAs[tracerr.Error](err); !ok {
		err = tracerr.Wrap(err)
	}

	return DpayError{
		err:       err,
		message:   message,
		errorCode: errCode,
		errorType: ErrorTypeUnknown,
	}
}

func NewAuthorizationError(err error, message string, errCode ErrorCode) DpayError {
	if _, ok := lo.ErrorsAs[tracerr.Error](err); !ok {
		err = tracerr.Wrap(err)
	}

	return DpayError{
		err:       err,
		message:   message,
		errorCode: errCode,
		errorType: ErrorTypeAuthorization,
	}
}

func NewIncorrectInputError(
	err error,
	message string,
	errCode ErrorCode,
) DpayError {
	if _, ok := lo.ErrorsAs[tracerr.Error](err); !ok {
		err = tracerr.Wrap(err)
	}

	return DpayError{
		err:       err,
		message:   message,
		errorCode: errCode,
		errorType: ErrorTypeIncorrectInput,
	}
}

func NewUnprocessableEntityError(
	err error,
	message string,
	errCode ErrorCode,
	errInfos ...api.ErrorInfo,
) DpayError {
	if _, ok := lo.ErrorsAs[tracerr.Error](err); !ok {
		err = tracerr.Wrap(err)
	}

	return DpayError{
		err:        err,
		message:    message,
		errorCode:  errCode,
		errorType:  ErrorTypeUnprocessableEntity,
		errorInfos: errInfos,
	}
}

func NewNotFoundError(err error, message string, errCode ErrorCode) DpayError {
	if _, ok := lo.ErrorsAs[tracerr.Error](err); !ok {
		err = tracerr.Wrap(err)
	}

	return DpayError{
		err:       err,
		message:   message,
		errorCode: errCode,
		errorType: ErrorTypeNotFound,
	}
}

func NewDatabaseError(err error, message string, errCode ErrorCode) DpayError {
	if _, ok := lo.ErrorsAs[tracerr.Error](err); !ok {
		err = tracerr.Wrap(err)
	}

	return DpayError{
		err:       err,
		message:   message,
		errorCode: errCode,
		errorType: ErrorTypeDatabase,
	}
}

func NewForbiddenError(err error, message string, errCode ErrorCode) DpayError {
	if _, ok := lo.ErrorsAs[tracerr.Error](err); !ok {
		err = tracerr.Wrap(err)
	}

	return DpayError{
		err:       err,
		message:   message,
		errorCode: errCode,
		errorType: ErrorTypeForbidden,
	}
}

func NewContextCancelledError(err error, message string, errCode ErrorCode) DpayError {
	if _, ok := lo.ErrorsAs[tracerr.Error](err); !ok {
		err = tracerr.Wrap(err)
	}

	return DpayError{
		err:       err,
		message:   message,
		errorCode: errCode,
		errorType: ErrorTypeContextCancelled,
	}
}
