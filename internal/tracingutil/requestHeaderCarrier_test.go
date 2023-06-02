package tracingutil

import (
	"context"
	"reflect"
	"testing"

	"github.com/valyala/fasthttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/exp/slices"
)

var (
	traceIDStr = "4bf92f3577b34da6a3ce929d0e0e4736"
	spanIDStr  = "00f067aa0ba902b7"

	__TEST_TRACE_ID = mustTraceIDFromHex(traceIDStr)
	__TEST_SPAN_ID  = mustSpanIDFromHex(spanIDStr)

	__TEST_PROPAGATOR = propagation.TraceContext{}
	__TEST_CONTEXT    = mustSpanContext()
)

func mustTraceIDFromHex(s string) (t trace.TraceID) {
	var err error
	t, err = trace.TraceIDFromHex(s)
	if err != nil {
		panic(err)
	}
	return
}

func mustSpanIDFromHex(s string) (t trace.SpanID) {
	var err error
	t, err = trace.SpanIDFromHex(s)
	if err != nil {
		panic(err)
	}
	return
}

func mustSpanContext() context.Context {
	sc := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    __TEST_TRACE_ID,
		SpanID:     __TEST_SPAN_ID,
		TraceFlags: 0,
	})
	return trace.ContextWithSpanContext(context.Background(), sc)
}

func TestRequestHeaderCarrier(t *testing.T) {
	var (
		header     = new(fasthttp.RequestHeader)
		propagator = propagation.TraceContext{}
	)

	carrier := NewRequestHeaderCarrier(header)

	// inject
	{
		propagator.Inject(__TEST_CONTEXT, carrier)

		header.DisableNormalizing()
		traceparent := header.Peek("traceparent")
		if len(traceparent) == 0 {
			t.Error("missing request header 'traceparent'")
		}
		header.EnableNormalizing()
	}

	// fields
	{
		keys := carrier.Keys()
		if !slices.Contains(keys, "traceparent") {
			t.Error("missing request header 'traceparent'")
		}
	}

	// extract
	{
		ctx := propagator.Extract(context.Background(), carrier)
		sc := trace.SpanContextFromContext(ctx)
		var expectedTraceID = __TEST_TRACE_ID
		if !reflect.DeepEqual(expectedTraceID, sc.TraceID()) {
			t.Errorf("TRACE ID expect: %v, got: %v", expectedTraceID, sc.TraceID())
		}
		var expectedSpanID = __TEST_SPAN_ID
		if !reflect.DeepEqual(expectedSpanID, sc.SpanID()) {
			t.Errorf("SPAN ID expect: %v, got: %v", expectedSpanID, sc.SpanID())
		}
	}
}
