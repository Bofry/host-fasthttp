package middleware

import (
	"github.com/Bofry/host"
	"github.com/Bofry/trace"

	. "github.com/Bofry/host-fasthttp/internal"
)

var _ host.Middleware = new(TracingMiddleware)

type TracingMiddleware struct {
	TracerProvider *trace.SeverityTracerProvider
}

// Init implements internal.Middleware
func (m *TracingMiddleware) Init(appCtx *host.AppContext) {
	var (
		fasthttphost = asFasthttpHost(appCtx.Host())
		registrar    = NewFasthttpHostRegistrar(fasthttphost)
	)

	tracingHandleModule := &TracingHandleModule{
		tp: m.TracerProvider,
	}
	registrar.RegisterRequestHandleModule(tracingHandleModule)
	registrar.RegisterRequestHandlerReprocessModule(tracingHandleModule)
}
