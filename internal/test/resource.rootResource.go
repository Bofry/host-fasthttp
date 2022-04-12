package test

import (
	"github.com/Bofry/host-fasthttp/response"
	"github.com/valyala/fasthttp"
)

type RootResource struct {
}

func (r *RootResource) Ping(ctx *fasthttp.RequestCtx) {

	response.Success(ctx, "text/plain", []byte("Pong"))
}
