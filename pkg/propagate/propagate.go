package propagate

import (
	"context"
)

type Values map[string]string

type Getter interface {
	Get(key string) (value string)
}

type Setter interface {
	Set(key, value string)
}

var headers = []string{
	"x-request-id",
	"x-b3-traceid",
	"x-b3-spanid",
	"x-b3-parentspanid",
	"x-b3-sampled",
	"x-b3-flags",
	"x-ot-span-context",
	"uber-trace-id",
}

type valuesKey struct{}

// NewContext creates a new context with propagated header values attached.
func NewContext(ctx context.Context, values Values) context.Context {
	return context.WithValue(ctx, valuesKey{}, values)
}

// FromContext returns the propagated header values in ctx if it exists.
func FromContext(ctx context.Context) (values Values, ok bool) {
	values, ok = ctx.Value(valuesKey{}).(Values)
	return
}

// FromRequest extracts propagation headers from an incoming HTTP request.
func FromRequest(getter Getter) Values {
	values := make(Values, len(headers))
	for i := range headers {
		value := getter.Get(headers[i])
		if value != "" {
			values[headers[i]] = value
		}
	}
	return values
}

// Propagate sets the propagation headers to an outgoing HTTP request.
func (v Values) Propagate(setter Setter) {
	for k, v := range v {
		setter.Set(k, v)
	}
}

// Incoming extracts propagation headers from an incoming HTTP request and stores them in a new context.
func Incoming(ctx context.Context, getter Getter) context.Context {
	values := FromRequest(getter)
	return NewContext(ctx, values)
}

// Outgoing propagates header values to an outgoing HTTP request.
func Outgoing(ctx context.Context, setter Setter) {
	values, ok := FromContext(ctx)
	if ok {
		values.Propagate(setter)
	}
}
