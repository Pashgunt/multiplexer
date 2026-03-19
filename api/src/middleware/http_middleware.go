package middleware

import (
	"fmt"
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Chain(h http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}

	return h
}

func AllowHttpMethodMiddleware(method string) Middleware {
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
