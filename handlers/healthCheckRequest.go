package handlers

import (
	fasthttp "github.com/Bofry/host-fasthttp"
	"github.com/Bofry/host-fasthttp/response"
	"github.com/Bofry/host-fasthttp/tracing"
)

type HealthCheckRequest struct{}

func (r *HealthCheckRequest) Ping(ctx *fasthttp.RequestCtx) {
	tracing.SpanFromRequestCtx(ctx).Disable(true)

	response.Text.Success(ctx, "PONG")
}

func (r *HealthCheckRequest) Head(ctx *fasthttp.RequestCtx) {
	tracing.SpanFromRequestCtx(ctx).Disable(true)

	response.Text.Success(ctx, "OK")
}

func (r *HealthCheckRequest) Get(ctx *fasthttp.RequestCtx) {
	tracing.SpanFromRequestCtx(ctx).Disable(true)

	response.Text.Success(ctx, "OK")
}
