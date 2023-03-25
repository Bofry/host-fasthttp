package test

import (
	"github.com/Bofry/host-fasthttp/response"
	"github.com/Bofry/trace"
	"github.com/valyala/fasthttp"
)

type TracingRequest struct {
}

func (r *TracingRequest) Ping(ctx *fasthttp.RequestCtx) {
	sp := trace.SpanFromContext(ctx)
	sp.Info("TracingRequest example starting")

	response.Success(ctx, "text/plain", []byte("OK"))
}
