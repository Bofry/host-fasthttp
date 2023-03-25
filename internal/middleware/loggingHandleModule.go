package middleware

import (
	"context"
	"runtime/debug"

	. "github.com/Bofry/host-fasthttp/internal"
	"github.com/Bofry/host-fasthttp/response"
)

var (
	_ RequestHandleModule = new(LoggingHandleModule)
)

type LoggingHandleModule struct {
	successor      RequestHandleModule
	loggingService LoggingService
}

// CanSetSuccessor implements RequestHandleModule
func (h *LoggingHandleModule) CanSetSuccessor() bool {
	return true
}

// SetSuccessor implements RequestHandleModule
func (h *LoggingHandleModule) SetSuccessor(successor RequestHandleModule) {
	h.successor = successor
}

// ProcessRequest implements RequestHandleModule
func (h *LoggingHandleModule) ProcessRequest(ctx *RequestCtx, recover *RecoverService) {
	if h.successor != nil {
		eventLog := h.loggingService.CreateEventLog()
		eventLog.WriteRequest(ctx)

		recover.
			Defer(func(err interface{}) {
				resp := response.ExtractResponseState(ctx)
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
			Do(func(Finalizer) {
				h.successor.ProcessRequest(ctx, recover)
			})
	}
}

// OnInitComplete implements RequestHandleModule
func (*LoggingHandleModule) OnInitComplete() {
	// ignored
}

// OnStart implements RequestHandleModule
func (*LoggingHandleModule) OnStart(ctx context.Context) error {
	// do nothing
	return nil
}

// OnStop implements RequestHandleModule
func (*LoggingHandleModule) OnStop(ctx context.Context) error {
	// do nothing
	return nil
}
