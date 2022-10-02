package test

import (
	"fmt"

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
