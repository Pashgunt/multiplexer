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
)

type HttpServer struct {
	server *http.Server
}

func NewHttpServer(config appconfig.Config) *HttpServer {
	router := http.NewServeMux()

	service := apiservice.NewTargetServiceService(
		repository.NewTargetServiceRepository(config.PgSql),
		factory.NewTargetServiceFactory(),
	)
	handler := apihandler.NewTargetServiceHandler(service)

	//todo add log middleware
	router.HandleFunc("/api/v1/target-services", middleware.Chain(
		handler.Create,
		middleware.AllowHttpMethodMiddleware(http.MethodPost),
	))

	router.HandleFunc(fmt.Sprintf("/api/v1/target-services/{%s}", apiutils.Uuid), middleware.Chain(
		handler.Delete,
		middleware.AllowHttpMethodMiddleware(http.MethodDelete),
		middleware.UUIDPathParamMiddleware(apiutils.Uuid),
	))

	return &HttpServer{
		server: &http.Server{
			Addr:    config.Environment.Get("PORT"),
			Handler: router,
		},
	}
}

func (s HttpServer) Start() error {
	return s.server.ListenAndServe()
}

func (s HttpServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
