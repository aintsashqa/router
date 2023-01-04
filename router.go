package router

import (
	"errors"
	"net/http"
	"strings"
)

var (
	ErrRouteNotFound    error = errors.New("Route " + http.StatusText(http.StatusNotFound))
	ErrMethodNotAllowed error = errors.New(http.StatusText(http.StatusMethodNotAllowed))
)

type HandlerFunc func(*Context) error

type ErrorHandlerFunc func(*Context, error)

func DefaultErrorHandlerFunc(ctx *Context, err error) {
	respond := func(ctx *Context, statusCode int, err error) {
		contentType := ctx.request.Header.Get(headerContentType)
		if strings.Contains(contentType, contentTypeJson) {
			ctx.Json(statusCode, map[string]string{"error": err.Error()})
		} else {
			ctx.Plain(statusCode, err.Error())
		}
	}

	switch err {
	case ErrUnsupportedMediaType:
		ctx.Plain(http.StatusUnsupportedMediaType, http.StatusText(http.StatusUnsupportedMediaType))
		break

	case ErrRouteNotFound:
		respond(ctx, http.StatusNotFound, err)
		break

	case ErrMethodNotAllowed:
		respond(ctx, http.StatusMethodNotAllowed, err)
		break

	default:
		respond(ctx, http.StatusInternalServerError, err)
	}
}

type Middleware func(HandlerFunc) HandlerFunc

type Router struct {
	routes      []route
	middlewares []Middleware

	ErrorHandlerFunc ErrorHandlerFunc
}

func New() Router {
	return Router{
		routes:           make([]route, 0),
		middlewares:      make([]Middleware, 0),
		ErrorHandlerFunc: DefaultErrorHandlerFunc,
	}
}

func (router *Router) Route(path, method string, handlerFn HandlerFunc) {
	router.routes = append(router.routes, newRoute(path, method, handlerFn))
}

func (router *Router) Use(middlewares ...Middleware) {
	for _, middleware := range middlewares {
		router.middlewares = append(router.middlewares, middleware)
	}
}

func (router *Router) Get(path string, handlerFn HandlerFunc) {
	router.Route(path, http.MethodGet, handlerFn)
}

func (router *Router) Post(path string, handlerFn HandlerFunc) {
	router.Route(path, http.MethodPost, handlerFn)
}

func (router *Router) Put(path string, handlerFn HandlerFunc) {
	router.Route(path, http.MethodPut, handlerFn)
}

func (router *Router) Patch(path string, handlerFn HandlerFunc) {
	router.Route(path, http.MethodPatch, handlerFn)
}

func (router *Router) Delete(path string, handlerFn HandlerFunc) {
	router.Route(path, http.MethodDelete, handlerFn)
}

func (router *Router) execMiddlewares(handlerFn HandlerFunc) HandlerFunc {
	for _, middleware := range router.middlewares {
		handlerFn = middleware(handlerFn)
	}

	return handlerFn
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	segments := strings.Split(r.URL.Path, "/")
	routeNotFound := true
	ctx := Context{
		ctx:     r.Context(),
		writer:  w,
		request: r,
		params:  make(Params),
	}

	for _, route := range router.routes {
		if !route.IsCurrentRoute(segments) {
			continue
		}

		routeNotFound = false
		if route.method != r.Method {
			continue
		}

		handlerFn := router.execMiddlewares(route.handlerFn)
		route.GetPathParameters(ctx.params, segments)
		if err := handlerFn(&ctx); err != nil {
			router.ErrorHandlerFunc(&ctx, err)
		}

		return
	}

	if routeNotFound {
		router.ErrorHandlerFunc(&ctx, ErrRouteNotFound)
		return
	}

	router.ErrorHandlerFunc(&ctx, ErrMethodNotAllowed)
}
