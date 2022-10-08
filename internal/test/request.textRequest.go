package test

import (
	fasthttp "github.com/Bofry/host-fasthttp"
	"github.com/Bofry/host-fasthttp/response"
)

type TextRequest struct {
}

func (r *TextRequest) Init() {
}

func (r *TextRequest) Ping(ctx *fasthttp.RequestCtx) {
	response.Text.Success(ctx, "OK")
}

func (r *TextRequest) Fail(ctx *fasthttp.RequestCtx) {
	response.Text.Failure(ctx, "UNKNOWN_ERROR", fasthttp.StatusBadRequest)
}
