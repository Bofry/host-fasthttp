package test

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

type EchoResource struct {
}

func (r *EchoResource) Init() {
}

func (r *EchoResource) Send(ctx *fasthttp.RequestCtx) {
	args := ctx.QueryArgs()
	input := args.Peek("input")
	fmt.Fprintf(ctx, "ECHO: %s", string(input))
}
