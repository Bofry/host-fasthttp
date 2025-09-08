package test

import (
	"github.com/Bofry/host-fasthttp/response"
	"github.com/valyala/fasthttp"
)

type ResponseCustomHeaderRequest struct {
}

func (r *ResponseCustomHeaderRequest) Get(ctx *fasthttp.RequestCtx) {

	ctx.Response.Header.Set("My-Custom-Header-Get", "foo")

	response.Success(ctx, "text/plain", []byte("OK"))
}

func (r *ResponseCustomHeaderRequest) Post(ctx *fasthttp.RequestCtx) {

	ctx.Response.Header.Set("My-Custom-Header-Post", "foo")

	response.Success(ctx, "text/plain", []byte("OK"))
}
