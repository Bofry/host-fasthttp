package test

import (
	"github.com/Bofry/host-fasthttp/response"
	"github.com/valyala/fasthttp"
)

type RootRequest struct {
}

func (r *RootRequest) Ping(ctx *fasthttp.RequestCtx) {

	response.Success(ctx, "text/plain", []byte("Pong"))
}
