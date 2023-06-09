package test

import (
	"github.com/Bofry/host-fasthttp/response"
	"github.com/Bofry/host-fasthttp/tracing"
	"github.com/valyala/fasthttp"
)

type TracingRequest struct {
}

func (r *TracingRequest) Ping(ctx *fasthttp.RequestCtx) {
	sp := tracing.SpanFromRequestCtx(ctx)
	sp.Info("TracingRequest example starting")

	response.Success(ctx, "text/plain", []byte("OK"))
}
