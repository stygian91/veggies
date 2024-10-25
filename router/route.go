package router

import (
	"context"
	"net/http"
	"slices"
)

type Route struct {
	name    string
	pattern string
	handler http.Handler

	middlewares     []Middleware
	skipMiddlewares map[string]empty
}

func (this *Route) Middleware(middlewares ...Middleware) *Route {
	this.middlewares = slices.Concat(this.middlewares, middlewares)

	return this
}

func (this *Route) SkipMiddleware(names ...string) *Route {
	for _, name := range names {
		this.skipMiddlewares[name] = empty{}
	}

	return this
}

func (this Route) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	newCtx := context.WithValue(ctx, "routeName", this.name)
	*r = *r.WithContext(newCtx)

	this.handler.ServeHTTP(w, r)
}

func (this Route) GetName() string {
	return this.name
}

func (this *Route) SetName(name string) {
	this.name = name
}
