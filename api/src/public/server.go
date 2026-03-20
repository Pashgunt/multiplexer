package public

import (
	"context"
	"fmt"
	"net/http"
	"transport/api/src/factory"
	apihandler "transport/api/src/handler"
	"transport/api/src/middleware"
	"transport/api/src/repository"
	apiservice "transport/api/src/service"
	apiutils "transport/api/src/utils"
	appconfig "transport/internal/infrastructure/config/app"
	"transport/pkg/logging"
	"transport/pkg/utils/backoff"
)

type IHttpServer interface {
	Start() error
	Shutdown(ctx context.Context) error
}

type HttpServer struct {
	server *http.Server
	logger logging.LoggerInterface
}

func NewHttpServer(config appconfig.Config) *HttpServer {
	router := http.NewServeMux()
	logger := config.Logger.GetLogger(backoff.ApiLogger)

	server := &HttpServer{
		server: &http.Server{
			Addr:    config.Environment.Get("PORT"),
			Handler: router,
		},
		logger: logger,
	}

	service := apiservice.NewTargetServiceService(
		repository.NewTargetServiceRepository(config.PgSql),
		factory.NewTargetServiceFactory(),
	)
	handler := apihandler.NewTargetServiceHandler(service)

	router.HandleFunc("/api/v1/target-services", middleware.Chain(
		handler.Create,
		middleware.LogHandlerMiddleware(logger),
		middleware.AllowHttpMethodMiddleware(http.MethodPost),
	))

	router.HandleFunc(fmt.Sprintf("/api/v1/target-services/{%s}", apiutils.Uuid), middleware.Chain(
		handler.Delete,
		middleware.LogHandlerMiddleware(logger),
		middleware.AllowHttpMethodMiddleware(http.MethodDelete),
		middleware.UUIDPathParamMiddleware(apiutils.Uuid),
	))

	return server
}

func (s HttpServer) Start() error {
	return s.server.ListenAndServe()
}

func (s HttpServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
