package test

import (
	"fmt"
	"log"

	fasthttp "github.com/Bofry/host-fasthttp"
	"github.com/Bofry/host-fasthttp/response"
)

var _ fasthttp.EventLog = EventLog{}

type EventLog struct {
	logger   *log.Logger
	evidence fasthttp.EventEvidence
}

func (l EventLog) OnError(ctx *fasthttp.RequestCtx, err interface{}, stackTrace []byte) {
	l.logger.Printf("EventLog.OnError(): %v\n", err)
}

func (l EventLog) OnProcessRequest(ctx *fasthttp.RequestCtx) {
	traceID := fmt.Sprintf("%s-%s",
		l.evidence.RequestTraceID(),
		l.evidence.RequestSpanID())

	l.logger.Printf("EventLog.OnProcessRequest(): (%s) %s %s\n", traceID, ctx.Method(), ctx.Path())
}

func (l EventLog) OnProcessRequestComplete(ctx *fasthttp.RequestCtx, flag response.ResponseFlag) {
	traceID := fmt.Sprintf("%s-%s",
		l.evidence.RequestTraceID(),
		l.evidence.RequestSpanID())

	l.logger.Printf("EventLog.OnProcessRequestComplete(): (%s) %d [%v]\n", traceID, ctx.Response.StatusCode(), flag)
}

func (l EventLog) Flush() {
	l.logger.Println("EventLog.Flush()")
}
