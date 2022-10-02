package test

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

type EchoRequest struct {
}

func (r *EchoRequest) Init() {
}

func (r *EchoRequest) Send(ctx *fasthttp.RequestCtx) {
	args := ctx.QueryArgs()
	input := args.Peek("input")
	fmt.Fprintf(ctx, "ECHO: %s", string(input))
}
