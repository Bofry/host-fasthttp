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
	if v, ok := ctx.(*http.RequestCtx); ok {
		return extractSpan(v)
	}
	return nil
}

func extractSpan(ctx *http.RequestCtx) *trace.SeveritySpan {
	obj := requestutil.ExtractSpan(ctx)
	v, ok := obj.(*trace.SeveritySpan)
	if ok {
		return v
	}
	return nil
}
