package test

import (
	"fmt"

	"github.com/Bofry/host-fasthttp/response/failure"
	"github.com/valyala/fasthttp"
)

type AccidentRequest struct {
}

func (r *AccidentRequest) Occur(ctx *fasthttp.RequestCtx) {
	panic("an error occurred")
}

func (r *AccidentRequest) Fail(ctx *fasthttp.RequestCtx) {
	panic(fmt.Errorf("FAIL"))
}

func (r *AccidentRequest) Fail2(ctx *fasthttp.RequestCtx) {
	failure.ThrowFailureMessage(failure.UNKNOWN_ERROR, "an error occurred")
}
