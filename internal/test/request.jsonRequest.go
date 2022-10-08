package test

import (
	fasthttp "github.com/Bofry/host-fasthttp"
	"github.com/Bofry/host-fasthttp/response"
)

type JsonRequest struct {
}

func (r *JsonRequest) Init() {
}

func (r *JsonRequest) Ping(ctx *fasthttp.RequestCtx) {
	response.Json.Success(ctx, struct {
		Message string `json:"message"`
	}{
		Message: "OK",
	})
}

func (r *JsonRequest) Fail(ctx *fasthttp.RequestCtx) {
	response.Json.Failure(ctx, struct {
		Message string `json:"message"`
	}{
		Message: "UNKNOWN_ERROR",
	}, fasthttp.StatusBadRequest)
}
