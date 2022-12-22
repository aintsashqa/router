package router

import (
	"net/http"
	"strings"
)

type route struct {
	method        string
	handlerFn     http.HandlerFunc
	segments      []string
	segmentsCount int
	segmentsKeys  map[int]string
}

func newRoute(path, method string, handlerFn http.HandlerFunc) route {
	segments := strings.Split(path, "/")

	segmentsKeys := make(map[int]string)
	for i, s := range segments {
		if strings.HasPrefix(s, ":") {
			segmentsKeys[i] = strings.TrimPrefix(s, ":")
		}
	}

	return route{
		method:        method,
		handlerFn:     handlerFn,
		segments:      segments,
		segmentsCount: len(segments),
		segmentsKeys:  segmentsKeys,
	}
}

func (r *route) IsCurrentRoute(segments []string) bool {
	if r.segmentsCount != len(segments) {
		return false
	}

	for i, s := range r.segments {
		if _, found := r.segmentsKeys[i]; !found && s != segments[i] {
			return false
		}
	}

	return true
}

func (r *route) GetPathParameters(segments []string) Params {
	params := make(Params)
	for i, k := range r.segmentsKeys {
		params.append(k, segments[i])
	}

	return params
}
