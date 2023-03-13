package middleware

import (
	"context"
	"fmt"

	. "github.com/Bofry/host-fasthttp/internal"
	"github.com/Bofry/trace"
)

var _ RequestHandleModule = new(TracingHandleModule)

type TracingHandleModule struct {
	tp        *trace.SeverityTracerProvider
	successor RequestHandleModule
}

// CanSetSuccessor implements RequestHandleModule
func (h *TracingHandleModule) CanSetSuccessor() bool {
	return true
}

// SetSuccessor implements RequestHandleModule
func (h *TracingHandleModule) SetSuccessor(successor RequestHandleModule) {
	h.successor = successor
}

// ProcessRequest implements RequestHandleModule
func (h *TracingHandleModule) ProcessRequest(ctx *RequestCtx, recover *RecoverService) {
	if h.successor != nil {
		fmt.Printf("Tracing Request: %s %s\n", string(ctx.Request.Header.Method()), string(ctx.Request.URI().Path()))

		recover.
			Defer(func(err interface{}) {
				if err != nil {
					defer func() {
						fmt.Printf("Tracing Error: %s %s [%v]\n", string(ctx.Request.Header.Method()), string(ctx.Request.URI().Path()), err)
					}()

					// NOTE: we should not handle error here, due to the underlying RequestHandler
					// will handle it.
				} else {
					fmt.Printf("Tracing Response: %s %s [%v]\n", string(ctx.Request.Header.Method()), string(ctx.Request.URI().Path()), ctx.Response.StatusCode())
				}
			}).
			Do(func() {
				h.successor.ProcessRequest(ctx, recover)
			})
	}
}

// OnInitComplete implements RequestHandleModule
func (h *TracingHandleModule) OnInitComplete() {
	// ignored
}

// OnStop implements RequestHandleModule
func (h *TracingHandleModule) OnStop(ctx context.Context) error {
	return h.tp.Shutdown(ctx)
}
