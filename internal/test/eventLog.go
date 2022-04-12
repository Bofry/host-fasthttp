package test

import (
	"fmt"

	fasthttp "github.com/Bofry/host-fasthttp"
	"github.com/Bofry/host-fasthttp/response"
)

var _ fasthttp.EventLog = new(EventLog)

type EventLog struct{}

func (l *EventLog) WriteError(ctx *fasthttp.RequestCtx, err interface{}, stackTrace []byte) {
	fmt.Println("EventLog.WriteError()")
}

func (l *EventLog) WriteRequest(ctx *fasthttp.RequestCtx) {
	fmt.Println("EventLog.WriteRequest()")
}

func (l *EventLog) WriteResponse(ctx *fasthttp.RequestCtx, flag response.ResponseFlag) {
	fmt.Println("EventLog.WriteResponse()")
}

func (l *EventLog) Flush() {
	fmt.Println("EventLog.Flush()")
}
