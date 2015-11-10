package routing

import (
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
)

const (
	contextKeyParams = "routing_params"
)

// NewParamsContext ...
func NewParamsContext(ctx context.Context, params httprouter.Params) context.Context {
	return context.WithValue(ctx, contextKeyParams, params)
}

// ParamsFromContext ...
func ParamsFromContext(ctx context.Context) (httprouter.Params, bool) {
	params, ok := ctx.Value(contextKeyParams).(httprouter.Params)

	return params, ok
}

// ParamFromContext ...
func ParamFromContext(ctx context.Context, paramName string) (string, bool) {
	var ok bool
	var params httprouter.Params

	if params, ok = ParamsFromContext(ctx); !ok {
		return "", ok
	}

	param := params.ByName(paramName)

	return param, param != ""
}

// ParamFromContextOr ...
func ParamFromContextOr(ctx context.Context, paramName, or string) string {
	if param, ok := ParamFromContext(ctx, paramName); ok {
		return param
	}

	return or
}
