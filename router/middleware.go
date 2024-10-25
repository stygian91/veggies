package router

import (
	"net/http"
	"slices"
)

type Middleware struct {
	Name    string
	Handler MiddlewareHandler
}

type MiddlewareHandler func(next http.Handler) http.Handler

func CombineMiddleware(middlewares []Middleware) MiddlewareHandler {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			m := middlewares[i]
			next = m.Handler(next)
		}

		return next
	}
}

func filterMiddleware(middlewares []Middleware, skips map[string]empty) []Middleware {
	return slices.DeleteFunc(middlewares, func(m Middleware) bool {
		_, isInSkips := skips[m.Name]
		return isInSkips
	})
}
