package httperr

import (
	"encoding/json"
	"net/http"

	"github.com/durianpay/dpay-common/api"
	"github.com/durianpay/dpay-common/logger"
	"github.com/layarda-durianpay/go-skeleton/pkg/common/errors"
	"github.com/samber/lo"
)

// errorMap only map error caused by 4xx, don't cover the 5xx error, let it became 500 Server Internal Error
var errorMap = map[errors.ErrorType]int{
	errors.ErrorTypeAuthorization:       http.StatusUnauthorized,
	errors.ErrorTypeIncorrectInput:      http.StatusBadRequest,
	errors.ErrorTypeUnprocessableEntity: http.StatusUnprocessableEntity,
	errors.ErrorTypeNotFound:            http.StatusNotFound,
	errors.ErrorTypeForbidden:           http.StatusForbidden,
	errors.ErrorTypeContextCancelled:    499, // client close connection

}

func ResponseWithError(err error, w http.ResponseWriter, r *http.Request) {
	dpayErr, ok := errors.GetDPayError(err)
	if !ok {
		httpResponseWithInternalServerError(err, w, r)
		return
	}

	status, ok := errorMap[dpayErr.ErrorType()]
	if !ok {
		// if slug error type not exist, better indicates as Internal Server error
		// the slug and message should be private to internal only for 5xx error.
		httpResponseWithInternalServerError(err, w, r)

		return
	}

	httpRespondWithError(err, w, r, dpayErr.Error(), dpayErr.ErrorCode(), status)

}

func httpResponseWithInternalServerError(err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(
		err,
		w,
		r,
		"Internal Server error",
		errors.DpayInternalError,
		http.StatusInternalServerError,
	)
}

func httpRespondWithError(
	err error,
	w http.ResponseWriter,
	r *http.Request,
	message string,
	errCode errors.ErrorCode,
	statusCode int,
) {
	logFields := make([]any, 0)

	resp := errorResponse{
		httpStatus: statusCode,

		Error:     message,
		ErrorCode: string(errCode),
	}

	dpayErr, ok := lo.ErrorsAs[errors.DpayError](err)
	if ok {
		resp.Errors = dpayErr.ErrorInfos()
	}

	logFields = append(logFields, "error_original", dpayErr)
	logFields = append(logFields, "error_code", dpayErr.ErrorCode())
	logFields = append(logFields, "status", http.StatusText(statusCode))
	logFields = append(logFields, "status_code", statusCode)

	loggerUsed := logger.Warnw

	if statusCode > 499 {
		loggerUsed = logger.Errorw
	}

	loggerUsed(r.Context(), "HTTP Respond With Error", logFields...)

	w.Header().Set("Content-Type", "application/json")

	if err := resp.Render(w, r); err != nil {
		panic(err)
	}
}

type errorResponse struct {
	httpStatus int

	Error     string `json:"error,omitempty"`
	ErrorCode string `json:"error_code,omitempty"`
	RequestID string `json:"request_id,omitempty"`

	// can be used when returning multiple form errors
	Errors []api.ErrorInfo `json:"errors,omitempty"`
}

func (e errorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(e.httpStatus)

	return json.NewEncoder(w).Encode(e)
}
