package middleware

import (
	"runtime/debug"

	. "github.com/Bofry/host-fasthttp/internal"
	"github.com/Bofry/host-fasthttp/response"
)

var _ RequestHandleModule = new(LoggingHandleModule)

type LoggingHandleModule struct {
	successor      RequestHandleModule
	loggingService LoggingService
}

func (h *LoggingHandleModule) CanSetSuccessor() bool {
	return true
}

func (h *LoggingHandleModule) SetSuccessor(successor RequestHandleModule) {
	h.successor = successor
}

func (h *LoggingHandleModule) ProcessRequest(ctx *RequestCtx, recover *RecoverService) {
	if h.successor != nil {
		eventLog := h.loggingService.CreateEventLog()
		eventLog.WriteRequest(ctx)

		recover.
			Defer(func(err interface{}) {
				resp := h.getResponse(ctx)
				if err != nil {
					defer func() {
						if resp != nil {
							eventLog.WriteResponse(ctx, resp.Flag())
						} else {
							eventLog.WriteError(ctx, err, debug.Stack())
						}
						eventLog.Flush()
					}()

					// NOTE: we should not handle error here, due to the underlying RequestHandler
					// will handle it.
				} else {
					if resp != nil {
						eventLog.WriteResponse(ctx, resp.Flag())
					} else {
						eventLog.WriteResponse(ctx, response.UNKNOWN)
					}
					eventLog.Flush()
				}
			}).
			Do(func() {
				h.successor.ProcessRequest(ctx, recover)
			})
	}
}

func (h *LoggingHandleModule) getResponse(ctx *RequestCtx) response.Response {
	obj := ctx.UserValue(response.RESPONSE_INVARIANT_NAME)
	v, ok := obj.(response.Response)
	if ok {
		return v
	}
	return nil
}
