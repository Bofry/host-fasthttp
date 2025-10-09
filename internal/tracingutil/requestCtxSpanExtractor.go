package tracingutil

import (
	"context"

	"github.com/Bofry/host-fasthttp/internal/requestutil"
	"github.com/Bofry/trace"
	http "github.com/valyala/fasthttp"
)

var (
	_ trace.SpanExtractor = RequestCtxSpanExtractor(0)
)

type RequestCtxSpanExtractor int

// Extract implements trace.SpanExtractor
func (RequestCtxSpanExtractor) Extract(ctx context.Context) *trace.SeveritySpan {
	if rx, ok := ctx.(*http.RequestCtx); ok {
		reply := requestutil.ExtractSpan(rx)
		span, ok := reply.(*trace.SeveritySpan)
		if ok {
			return span
		}
	}
	// FIX: Return nil instead of calling trace.SpanFromContext to avoid infinite recursion
	// When this extractor is registered globally, calling SpanFromContext would trigger
	// this Extract method again, causing a stack overflow
	return nil
}
