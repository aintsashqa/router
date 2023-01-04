package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRouterFuncServeHTTP(t *testing.T) {
	testCases := []struct {
		desc    string
		request struct {
			path   string
			method string
		}
		route struct {
			path             string
			method           string
			handlerFn        HandlerFunc
			errorHandlerFunc ErrorHandlerFunc
			middlewares      []Middleware
		}
		response struct {
			status int
			body   string
		}
	}{
		{
			desc: "_success_request_method",
			request: struct {
				path   string
				method string
			}{
				path:   "/",
				method: http.MethodGet,
			},
			route: struct {
				path             string
				method           string
				handlerFn        HandlerFunc
				errorHandlerFunc ErrorHandlerFunc
				middlewares      []Middleware
			}{
				path:   "/",
				method: http.MethodGet,
				handlerFn: func(ctx *Context) error {
					return ctx.Plain(http.StatusOK, "success")
				},
			},
			response: struct {
				status int
				body   string
			}{
				status: http.StatusOK,
				body:   "success",
			},
		},
		{
			desc: "_success_request_method_use_middleware",
			request: struct {
				path   string
				method string
			}{
				path:   "/",
				method: http.MethodGet,
			},
			route: struct {
				path             string
				method           string
				handlerFn        HandlerFunc
				errorHandlerFunc ErrorHandlerFunc
				middlewares      []Middleware
			}{
				path:   "/",
				method: http.MethodGet,
				handlerFn: func(ctx *Context) error {
					return ctx.Plain(http.StatusOK, "from_middleware_"+ctx.Context().Value("middleware").(string))
				},
				middlewares: []Middleware{
					func(next HandlerFunc) HandlerFunc {
						return func(ctx *Context) error {
							ctx.WithValue("middleware", "value")
							return next(ctx)
						}
					},
				},
			},
			response: struct {
				status int
				body   string
			}{
				status: http.StatusOK,
				body:   "from_middleware_value",
			},
		},
		{
			desc: "_success_request_method_with_params",
			request: struct {
				path   string
				method string
			}{
				path:   "/parameter_value",
				method: http.MethodGet,
			},
			route: struct {
				path             string
				method           string
				handlerFn        HandlerFunc
				errorHandlerFunc ErrorHandlerFunc
				middlewares      []Middleware
			}{
				path:   "/:parameter_key",
				method: http.MethodGet,
				handlerFn: func(ctx *Context) error {
					return ctx.Plain(http.StatusOK, "parameter_key:"+ctx.Param("parameter_key"))
				},
			},
			response: struct {
				status int
				body   string
			}{
				status: http.StatusOK,
				body:   "parameter_key:parameter_value",
			},
		},
		{
			desc: "_failure_request_route_not_found_default",
			request: struct {
				path   string
				method string
			}{
				path:   "/invalid_route",
				method: http.MethodGet,
			},
			route: struct {
				path             string
				method           string
				handlerFn        HandlerFunc
				errorHandlerFunc ErrorHandlerFunc
				middlewares      []Middleware
			}{
				path:      "/",
				method:    http.MethodGet,
				handlerFn: func(ctx *Context) error { return nil },
			},
			response: struct {
				status int
				body   string
			}{
				status: http.StatusNotFound,
				body:   ErrRouteNotFound.Error(),
			},
		},
		{
			desc: "_failure_request_route_not_found_custom",
			request: struct {
				path   string
				method string
			}{
				path:   "/invalid_route",
				method: http.MethodGet,
			},
			route: struct {
				path             string
				method           string
				handlerFn        HandlerFunc
				errorHandlerFunc ErrorHandlerFunc
				middlewares      []Middleware
			}{
				path:      "/",
				method:    http.MethodGet,
				handlerFn: func(ctx *Context) error { return nil },
				errorHandlerFunc: func(ctx *Context, err error) {
					ctx.Plain(http.StatusNotFound, "not_found_custom_handler")
				},
			},
			response: struct {
				status int
				body   string
			}{
				status: http.StatusNotFound,
				body:   "not_found_custom_handler",
			},
		},
		{
			desc: "_failure_request_method_not_allow_default",
			request: struct {
				path   string
				method string
			}{
				path:   "/",
				method: http.MethodGet,
			},
			route: struct {
				path             string
				method           string
				handlerFn        HandlerFunc
				errorHandlerFunc ErrorHandlerFunc
				middlewares      []Middleware
			}{
				path:      "/",
				method:    http.MethodPost,
				handlerFn: func(ctx *Context) error { return nil },
			},
			response: struct {
				status int
				body   string
			}{
				status: http.StatusMethodNotAllowed,
				body:   ErrMethodNotAllowed.Error(),
			},
		},
		{
			desc: "_failure_request_method_not_allow_custom",
			request: struct {
				path   string
				method string
			}{
				path:   "/",
				method: http.MethodGet,
			},
			route: struct {
				path             string
				method           string
				handlerFn        HandlerFunc
				errorHandlerFunc ErrorHandlerFunc
				middlewares      []Middleware
			}{
				path:      "/",
				method:    http.MethodPost,
				handlerFn: func(ctx *Context) error { return nil },
				errorHandlerFunc: func(ctx *Context, err error) {
					ctx.Plain(http.StatusMethodNotAllowed, "method_not_allowed_custom_handler")
				},
			},
			response: struct {
				status int
				body   string
			}{
				status: http.StatusMethodNotAllowed,
				body:   "method_not_allowed_custom_handler",
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			request := httptest.NewRequest(tC.request.method, tC.request.path, nil)
			recorder := httptest.NewRecorder()

			router := New()
			if tC.route.errorHandlerFunc != nil {
				router.ErrorHandlerFunc = tC.route.errorHandlerFunc
			}

			router.Use(tC.route.middlewares...)
			if tC.route.handlerFn != nil {
				router.Route(tC.route.path, tC.route.method, tC.route.handlerFn)
			}

			router.ServeHTTP(recorder, request)

			if tC.response.status != recorder.Code {
				t.Errorf("expects status '%d' but got '%d'", tC.response.status, recorder.Code)
			}

			actualBody := strings.TrimSpace(string(recorder.Body.Bytes()))
			if tC.response.body != actualBody {
				t.Errorf("expects body '%s' but got '%s'", tC.response.body, actualBody)
			}
		})
	}
}
