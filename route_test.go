package router

import (
	"strings"
	"testing"
)

func TestRouteFuncIsCurrentRoute(t *testing.T) {
	testCases := []struct {
		desc   string
		path   string
		input  string
		result bool
	}{
		{
			desc:   "_success_blank_path",
			path:   "/",
			input:  "/",
			result: true,
		},
		{
			desc:   "_success_not_blank_path",
			path:   "/route",
			input:  "/route",
			result: true,
		},
		{
			desc:   "_success_path_with_parameter",
			path:   "/route/:parameter_key",
			input:  "/route/parameter-value",
			result: true,
		},
		{
			desc:   "_failure",
			path:   "/route/success",
			input:  "/route/failure",
			result: false,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			route := newRoute(tC.path, "", nil)

			result := route.IsCurrentRoute(strings.Split(tC.input, "/"))

			if tC.result != result {
				t.Errorf("expects result '%+v' but got '%+v'", tC.result, result)
			}
		})
	}
}

func TestRouteFuncGetPathParameters(t *testing.T) {
	testCases := []struct {
		desc   string
		path   string
		input  string
		result map[string]any
	}{
		{
			desc:   "_success_with_out_parameters",
			path:   "/route",
			input:  "/route",
			result: nil,
		},
		{
			desc:  "_success_with_single_parameter",
			path:  "/route/:route_id",
			input: "/route/1",
			result: map[string]any{
				"route_id": "1",
			},
		},
		{
			desc:  "_success_with_many_parameters",
			path:  "/route/:route_id/:subroute_id",
			input: "/route/1/subroute",
			result: map[string]any{
				"route_id":    "1",
				"subroute_id": "subroute",
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			route := newRoute(tC.path, "", nil)
			params := make(Params)

			route.GetPathParameters(params, strings.Split(tC.input, "/"))

			for key, expectValue := range tC.result {
				actualValue, found := params[key]
				if !found {
					t.Errorf("expected key '%s' not found", key)
				}

				if expectValue != actualValue {
					t.Errorf("expects value '%+v' but got '%+v'", expectValue, actualValue)
				}
			}
		})
	}
}
