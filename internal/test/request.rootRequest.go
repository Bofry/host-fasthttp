package test

import (
	"github.com/Bofry/host-fasthttp/response"
	"github.com/Bofry/trace"
	"github.com/valyala/fasthttp"
)

type RootRequest struct {
	Tracer *trace.SeverityTracer
}

func (r *RootRequest) Ping(ctx *fasthttp.RequestCtx) {

	response.Success(ctx, "text/plain", []byte("Pong"))
}
