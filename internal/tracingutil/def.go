package tracingutil

import (
	"github.com/Bofry/host-fasthttp/internal/requestutil"
	"github.com/Bofry/trace"
	http "github.com/valyala/fasthttp"
)

func ExtractSpan(ctx *http.RequestCtx) *trace.SeveritySpan {
	obj := requestutil.ExtractSpan(ctx)
	v, ok := obj.(*trace.SeveritySpan)
	if ok {
		return v
	}
	return nil
}
