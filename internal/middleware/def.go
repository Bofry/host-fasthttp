package middleware

import (
	"log"
	"reflect"

	"github.com/Bofry/host-fasthttp/internal"
	"github.com/Bofry/host-fasthttp/response"
)

var (
	typeOfHost           = reflect.TypeOf(internal.FasthttpHost{})
	typeOfRequestHandler = reflect.TypeOf(internal.RequestHandler(nil))

	TAG_URL = "url"
)

type (
	LoggingService interface {
		CreateEventLog() EventLog
		ConfigureLogger(l *log.Logger)
	}

	EventLog interface {
		WriteRequest(ctx *internal.RequestCtx)
		WriteError(ctx *internal.RequestCtx, err interface{}, stackTrace []byte)
		WriteResponse(ctx *internal.RequestCtx, flag response.ResponseFlag)
		Flush()
	}
)
