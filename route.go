package router

import (
	"net/http"
	"strings"
)

type route struct {
	path      string
	method    string
	handlerFn http.HandlerFunc
}

func (r route) IsPathEquals(path string) bool {
	currentPathParts := strings.Split(r.path, "/")
	inputPathParts := strings.Split(path, "/")

	if len(currentPathParts) != len(inputPathParts) {
		return false
	}

	for i := 0; i < len(currentPathParts); i++ {
		if strings.HasPrefix(currentPathParts[i], ":") {
			continue
		}

		if currentPathParts[i] != inputPathParts[i] {
			return false
		}
	}

	return true
}

func (r route) GetPathParameters(path string) (result map[string]any) {
	currentPathParts := strings.Split(r.path, "/")
	inputPathParts := strings.Split(path, "/")

	for i := 0; i < len(currentPathParts); i++ {
		if !strings.HasPrefix(currentPathParts[i], ":") {
			continue
		}

		if result == nil {
			result = make(map[string]any)
		}

		key := strings.TrimPrefix(currentPathParts[i], ":")
		result[key] = inputPathParts[i]
	}

	return
}
