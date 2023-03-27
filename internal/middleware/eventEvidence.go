package middleware

import (
	"github.com/Bofry/host-fasthttp/internal"
	"github.com/Bofry/trace"
)

type EventEvidence struct {
	traceID   trace.TraceID
	spanID    trace.SpanID
	routePath *internal.RoutePath
}

func (e EventEvidence) RequestTraceID() trace.TraceID {
	return e.traceID
}

func (e EventEvidence) RequestSpanID() trace.SpanID {
	return e.spanID
}

func (e EventEvidence) RequestRoutePath() internal.RoutePath {
	return internal.RoutePath{
		Method: e.routePath.Method,
		Path:   e.routePath.Path,
	}
}
