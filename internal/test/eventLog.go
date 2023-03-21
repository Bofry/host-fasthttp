package test

import (
	"log"

	fasthttp "github.com/Bofry/host-fasthttp"
	"github.com/Bofry/host-fasthttp/response"
)

var _ fasthttp.EventLog = new(EventLog)

type EventLog struct {
	logger *log.Logger
}

func (l *EventLog) WriteError(ctx *fasthttp.RequestCtx, err interface{}, stackTrace []byte) {
	l.logger.Printf("EventLog.WriteError(): %v\n", err)
}

func (l *EventLog) WriteRequest(ctx *fasthttp.RequestCtx) {
	l.logger.Printf("EventLog.WriteRequest(): %s %s\n", string(ctx.Method()), string(ctx.Path()))
}

func (l *EventLog) WriteResponse(ctx *fasthttp.RequestCtx, flag response.ResponseFlag) {
	l.logger.Printf("EventLog.WriteResponse(): %d [%v]\n", ctx.Response.StatusCode(), flag)
}

func (l *EventLog) Flush() {
	l.logger.Println("EventLog.Flush()")
}
