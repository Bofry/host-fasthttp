package middleware

import (
	"github.com/Bofry/host"
	"github.com/Bofry/trace"

	. "github.com/Bofry/host-fasthttp/internal"
)

var _ host.Middleware = new(TracingMiddleware)

type TracingMiddleware struct {
	Enabled bool
}

// Init implements internal.Middleware
func (m *TracingMiddleware) Init(app *host.AppModule) {
	var (
		fasthttphost = asFasthttpHost(app.Host())
		registrar    = NewFasthttpHostRegistrar(fasthttphost)

		tp *trace.SeverityTracerProvider
	)

	if m.Enabled {
		tp = fasthttphost.TracerProvider
	} else {
		var err error
		tp, err = trace.NoopProvider()
		if err != nil {
			FasthttpHostLogger.Fatalf("%v", err)
		}
	}

	if tp == nil {
		FasthttpHostLogger.Fatal("cannot found valid TracerProvider")
	}

	tracingHandleModule := &TracingHandleModule{
		tp: tp,
	}
	registrar.RegisterRequestHandleModule(tracingHandleModule)
	registrar.RegisterRequestHandlerReprocessModule(tracingHandleModule)
}
