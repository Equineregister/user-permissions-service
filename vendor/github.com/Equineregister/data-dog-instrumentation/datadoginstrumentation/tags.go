package datadoginstrumentation

import (
	"context"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func AddSpanTag(ctx context.Context, key string, value any) bool {
	if span, ok := tracer.SpanFromContext(ctx); ok {
		span.SetTag(key, value)
		return true
	}
	return false
}
