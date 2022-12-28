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
			handlerFn        http.HandlerFunc
			notFound         http.HandlerFunc
			methodNotAllowed http.HandlerFunc
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
				handlerFn        http.HandlerFunc
				notFound         http.HandlerFunc
				methodNotAllowed http.HandlerFunc
			}{
				path:   "/",
				method: http.MethodGet,
				handlerFn: func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("success"))
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
				handlerFn        http.HandlerFunc
				notFound         http.HandlerFunc
				methodNotAllowed http.HandlerFunc
			}{
				path:   "/:parameter_key",
				method: http.MethodGet,
				handlerFn: func(w http.ResponseWriter, r *http.Request) {
					params := ParamsFromContext(r.Context())

					w.WriteHeader(http.StatusOK)
					w.Write([]byte("parameter_key:" + params.Get("parameter_key")))
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
				handlerFn        http.HandlerFunc
				notFound         http.HandlerFunc
				methodNotAllowed http.HandlerFunc
			}{
				path:      "/",
				method:    http.MethodGet,
				handlerFn: func(w http.ResponseWriter, r *http.Request) {},
			},
			response: struct {
				status int
				body   string
			}{
				status: http.StatusNotFound,
				body:   http.StatusText(http.StatusNotFound),
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
				handlerFn        http.HandlerFunc
				notFound         http.HandlerFunc
				methodNotAllowed http.HandlerFunc
			}{
				path:      "/",
				method:    http.MethodGet,
				handlerFn: func(w http.ResponseWriter, r *http.Request) {},
				notFound: func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte("not_found_custom_handler"))
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
				handlerFn        http.HandlerFunc
				notFound         http.HandlerFunc
				methodNotAllowed http.HandlerFunc
			}{
				path:      "/",
				method:    http.MethodPost,
				handlerFn: func(w http.ResponseWriter, r *http.Request) {},
			},
			response: struct {
				status int
				body   string
			}{
				status: http.StatusMethodNotAllowed,
				body:   http.StatusText(http.StatusMethodNotAllowed),
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
				handlerFn        http.HandlerFunc
				notFound         http.HandlerFunc
				methodNotAllowed http.HandlerFunc
			}{
				path:      "/",
				method:    http.MethodPost,
				handlerFn: func(w http.ResponseWriter, r *http.Request) {},
				methodNotAllowed: func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusMethodNotAllowed)
					w.Write([]byte("method_not_allowed_custom_handler"))
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

			router := Router{}
			router.NotFound(tC.route.notFound)
			router.MethodNotAllowed(tC.route.methodNotAllowed)
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
