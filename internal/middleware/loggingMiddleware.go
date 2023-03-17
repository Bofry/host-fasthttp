package middleware

import (
	"github.com/Bofry/host"
	. "github.com/Bofry/host-fasthttp/internal"
)

var _ host.Middleware = new(LoggingMiddleware)

type LoggingMiddleware struct {
	LoggingService LoggingService
}

// Init implements internal.Middleware
func (m *LoggingMiddleware) Init(appCtx *host.AppContext) {
	var (
		fasthttphost = asFasthttpHost(appCtx.Host())
		registrar    = NewFasthttpHostRegistrar(fasthttphost)
	)

	loggingHandleModule := &LoggingHandleModule{
		loggingService: m.LoggingService,
	}
	registrar.RegisterRequestHandleModule(loggingHandleModule)
}
