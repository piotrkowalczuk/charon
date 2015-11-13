package charon

import "golang.org/x/net/context"

const (
	contextKeyRPCClient = "charon_rpc_client"
)

// NewContext returns a new Context that carries RPCClient instance.
func NewContext(ctx context.Context, c RPCClient) context.Context {
	return context.WithValue(ctx, contextKeyRPCClient, c)
}

// FromContext returns the RPCClient instance stored in context, if any.
func FromContext(ctx context.Context) (RPCClient, bool) {
	c, ok := ctx.Value(contextKeyRPCClient).(RPCClient)
	return c, ok
}

// Error ...
type Error struct {
	StatusCode   int
	InternalCode int
	Message      string
	Validation   map[string][]string
}

// Error ...
func (e *Error) Error() string {
	return e.Message
}

// AddValidation ...
func (e *Error) AddValidation(key, value string) {
	if e.Validation[key] == nil {
		e.Validation[key] = make([]string, 0, 1)
	}

	e.Validation[key] = append(e.Validation[key], value)
}
