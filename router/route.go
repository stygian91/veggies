package router

import (
	"context"
	"net/http"
	"slices"

	m "github.com/stygian91/veggies/router/middleware"
)

type Route struct {
	name    string
	pattern string
	handler http.Handler

	middlewares     []m.Middleware
	skipMiddlewares map[string]struct{}
}

func (this *Route) Middleware(middlewares ...m.Middleware) *Route {
	this.middlewares = slices.Concat(this.middlewares, middlewares)

	return this
}

func (this *Route) SkipMiddleware(names ...string) *Route {
	for _, name := range names {
		this.skipMiddlewares[name] = struct{}{}
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

func (this *Route) SetName(name string) *Route {
	this.name = name

	return this
}
