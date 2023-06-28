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

	TAG_URL            = "url"
	TAG_OPT_EXPAND_ENV = "@ExpandEnv"
)

type (
	LoggingService interface {
		CreateEventLog(ev EventEvidence) EventLog
		ConfigureLogger(l *log.Logger)
	}

	EventLog interface {
		OnError(ctx *internal.RequestCtx, err interface{}, stackTrace []byte)
		OnProcessRequest(ctx *internal.RequestCtx)
		OnProcessRequestComplete(ctx *internal.RequestCtx, flag response.ResponseFlag)
		Flush()
	}
)
