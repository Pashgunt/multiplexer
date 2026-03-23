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
	"transport/internal/infrastructure/db"
	"transport/pkg/logging"
)

type IHttpServer interface {
	Start() error
	Shutdown(ctx context.Context) error
}

type HTTPServer struct {
	server *http.Server
	logger logging.LoggerInterface
}

func NewHTTPServer(
	config appconfig.Config,
	logger logging.LoggerInterface,
	db db.IDB, // todo add registry
) *HTTPServer {
	router := http.NewServeMux()

	server := &HTTPServer{
		server: &http.Server{
			Addr:              config.Environment.Get("PORT"),
			Handler:           router,
			ReadHeaderTimeout: 10 * time.Second,
		},
		logger: logger,
	}

	service := apiservice.NewTargetServiceService(
		repository.NewTargetServiceRepository(db),
		factory.NewTargetServiceFactory(),
	)
	handler := apihandler.NewTargetServiceHandler(service)

	router.HandleFunc("/api/v1/target-services", middleware.Chain(
		handler.Create,
		middleware.LogHandlerMiddleware(server.logger),
		middleware.AllowHTTPMethodMiddleware(http.MethodPost),
	))

	router.HandleFunc(fmt.Sprintf("/api/v1/target-services/{%s}", apiutils.UUID.String()), middleware.Chain(
		handler.Delete,
		middleware.LogHandlerMiddleware(server.logger),
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
