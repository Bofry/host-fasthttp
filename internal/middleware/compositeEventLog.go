package middleware

import (
	"github.com/Bofry/host-fasthttp/internal/responseutil"
	"github.com/valyala/fasthttp"
)

var _ EventLog = CompositeEventLog{}

type CompositeEventLog struct {
	eventLogs []EventLog
}

// Flush implements EventLog.
func (l CompositeEventLog) Flush() {
	for _, log := range l.eventLogs {
		log.Flush()
	}
}

// OnError implements EventLog.
func (l CompositeEventLog) OnError(ctx *fasthttp.RequestCtx, err interface{}, stackTrace []byte) {
	for _, log := range l.eventLogs {
		log.OnError(ctx, err, stackTrace)
	}
}

// OnProcessRequest implements EventLog.
func (l CompositeEventLog) OnProcessRequest(ctx *fasthttp.RequestCtx) {
	for _, log := range l.eventLogs {
		log.OnProcessRequest(ctx)
	}
}

// OnProcessRequestComplete implements EventLog.
func (l CompositeEventLog) OnProcessRequestComplete(ctx *fasthttp.RequestCtx, flag responseutil.ResponseFlag) {
	for _, log := range l.eventLogs {
		log.OnProcessRequestComplete(ctx, flag)
	}
}
