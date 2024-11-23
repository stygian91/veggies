package middleware

import (
	"fmt"
	"log"
	"net/http"
)

var logCallback LogFunction = func(r *http.Request, status int) string {
	return fmt.Sprintf("method: %s, uri: %s, status: %d, addr: %s", r.Method, r.RequestURI, status, r.RemoteAddr)
}

var logger *log.Logger = log.Default()

type LogFunction func(r *http.Request, status int) string

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

var LogMiddleware Middleware = Middleware{
	Name: "log",
	Handler: func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			writer := StatusResponseWriter{ResponseWriter: w}
			next.ServeHTTP(&writer, r)
			logger.Println(logCallback(r, writer.status))
		})
	},
}

func SetLogCallback(cb LogFunction) {
	logCallback = cb
}

func SetLogger(l *log.Logger) {
	logger = l
}
