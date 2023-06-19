package test

import (
	"context"

	"github.com/Bofry/host-fasthttp/response"
	"github.com/Bofry/host-fasthttp/tracing"
	"github.com/valyala/fasthttp"
)

type TracingRequest struct {
	counter *TracingPingCounter
}

func (r *TracingRequest) Init() {
	r.counter = new(TracingPingCounter)
}

func (r *TracingRequest) Ping(ctx *fasthttp.RequestCtx) {
	sp := tracing.SpanFromRequestCtx(ctx)
	sp.Info("TracingRequest example starting")

	r.counter.increase(sp.Context())

	response.Success(ctx, "text/plain", []byte("OK"))
}

type TracingPingCounter struct {
	count int
}

func (c *TracingPingCounter) increase(ctx context.Context) int {
	tr := tracing.GetTracer(c)
	sp := tr.Start(ctx, "increase()")
	defer sp.End()

	c.count++
	return c.count
}
