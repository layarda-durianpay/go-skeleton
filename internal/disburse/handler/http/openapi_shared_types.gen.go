// Package httphandler provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.2.0 DO NOT EDIT.
package httphandler

// BadRequestError defines model for BadRequestError.
type BadRequestError struct {
	// Error message error description
	Error string `json:"error"`

	// ErrorCode durianpay error code
	ErrorCode *string `json:"error_code,omitempty"`

	// RequestId durianpay request_id for reconciliation
	RequestId *string `json:"request_id,omitempty"`
}

// Error defines model for Error.
type Error struct {
	// Error message error description
	Error string `json:"error"`

	// ErrorCode durianpay error code
	ErrorCode *string `json:"error_code,omitempty"`

	// RequestId durianpay request_id for reconciliation
	RequestId *string `json:"request_id,omitempty"`
}

// NotFoundError defines model for NotFoundError.
type NotFoundError struct {
	// Error message error description
	Error string `json:"error"`

	// ErrorCode durianpay error code
	ErrorCode *string `json:"error_code,omitempty"`

	// RequestId durianpay request_id for reconciliation
	RequestId *string `json:"request_id,omitempty"`
}

// BadRequestResponse defines model for BadRequestResponse.
type BadRequestResponse = BadRequestError

// NotFoundRequest defines model for NotFoundRequest.
type NotFoundRequest = NotFoundError

// UnexpectedErrorRequest defines model for UnexpectedErrorRequest.
type UnexpectedErrorRequest = Error
