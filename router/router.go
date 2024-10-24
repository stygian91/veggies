package router

import (
	"net/http"
	"slices"
	"strings"
)

type Router struct {
	groups []*Group
	mux    *http.ServeMux
}

type Group struct {
	prefix    string
	routes    []*Route
	subgroups []*Group
	mux       *http.ServeMux

	middlewares     []Middleware
	skipMiddlewares map[string]empty
}

type Middleware struct {
	Name    string
	Handler MiddlewareHandler
}

type empty struct{}

type Route struct {
	pattern string
	handler http.Handler

	middlewares     []Middleware
	skipMiddlewares map[string]empty
}

type MiddlewareHandler func(next http.Handler) http.Handler

var (
	router *Router
)

func init() {
	r := NewRouter()
	router = &r
}

func Get() *Router {
	return router
}

func NewRouter() Router {
	mux := http.NewServeMux()
	return Router{groups: []*Group{}, mux: mux}
}

func NewGroup() Group {
	mux := http.NewServeMux()
	return Group{
		mux:             mux,
		middlewares:     []Middleware{},
		routes:          []*Route{},
		skipMiddlewares: map[string]empty{},
	}
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

func (this *Router) Group(cb func(*Group)) *Group {
	subgroup := NewGroup()

	cb(&subgroup)
	this.groups = append(this.groups, &subgroup)

	return &subgroup
}

func (this *Router) Boot() {
	for _, g := range this.groups {
		g.boot(filterMiddleware(g.middlewares, g.skipMiddlewares))

		cleanPrefix := strings.Trim(g.GetPrefix(), "/")
		if len(cleanPrefix) == 0 {
			this.mux.Handle("/", g.mux)
		} else {
			this.mux.Handle("/"+cleanPrefix+"/", http.StripPrefix("/"+cleanPrefix, g.mux))
		}
	}
}

func (this Router) Mux() *http.ServeMux {
	return this.mux
}

func (this *Group) boot(middlewares []Middleware) {
	for _, subgroup := range this.subgroups {
		subgroup.boot(filterMiddleware(slices.Concat(middlewares, subgroup.middlewares), subgroup.skipMiddlewares))

		cleanPrefix := strings.Trim(subgroup.GetPrefix(), "/")
		if len(cleanPrefix) == 0 {
			this.mux.Handle("/", subgroup.mux)
		} else {
			this.mux.Handle("/"+cleanPrefix+"/", http.StripPrefix("/"+cleanPrefix, subgroup.mux))
		}
	}

	for _, r := range this.routes {
		routeMiddleware := CombineMiddleware(filterMiddleware(slices.Concat(middlewares, r.middlewares), r.skipMiddlewares))
		handler := routeMiddleware(r.handler)
		this.mux.Handle(r.pattern, handler)
	}
}

func (this *Group) Group(cb func(*Group)) *Group {
	subgroup := NewGroup()
	cb(&subgroup)
	this.subgroups = append(this.subgroups, &subgroup)

	return &subgroup
}

func (this Group) GetPrefix() string {
	return this.prefix
}

func (this *Group) SetPrefix(prefix string) *Group {
	this.prefix = prefix

	return this
}

func (this *Group) Middleware(middlewares ...Middleware) *Group {
	this.middlewares = slices.Concat(this.middlewares, middlewares)

	return this
}

func (this *Group) SkipMiddleware(names ...string) *Group {
	for _, name := range names {
		this.skipMiddlewares[name] = empty{}
	}

	return this
}

func (this *Group) Handle(pattern string, handler http.Handler) *Route {
	route := Route{pattern: pattern, handler: handler, middlewares: []Middleware{}, skipMiddlewares: map[string]empty{}}
	this.routes = append(this.routes, &route)

	return &route
}

func (this *Group) HandleFunc(pattern string, handler http.HandlerFunc) *Route {
	route := Route{pattern: pattern, handler: http.HandlerFunc(handler), middlewares: []Middleware{}, skipMiddlewares: map[string]empty{}}
	this.routes = append(this.routes, &route)

	return &route
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

func filterMiddleware(middlewares []Middleware, skips map[string]empty) []Middleware {
	return slices.DeleteFunc(middlewares, func(m Middleware) bool {
		_, isInSkips := skips[m.Name]
		return isInSkips
	})
}
