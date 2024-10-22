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
	prefix      string
	routes      []*Route
	middlewares []Middleware
	subgroups   []*Group
	mux         *http.ServeMux
}

type Middleware struct {
	Name    string
	Handler MiddlewareHandler
}

type Route struct {
	pattern     string
	handler     http.Handler
	middlewares []Middleware
}

type MiddlewareHandler func(next http.Handler) http.Handler

func NewRouter() Router {
	mux := http.NewServeMux()
	return Router{groups: []*Group{}, mux: mux}
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

func NewGroup() Group {
	mux := http.NewServeMux()
	return Group{
		mux:         mux,
		middlewares: []Middleware{},
		routes:      []*Route{},
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
		g.Boot()

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

func (this *Group) Boot() {
	for _, subgroup := range this.subgroups {
		for _, g := range subgroup.subgroups {
			g.Boot()
		}

		cleanPrefix := strings.Trim(subgroup.GetPrefix(), "/")
		if len(cleanPrefix) == 0 {
			this.mux.Handle("/", subgroup.mux)
		} else {
			this.mux.Handle("/"+cleanPrefix+"/", http.StripPrefix("/"+cleanPrefix, subgroup.mux))
		}

		groupMiddleware := CombineMiddleware(this.middlewares)

		for _, r := range subgroup.routes {
			routeMiddleware := CombineMiddleware(r.middlewares)
			handler := groupMiddleware(routeMiddleware(r.handler))
			subgroup.mux.Handle(r.pattern, handler)
		}
	}

	groupMiddleware := CombineMiddleware(this.middlewares)

	for _, r := range this.routes {
		routeMiddleware := CombineMiddleware(r.middlewares)
		handler := groupMiddleware(routeMiddleware(r.handler))
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

func (this *Group) Middleware(middlewares []Middleware) *Group {
	this.middlewares = slices.Concat(this.middlewares, middlewares)

	return this
}

func (this *Group) Handle(pattern string, handler http.Handler) *Route {
	route := Route{pattern: pattern, handler: handler, middlewares: []Middleware{}}
	this.routes = append(this.routes, &route)

	return &route
}

func (this *Group) HandleFunc(pattern string, handler http.HandlerFunc) *Route {
	route := Route{pattern: pattern, handler: http.HandlerFunc(handler)}
	this.routes = append(this.routes, &route)

	return &route
}

func (this *Route) Middleware(middlewares []Middleware) *Route {
	this.middlewares = slices.Concat(this.middlewares, middlewares)

	return this
}
