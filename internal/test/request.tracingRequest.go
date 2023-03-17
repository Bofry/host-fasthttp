package test

import (
	"context"

	"github.com/Bofry/host-fasthttp/response"
	"github.com/Bofry/trace"
	"github.com/valyala/fasthttp"
)

type TracingRequest struct {
	Tracer *trace.SeverityTracer
}

func (r *TracingRequest) Ping(ctx *fasthttp.RequestCtx) {
	sp := r.Tracer.Open(context.Background(), "PING /Tracing")
	defer sp.End()
	sp.Info("example starting")

	response.Success(ctx, "text/plain", []byte("OK"))
}
