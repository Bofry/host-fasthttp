package tracingutil

import (
	"context"

	"github.com/Bofry/trace"
	http "github.com/valyala/fasthttp"
)

var (
	_ trace.SpanExtractor = RequestCtxSpanExtractor(0)
)

type RequestCtxSpanExtractor int

// Extract implements trace.SpanExtractor
func (RequestCtxSpanExtractor) Extract(ctx context.Context) *trace.SeveritySpan {
	if v, ok := ctx.(*http.RequestCtx); ok {
		return ExtractSpan(v)
	}
	return nil
}
