package internal

import "github.com/Bofry/trace"

type RequestState struct {
	Tracer    *trace.SeverityTracer
	Span      *trace.SeveritySpan
	RoutePath *RoutePath
}
