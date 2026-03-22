package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	apiutils "transport/api/src/utils"
	logging2 "transport/pkg/logging"

	"github.com/google/uuid"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Chain(h http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}

	return h
}

func LogHandlerMiddleware(logger logging2.LoggerInterface) Middleware {
	return func(_ http.HandlerFunc) http.HandlerFunc {
		return func(_ http.ResponseWriter, r *http.Request) {
			logger.Info(logging2.NewAPILogEntity(fmt.Sprintf("Handled HTTP Request: %s %s", r.Method, r.URL.Path)))
		}
	}
}

func AllowHTTPMethodMiddleware(method string) Middleware {
	return func(handlerFunc http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method == method {
				handlerFunc(w, r)

				return
			}

			w.Header().Set("Allow", method)

			http.Error(
				w,
				fmt.Sprintf("Method %s not allowed", r.Method),
				http.StatusMethodNotAllowed,
			)
		}
	}
}

func UUIDPathParamMiddleware(pathParam apiutils.PathParam) Middleware {
	return func(handlerFunc http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			parts := strings.Split(strings.TrimSuffix(r.URL.Path, "/"), "/")
			pathParamStr := parts[len(parts)-1]

			if _, err := uuid.Parse(pathParamStr); err == nil {
				handlerFunc.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), pathParam, pathParamStr)))

				return
			}

			http.Error(
				w,
				fmt.Sprintf("Invalid path params for %s", pathParam),
				http.StatusBadRequest,
			)
		}
	}
}
