package router

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

const (
	headerContentType string = "Content-Type"

	contentTypePlain string = "text/plain"
	contentTypeJson  string = "application/json"
)

var (
	ErrUnsupportedMediaType error = errors.New(http.StatusText(http.StatusUnsupportedMediaType))
)

type Context struct {
	ctx     context.Context
	writer  http.ResponseWriter
	request *http.Request
	params  Params
}

func (c *Context) Request() *http.Request {
	return c.request
}

func (c *Context) Context() context.Context {
	return c.ctx
}

func (c *Context) WithValue(key any, value any) {
	c.ctx = context.WithValue(c.ctx, key, value)
}

func (c *Context) Param(key string) string {
	if c.params.Has(key) {
		return c.params.Get(key)
	}

	return ""
}

func (c *Context) Query(key string) string {
	if q := c.request.URL.Query(); q.Has(key) {
		return q.Get(key)
	}

	return ""
}

func (c *Context) ParseJson(v any) error {
	contentType := c.request.Header.Get(headerContentType)
	if !strings.Contains(contentType, contentTypeJson) {
		return ErrUnsupportedMediaType
	}

	return json.NewDecoder(c.request.Body).Decode(v)
}

func (c *Context) Plain(statusCode int, v string) error {
	c.writer.Header().Set(headerContentType, contentTypePlain)
	c.writer.WriteHeader(statusCode)
	c.writer.Write([]byte(v))
	return nil
}

func (c *Context) Json(statusCode int, v any) error {
	c.writer.Header().Set(headerContentType, contentTypeJson)
	c.writer.WriteHeader(statusCode)
	return json.NewEncoder(c.writer).Encode(v)
}
