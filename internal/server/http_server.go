package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	commoncfg "github.com/durianpay/dpay-common/config"
	"github.com/durianpay/dpay-common/constants"
	"github.com/durianpay/dpay-common/dprouter"
	"github.com/durianpay/dpay-common/logger"
	"github.com/durianpay/dpay-common/middleware"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/layarda-durianpay/go-skeleton/internal/disburse/app"
	disbursehttphandler "github.com/layarda-durianpay/go-skeleton/internal/disburse/handler/http"
	"github.com/samber/lo"
)

func startHTTPServer(server *http.Server) error {
	logger.Infof(context.TODO(), "starting API server on %s", server.Addr)

	err := server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func buildHTTPServer(apps *app.Application) *http.Server {
	port := commoncfg.AppPort()
	addr := fmt.Sprintf(":%s", strconv.Itoa(port))

	muxRouter := initRouter(apps)

	headersOk := handlers.AllowedHeaders([]string{constants.ContentType, constants.Authorization, constants.VerificationToken, constants.UserAgent})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	corsHandler := handlers.CORS(headersOk, originsOk, methodsOk)(muxRouter)
	logHandler := middleware.RequestDefaultHandler(corsHandler)
	compressHandler := handlers.CompressHandler(logHandler)

	// metrics middleware
	handler := middleware.PrometheusMiddleware(compressHandler, middleware.HTTPStats{IsLatencyGaugeEnabled: false})

	return &http.Server{
		Addr:    addr,
		Handler: handler,
	}
}

func initRouter(apps *app.Application) *mux.Router {
	router := dprouter.NewDpayRouter(
		middleware.CoreAPIAuthenticator(apps.Dependencies.MerchantGRPCClient.Client),
		middleware.AuthenticateUser,
		nil,
		[]string{
			"/health",
			"/disbursements/disburse",
		},
		middleware.MerchantSnapAPIAuthenticator(apps.Dependencies.MerchantGRPCClient.Client),
	)

	router.EnableTracing("disbursement-service-http")

	for _, route := range getRoutes(apps) {
		router.HandleRoute(route)
	}

	return router.Get()
}

func getRoutes(apps *app.Application) []dprouter.Route {
	disburseServer := disbursehttphandler.ServerInterfaceWrapper{
		Handler: disbursehttphandler.NewHTTPServer(apps),
		ErrorHandlerFunc: func(w http.ResponseWriter, _ *http.Request, err error) {
			if _, ok := lo.ErrorsAs[*disbursehttphandler.InvalidParamFormatError](err); ok {
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
		},
	}

	return []dprouter.Route{
		{
			Path:   "/health",
			Method: http.MethodGet,
			HTTPHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
			Version: "v1",
		},
		{
			Path:        "/disbursements/disburse",
			Method:      http.MethodPost,
			HTTPHandler: http.HandlerFunc(disburseServer.Disburse),
			Version:     "v1",
		},
	}
}
