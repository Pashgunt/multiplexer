package public

import (
	"context"
	"fmt"
	"net/http"
	"time"
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

type HTTPServer struct {
	server *http.Server
	logger logging.LoggerInterface
}

func NewHTTPServer(config appconfig.Config) *HTTPServer {
	router := http.NewServeMux()
	logger := config.Logger.GetLogger(backoff.APILogger)

	server := &HTTPServer{
		server: &http.Server{
			Addr:              config.Environment.Get("PORT"),
			Handler:           router,
			ReadHeaderTimeout: 10 * time.Second,
		},
		logger: logger,
	}

	service := apiservice.NewTargetServiceService(
		repository.NewTargetServiceRepository(config.PgSQL),
		factory.NewTargetServiceFactory(),
	)
	handler := apihandler.NewTargetServiceHandler(service)

	router.HandleFunc("/api/v1/target-services", middleware.Chain(
		handler.Create,
		middleware.LogHandlerMiddleware(logger),
		middleware.AllowHTTPMethodMiddleware(http.MethodPost),
	))

	router.HandleFunc(fmt.Sprintf("/api/v1/target-services/{%s}", apiutils.UUID.String()), middleware.Chain(
		handler.Delete,
		middleware.LogHandlerMiddleware(logger),
		middleware.AllowHTTPMethodMiddleware(http.MethodDelete),
		middleware.UUIDPathParamMiddleware(apiutils.UUID),
	))

	return server
}

func (s HTTPServer) Start() error {
	return s.server.ListenAndServe()
}

func (s HTTPServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
