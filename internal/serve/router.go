package serve

import (
	"fmt"
	"net/http"
)

type Router struct {
	index       string
	mux         *http.ServeMux
	middlewares []Middleware
}

func NewRouter(mux *http.ServeMux) *Router {
	return &Router{
		mux:         mux,
		middlewares: []Middleware{},
	}
}

func (r *Router) Route(pattern string) *Router {
	middlewares := []Middleware{}
	copy(middlewares, r.middlewares)

	return &Router{
		mux:         r.mux,
		index:       r.index + pattern,
		middlewares: middlewares,
	}
}

func (r *Router) Use(middlewares ...Middleware) {
	r.middlewares = append(r.middlewares, middlewares...)
}

func (r *Router) Method(method, pattern string, h http.Handler) {
	for _, mw := range r.middlewares {
		h = mw.Next(h)
	}
	fmt.Println(method + " " + r.index + pattern)
	r.mux.Handle(method+" "+r.index+pattern, h)
}

func (r *Router) Get(pattern string, h http.Handler) {
	r.Method(http.MethodGet, pattern, h)
}

func (r *Router) Post(pattern string, h http.Handler) {
	r.Method(http.MethodPost, pattern, h)
}

func (r *Router) Put(pattern string, h http.Handler) {
	r.Method(http.MethodPut, pattern, h)
}

func (r *Router) Patch(pattern string, h http.Handler) {
	r.Method(http.MethodPatch, pattern, h)
}

func (r *Router) Delete(pattern string, h http.Handler) {
	r.Method(http.MethodDelete, pattern, h)
}

func (r *Router) Options(pattern string, h http.Handler) {
	r.Method(http.MethodOptions, pattern, h)
}
