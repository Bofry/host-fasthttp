package test

import (
	"log"

	fasthttp "github.com/Bofry/host-fasthttp"
	"github.com/Bofry/host-fasthttp/internal/responseutil"
)

var _ fasthttp.LoggingService = new(MultiLoggerService)

type MultiLoggerService struct {
	LoggingServices []fasthttp.LoggingService
}

func (s *MultiLoggerService) CreateEventLog(ev fasthttp.EventEvidence) fasthttp.EventLog {
	var eventlogs []fasthttp.EventLog
	for _, svc := range s.LoggingServices {
		eventlogs = append(eventlogs, svc.CreateEventLog(ev))
	}

	return MultiEventLog{
		EventLogs: eventlogs,
	}
}

func (s *MultiLoggerService) ConfigureLogger(l *log.Logger) {
	for _, svc := range s.LoggingServices {
		svc.ConfigureLogger(l)
	}
}

var _ fasthttp.EventLog = MultiEventLog{}

type MultiEventLog struct {
	EventLogs []fasthttp.EventLog
}

// Flush implements middleware.EventLog.
func (l MultiEventLog) Flush() {
	for _, log := range l.EventLogs {
		log.Flush()
	}
}

// WriteError implements middleware.EventLog.
func (l MultiEventLog) WriteError(ctx *fasthttp.RequestCtx, err interface{}, stackTrace []byte) {
	for _, log := range l.EventLogs {
		log.WriteError(ctx, err, stackTrace)
	}
}

// WriteRequest implements middleware.EventLog.
func (l MultiEventLog) WriteRequest(ctx *fasthttp.RequestCtx) {
	for _, log := range l.EventLogs {
		log.WriteRequest(ctx)
	}
}

// WriteResponse implements middleware.EventLog.
func (l MultiEventLog) WriteResponse(ctx *fasthttp.RequestCtx, flag responseutil.ResponseFlag) {
	for _, log := range l.EventLogs {
		log.WriteResponse(ctx, flag)
	}
}
