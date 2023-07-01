package middleware

import (
	"github.com/Bofry/host-fasthttp/internal/responseutil"
	"github.com/valyala/fasthttp"
)

var _ EventLog = NoopEventLog(0)

type NoopEventLog int

// Flush implements EventLog.
func (NoopEventLog) Flush() {}

// OnError implements EventLog.
func (NoopEventLog) OnError(ctx *fasthttp.RequestCtx, err interface{}, stackTrace []byte) {}

// OnProcessRequest implements EventLog.
func (NoopEventLog) OnProcessRequest(ctx *fasthttp.RequestCtx) {}

// OnProcessRequestComplete implements EventLog.
func (NoopEventLog) OnProcessRequestComplete(ctx *fasthttp.RequestCtx, flag responseutil.ResponseFlag) {
}
