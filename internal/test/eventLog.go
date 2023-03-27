package test

import (
	"fmt"
	"log"

	fasthttp "github.com/Bofry/host-fasthttp"
	"github.com/Bofry/host-fasthttp/response"
	"github.com/Bofry/trace"
)

var _ fasthttp.EventLog = new(EventLog)

type EventLog struct {
	logger *log.Logger
}

func (l *EventLog) WriteError(ctx *fasthttp.RequestCtx, err interface{}, stackTrace []byte) {
	l.logger.Printf("EventLog.WriteError(): %v\n", err)
}

func (l *EventLog) WriteRequest(ctx *fasthttp.RequestCtx) {
	sp := trace.SpanFromContext(ctx)

	trace := fmt.Sprintf("%s-%s",
		sp.TraceID(),
		sp.SpanID())

	l.logger.Printf("EventLog.WriteRequest(): (%s) %s %s\n", trace, ctx.Method(), ctx.Path())
}

func (l *EventLog) WriteResponse(ctx *fasthttp.RequestCtx, flag response.ResponseFlag) {
	sp := trace.SpanFromContext(ctx)

	trace := fmt.Sprintf("%s-%s",
		sp.TraceID(),
		sp.SpanID())

	l.logger.Printf("EventLog.WriteResponse(): (%s) %d [%v]\n", trace, ctx.Response.StatusCode(), flag)
}

func (l *EventLog) Flush() {
	l.logger.Println("EventLog.Flush()")
}
