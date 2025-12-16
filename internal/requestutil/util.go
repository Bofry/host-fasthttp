package requestutil

import (
	"github.com/valyala/fasthttp"
)

const (
	USER_STORE_KEY_RESPONSE_STATE         = "github.com/Bofry/host-fasthttp/internal/response::ResponseState"
	USER_STORE_KEY_SEVERITY_TRACER        = "github.com/Bofry/trace::SeverityTracer"
	USER_STORE_KEY_SEVERITY_SPAN          = "github.com/Bofry/trace::SeveritySpan"
	USER_STORE_KEY_DISABLE_RESET_RESPONSE = "github.com/Bofry/host-fasthttp::DisableResetResponse"

	TRUE  = "true"
	FALSE = "false"
)

func InjectResponseState(ctx *fasthttp.RequestCtx, responseState interface{}) {
	ctx.SetUserValue(USER_STORE_KEY_RESPONSE_STATE, responseState)
}

func ExtractResponseState(ctx *fasthttp.RequestCtx) interface{} {
	return ctx.UserValue(USER_STORE_KEY_RESPONSE_STATE)
}

func InjectTracer(ctx *fasthttp.RequestCtx, tracer interface{}) {
	ctx.SetUserValue(USER_STORE_KEY_SEVERITY_TRACER, tracer)
}

func ExtractTracer(ctx *fasthttp.RequestCtx) interface{} {
	return ctx.UserValue(USER_STORE_KEY_SEVERITY_TRACER)
}

func InjectSpan(ctx *fasthttp.RequestCtx, span interface{}) {
	ctx.SetUserValue(USER_STORE_KEY_SEVERITY_SPAN, span)
}

func ExtractSpan(ctx *fasthttp.RequestCtx) interface{} {
	return ctx.UserValue(USER_STORE_KEY_SEVERITY_SPAN)
}

func DisableResetResponse(ctx *fasthttp.RequestCtx) {
	ctx.SetUserValue(USER_STORE_KEY_DISABLE_RESET_RESPONSE, TRUE)
}

func EnableResetResponse(ctx *fasthttp.RequestCtx) {
	ctx.SetUserValue(USER_STORE_KEY_DISABLE_RESET_RESPONSE, FALSE)
}

func CanResetResponse(ctx *fasthttp.RequestCtx) bool {
	return ctx.UserValue(USER_STORE_KEY_DISABLE_RESET_RESPONSE) == TRUE
}
