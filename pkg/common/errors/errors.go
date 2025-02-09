package errors

import (
	"errors"

	"github.com/durianpay/dpay-common/api"
	"github.com/samber/lo"
	"github.com/ztrue/tracerr"
)

type ErrorType string

var (
	ErrorTypeUnknown             = ErrorType("unknown")
	ErrorTypeDatabase            = ErrorType("database")
	ErrorTypeAuthorization       = ErrorType("authorization")
	ErrorTypeIncorrectInput      = ErrorType("incorrect-input")
	ErrorTypeUnprocessableEntity = ErrorType("unprocessable-entity")
	ErrorTypeNotFound            = ErrorType("not-found")
	ErrorTypeForbidden           = ErrorType("forbidden")
	ErrorTypeContextCancelled    = ErrorType("context-cancelled")
)

type ErrorCode string

const (
	DpayInternalError  ErrorCode = ErrorCode("DPAY_INTERNAL_ERROR")
	DpayInvalidRequest ErrorCode = ErrorCode("DPAY_INVALID_REQUEST")
	DpayCancelled      ErrorCode = ErrorCode("DPAY_CANCELLED")
)

// mapClientErrorType mapping the 4xx error as true
var mapClientErrorType = map[ErrorType]bool{
	ErrorTypeAuthorization:       true,
	ErrorTypeIncorrectInput:      true,
	ErrorTypeUnprocessableEntity: true,
	ErrorTypeNotFound:            true,
	ErrorTypeForbidden:           true,
	ErrorTypeUnknown:             false,
	ErrorTypeDatabase:            false,
	ErrorTypeContextCancelled:    true,
}

type DpayError struct {
	err        error
	message    string
	errorCode  ErrorCode
	errorType  ErrorType
	errorInfos []api.ErrorInfo
}

func (s DpayError) Unwrap() error {
	return s.err
}

func (s DpayError) Error() string {
	return s.message
}

func (s DpayError) ErrorCode() ErrorCode {
	return s.errorCode
}

func (s DpayError) ErrorType() ErrorType {
	return s.errorType
}

func (s DpayError) ErrorInfos() []api.ErrorInfo {
	return s.errorInfos
}

func (s DpayError) StackTrace() []tracerr.Frame {
	return tracerr.StackTrace(s.err)
}

// WrapDpayErrTrace wraps an error with trace information using tracerr.Wrap.
// If the provided error is not a DpayError, it will be wrapped into
// a DpayError with an "unknown" error type. This function ensures that the returned
// error is always of type DpayError.
func WrapDpayErrTrace(err error) error {
	if err == nil {
		return nil
	}

	dpayErr, ok := lo.ErrorsAs[DpayError](err)
	if !ok {
		dpayErr = NewDpayError(err, err.Error(), DpayInternalError)
	}

	return NewCustomDpayError(
		tracerr.Wrap(dpayErr.err),
		dpayErr.message,
		dpayErr.errorCode,
		dpayErr.errorType,
	)
}

// GetDPayError will get the DpayError from the err, will return false if the error or the unwrapped is not DpayError type
func GetDPayError(err error) (DpayError, bool) {
	if err == nil {
		return NewDpayError(err, "", ""), false
	}

	dpayErr, ok := lo.ErrorsAs[DpayError](err)
	if ok {
		return dpayErr, true
	}

	return GetDPayError(errors.Unwrap(err))
}

// GetTracerrErr will get the TracerrErr from the err, will return false if the error or the unwrapped is not GetTracerrErr type
func GetTracerrErr(err error) tracerr.Error {
	if err == nil {
		return nil
	}

	tracerrErr, ok := lo.ErrorsAs[tracerr.Error](err)
	if ok {
		return tracerrErr
	}

	return GetTracerrErr(errors.Unwrap(err))
}

// GetOriginalErr will get the original error that's not wrapped by any type error
func GetOriginalErr(err error) error {
	return unwrapError(err)
}

// IsClientError checks if an error is a client error
func IsClientError(err error) bool {
	dpayErr, ok := lo.ErrorsAs[DpayError](err)
	if !ok {
		// The error does not implement the SlugError interface
		return false
	}

	errorType := dpayErr.ErrorType()

	isClientErr, mapped := mapClientErrorType[errorType]
	if !mapped {
		return false
	}

	return isClientErr
}

// unwrapError recursively unwraps errors that implement the Unwrap method.
// If the error does not implement the Unwrap method, it is returned as is.
func unwrapError(err error) error {
	uErr, ok := err.(interface {
		Unwrap() error
	})
	if ok {
		return unwrapError(uErr.Unwrap())
	}

	return err
}
