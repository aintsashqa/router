package router

import (
	"context"
)

const (
	_paramsCtxKey = "params_context_key"
)

type Params map[string]string

func ParamsFromContext(ctx context.Context) Params {
	value := ctx.Value(_paramsCtxKey)
	if value == nil {
		return Params{}
	}

	params, success := value.(Params)
	if !success {
		return Params{}
	}

	return params
}

func (p Params) append(key, value string) {
	p[key] = value
}

func (p Params) Has(key string) bool {
	_, found := p[key]
	return found
}

func (p Params) Get(key string) string {
	return p[key]
}
