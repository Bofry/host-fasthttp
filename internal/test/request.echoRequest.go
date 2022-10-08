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

	/* NOTE:
	 *   the formal response result should as following:
	 *
	 *     buf := fmt.Sprintf("ECHO: %s", string(input))
	 *     response.Text.Success(ctx, buf)
	 */
	fmt.Fprintf(ctx, "ECHO: %s", string(input))
}
