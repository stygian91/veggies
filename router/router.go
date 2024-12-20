package router

import (
	"net/http"
	"slices"
	"strings"

	m "github.com/stygian91/veggies/router/middleware"
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

	middlewares     []m.Middleware
	skipMiddlewares map[string]struct{}
}

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
		middlewares:     []m.Middleware{},
		routes:          []*Route{},
		skipMiddlewares: map[string]struct{}{},
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
		g.boot(m.FilterMiddleware(g.middlewares, g.skipMiddlewares))

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

func (this *Group) boot(middlewares []m.Middleware) {
	for _, subgroup := range this.subgroups {
		subgroup.boot(
			m.FilterMiddleware(
				slices.Concat(middlewares, subgroup.middlewares),
				subgroup.skipMiddlewares,
			),
		)

		cleanPrefix := strings.Trim(subgroup.GetPrefix(), "/")
		if len(cleanPrefix) == 0 {
			this.mux.Handle("/", subgroup.mux)
		} else {
			this.mux.Handle("/"+cleanPrefix+"/", http.StripPrefix("/"+cleanPrefix, subgroup.mux))
		}
	}

	for _, r := range this.routes {
		routeMiddleware := m.CombineMiddleware(
			m.FilterMiddleware(
				slices.Concat(middlewares, r.middlewares),
				r.skipMiddlewares,
			),
		)
		r.handler = routeMiddleware(r.handler)
		this.mux.Handle(r.pattern, r)
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

func (this *Group) Middleware(middlewares ...m.Middleware) *Group {
	this.middlewares = slices.Concat(this.middlewares, middlewares)

	return this
}

func (this *Group) SkipMiddleware(names ...string) *Group {
	for _, name := range names {
		this.skipMiddlewares[name] = struct{}{}
	}

	return this
}

func (this *Group) Handle(pattern string, handler http.Handler) *Route {
	route := Route{
		pattern:         pattern,
		handler:         handler,
		middlewares:     []m.Middleware{},
		skipMiddlewares: map[string]struct{}{},
	}
	this.routes = append(this.routes, &route)

	return &route
}

func (this *Group) HandleFunc(pattern string, handler http.HandlerFunc) *Route {
	route := Route{
		pattern:         pattern,
		handler:         http.HandlerFunc(handler),
		middlewares:     []m.Middleware{},
		skipMiddlewares: map[string]struct{}{},
	}
	this.routes = append(this.routes, &route)

	return &route
}
