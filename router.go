package router

import (
	"context"
	"net/http"
)

type Router struct {
	routes []route
}

func New() Router {
	return Router{}
}

func (router *Router) Route(path, method string, handlerFn http.HandlerFunc) {
	router.routes = append(router.routes, route{
		path:      path,
		method:    method,
		handlerFn: handlerFn,
	})
}

func (router *Router) Get(path string, handlerFn http.HandlerFunc) {
	router.Route(path, http.MethodGet, handlerFn)
}

func (router *Router) Post(path string, handlerFn http.HandlerFunc) {
	router.Route(path, http.MethodPost, handlerFn)
}

func (router *Router) Put(path string, handlerFn http.HandlerFunc) {
	router.Route(path, http.MethodPut, handlerFn)
}

func (router *Router) Patch(path string, handlerFn http.HandlerFunc) {
	router.Route(path, http.MethodPatch, handlerFn)
}

func (router *Router) Delete(path string, handlerFn http.HandlerFunc) {
	router.Route(path, http.MethodDelete, handlerFn)
}

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range router.routes {
		if route.method != r.Method {
			continue
		}

		if route.IsPathEquals(r.URL.Path) {
			params := route.GetPathParameters(r.URL.Path)
			for key, value := range params {
				withParamCtx := context.WithValue(r.Context(), key, value)
				r = r.WithContext(withParamCtx)
			}

			route.handlerFn(w, r)
			return
		}
	}

	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}
