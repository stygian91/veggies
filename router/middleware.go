package router

import (
	"log"
	"net/http"
	"slices"
)

type Middleware struct {
	Name    string
	Handler MiddlewareHandler
}

type MiddlewareHandler func(next http.Handler) http.Handler

type StatusResponseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (this *StatusResponseWriter) Write(b []byte) (int, error) {
	if !this.wroteHeader {
		this.WriteHeader(http.StatusOK)
	}

	return this.ResponseWriter.Write(b)
}

func (this *StatusResponseWriter) WriteHeader(code int) {
	this.status = code
	this.wroteHeader = true
	this.ResponseWriter.WriteHeader(code)
}

// TODO:
// make message configurable through a callback
// make logger that's used configurable, instead of using default logger
// move middlewares to a subfolder
var LogMiddleware Middleware = Middleware{
	Name: "log",
	Handler: func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			writer := StatusResponseWriter{ResponseWriter: w}
			next.ServeHTTP(&writer, r)
			log.Printf("uri: %s, status: %d, remote addr: %s, user agent: %s", r.RequestURI, writer.status, r.RemoteAddr, r.UserAgent())
		})
	},
}

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
