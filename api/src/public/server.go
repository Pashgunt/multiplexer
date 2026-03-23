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
	HandleFunc(db db.IDB)
}

type HTTPServer struct {
	server *http.Server
	logger logging.LoggerInterface
	router *http.ServeMux
}

func NewHTTPServer(
	config appconfig.Config,
	logger logging.LoggerInterface,
) IHttpServer {
	router := http.NewServeMux()

	return &HTTPServer{
		server: &http.Server{
			Addr:              config.Environment.Get("PORT"),
			Handler:           router,
			ReadHeaderTimeout: 10 * time.Second,
		},
		logger: logger,
		router: router,
	}
}

func (s HTTPServer) HandleFunc(
	db db.IDB, //todo add DI
) {
	service := apiservice.NewTargetServiceService(
		repository.NewTargetServiceRepository(db),
		factory.NewTargetServiceFactory(),
	)
	handler := apihandler.NewTargetServiceHandler(service)

	s.router.HandleFunc("/api/v1/target-services", middleware.Chain(
		handler.Create,
		middleware.LogHandlerMiddleware(s.logger),
		middleware.AllowHTTPMethodMiddleware(http.MethodPost),
	))

	s.router.HandleFunc(fmt.Sprintf("/api/v1/target-services/{%s}", apiutils.UUID.String()), middleware.Chain(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				handler.Get(w, r)
			case http.MethodDelete:
				handler.Delete(w, r)
			default:
				http.Error(
					w,
					fmt.Sprintf("Method %s not allowed", r.Method),
					http.StatusMethodNotAllowed,
				)
			}
		},
		middleware.LogHandlerMiddleware(s.logger),
		middleware.UUIDPathParamMiddleware(apiutils.UUID),
	))
}

func (s HTTPServer) Start() error {
	return s.server.ListenAndServe()
}

func (s HTTPServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
