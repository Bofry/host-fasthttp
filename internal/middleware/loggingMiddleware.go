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
func (m *LoggingMiddleware) Init(app *host.AppModule) {
	var (
		fasthttphost = asFasthttpHost(app.Host())
		registrar    = NewFasthttpHostRegistrar(fasthttphost)
	)

	m.LoggingService.ConfigureLogger(fasthttphost.Logger)

	loggingHandleModule := &LoggingHandleModule{
		loggingService: m.LoggingService,
	}
	registrar.RegisterRequestHandleModule(loggingHandleModule)
}
