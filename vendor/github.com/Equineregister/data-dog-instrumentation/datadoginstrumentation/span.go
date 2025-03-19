package datadoginstrumentation

import (
	"context"
	"errors"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var (
	ErrNoSpanFound = errors.New("no span found in context")
)

// StartSpan starts a new span with the given operation name and resource name, within the supplied context.
// Typically you would use this function to start a new root span.
// It will overwrite any existing span in the context.
// The new span is contained in the returned context - use this context with FinishSpan to finish the span.
func StartSpan(ctx context.Context, operationName string, resourceName string) context.Context {
	span := tracer.StartSpan(operationName,
		tracer.ResourceName(resourceName),
	)
	ctx = tracer.ContextWithSpan(ctx, span)

	return ctx
}

// StartNestedSpan starts a new span with the given operation name and resource name, as a nested span found within the supplied context.
// If no parent span is found in the context, a new root span is started.
// The new nested span is contained in the returned context - use this context with FinishSpan to finish the nested span.
func StartNestedSpan(ctx context.Context, operationName string, resourceName string) context.Context {
	_, ctx = tracer.StartSpanFromContext(ctx, operationName,
		tracer.ResourceName(resourceName),
	)
	return ctx
}

// FinishSpan finishes the span found within the supplied context.
// If no span is found, an error is returned.
// If an error is supplied as parameter opErr, it is set on the finished span.
func FinishSpan(ctx context.Context, opErr error) error {
	span, ok := tracer.SpanFromContext(ctx)
	if !ok {
		return ErrNoSpanFound
	}

	span.Finish(tracer.WithError(opErr))
	return nil
}
