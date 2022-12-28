package router

import (
	"context"
	"net/http"
	"strings"
)

type Router struct {
	routes           []route
	notFound         http.HandlerFunc
	methodNotAllowed http.HandlerFunc
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

func (router *Router) NotFound(handlerFn http.HandlerFunc) {
	router.notFound = handlerFn
}

func (router *Router) MethodNotAllowed(handlerFn http.HandlerFunc) {
	router.methodNotAllowed = handlerFn
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	segments := strings.Split(r.URL.Path, "/")
	routeNotFound := true

	for _, route := range router.routes {
		if !route.IsCurrentRoute(segments) {
			continue
		}

		routeNotFound = false
		if route.method != r.Method {
			continue
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

	if routeNotFound {
		if router.notFound != nil {
			router.notFound(w, r)
		} else {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}

		return
	}

	if router.methodNotAllowed != nil {
		router.methodNotAllowed(w, r)
	} else {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
