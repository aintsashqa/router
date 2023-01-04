package router

import (
	"net/http"
	"testing"
)

func BenchmarkGetWithoutParams(b *testing.B) {
	router := Router{}
	router.Get("/foo", func(ctx *Context) error { return nil })

	for i := 0; i < b.N; i++ {
		request, _ := http.NewRequest(http.MethodGet, "/foo", nil)

		router.ServeHTTP(nil, request)
	}
}

func BenchmarkGetWithParams(b *testing.B) {
	router := Router{}
	router.Get("/foo/:baz", func(ctx *Context) error { return nil })

	for i := 0; i < b.N; i++ {
		request, _ := http.NewRequest(http.MethodGet, "/foo/baz", nil)

		router.ServeHTTP(nil, request)
	}
}
