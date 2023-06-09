package tracing

import (
	"github.com/Bofry/host-fasthttp/internal/tracingutil"
	"github.com/Bofry/trace"
	http "github.com/valyala/fasthttp"
)

var (
	defaultRequestCtxSpanExtractor = tracingutil.RequestCtxSpanExtractor(0)
)

func SpanFromRequestCtx(ctx *http.RequestCtx) *trace.SeveritySpan {
	return defaultRequestCtxSpanExtractor.Extract(ctx)
}
