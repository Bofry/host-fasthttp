package test

import (
	"bytes"
	"log"

	fasthttp "github.com/Bofry/host-fasthttp"
	"github.com/Bofry/host-fasthttp/response"
)

var _ fasthttp.LoggingService = new(BlackholeLoggerService)

type BlackholeLoggerService struct {
	Buffer *bytes.Buffer
}

func (s *BlackholeLoggerService) CreateEventLog(ev fasthttp.EventEvidence) fasthttp.EventLog {
	s.Buffer.WriteString("CreateEventLog()")
	s.Buffer.WriteByte('\n')
	return &BlackholeEventLog{
		buffer: s.Buffer,
	}
}

func (*BlackholeLoggerService) ConfigureLogger(l *log.Logger) {
}

var _ fasthttp.EventLog = new(BlackholeEventLog)

type BlackholeEventLog struct {
	buffer *bytes.Buffer
}

func (l *BlackholeEventLog) OnError(ctx *fasthttp.RequestCtx, err interface{}, stackTrace []byte) {
	l.buffer.WriteString("OnError()")
	l.buffer.WriteByte('\n')
}

func (l *BlackholeEventLog) OnProcessRequest(ctx *fasthttp.RequestCtx) {
	l.buffer.WriteString("OnProcessRequest()")
	l.buffer.WriteByte('\n')
}

func (l *BlackholeEventLog) OnProcessRequestComplete(ctx *fasthttp.RequestCtx, flag response.ResponseFlag) {
	l.buffer.WriteString("OnProcessRequestComplete()")
	l.buffer.WriteByte('\n')
}

func (l *BlackholeEventLog) Flush() {
	l.buffer.WriteString("Flush()")
	l.buffer.WriteByte('\n')
}
