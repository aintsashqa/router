package router

import (
	"context"
	"net/http"
	"strings"
)

type Router struct {
	routes []route
}

func New() Router {
	return Router{}
}

func (router *Router) Route(path, method string, handlerFn http.HandlerFunc) {
	router.routes = append(router.routes, newRoute(path, method, handlerFn))
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

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	segments := strings.Split(r.URL.Path, "/")

	for _, route := range router.routes {
		if !route.IsCurrentRoute(segments) {
			continue
		}

		if route.method != r.Method {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		if len(route.segmentsKeys) == 0 {
			route.handlerFn(w, r)
			return
		}

		params := route.GetPathParameters(segments)
		r = r.WithContext(context.WithValue(r.Context(), _paramsCtxKey, params))

		route.handlerFn(w, r)
		return
	}
}
