package middleware

import (
	"github.com/Bofry/host"
	. "github.com/Bofry/host-fasthttp/internal"
)

var _ host.Middleware = new(LoggingMiddleware)

type LoggingMiddleware struct {
	LoggingService LoggingService
}

func (m *LoggingMiddleware) Init(appCtx *host.AppContext) {
	var (
		fasthttphost = asFasthttpHost(appCtx.Host())
		preparer     = NewFasthttpHostPreparer(fasthttphost)
	)

	loggingHandleModule := &LoggingHandleModule{
		loggingService: m.LoggingService,
	}
	preparer.RegisterRequestHandleModule(loggingHandleModule)
}
