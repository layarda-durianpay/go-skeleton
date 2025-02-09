package httphandler

import (
	"encoding/json"
	"net/http"

	"github.com/durianpay/dpay-common/api"
	"github.com/durianpay/dpay-common/dcerrors"
	"github.com/layarda-durianpay/go-skeleton/internal/disburse/app"
	"github.com/layarda-durianpay/go-skeleton/internal/disburse/app/command"
	"github.com/layarda-durianpay/go-skeleton/pkg/common/errors"
	"github.com/layarda-durianpay/go-skeleton/pkg/common/httperr"
)

type httpServer struct {
	app *app.Application
}

func NewHTTPServer(apps *app.Application) ServerInterface {
	return &httpServer{
		app: apps,
	}
}

// (POST /disburse)
func (h httpServer) Disburse(w http.ResponseWriter, r *http.Request) {
	var body PostDisburseBody

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		httperr.ResponseWithError(
			errors.NewIncorrectInputError(
				dcerrors.ErrReadingRequestBody,
				dcerrors.ErrReadingRequestBody.Error(),
				dcerrors.DpayInvalidRequest,
			),
			w, r,
		)
		return
	}

	err = h.app.Commands.Disburse.Handle(r.Context(), &command.DisburseParam{
		Amount: body.Amount,
	})
	if err != nil {
		httperr.ResponseWithError(err, w, r)
		return
	}

	api.RespondWithJSON(w, http.StatusCreated, api.Response{
		Message: "Success process your disburse.",
	})
}
