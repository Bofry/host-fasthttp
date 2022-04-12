package test

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

type AccidentResource struct {
}

func (r *AccidentResource) Occur(ctx *fasthttp.RequestCtx) {
	panic("an error occurred")
}

func (r *AccidentResource) Fail(ctx *fasthttp.RequestCtx) {
	panic(fmt.Errorf("FAIL"))
}
