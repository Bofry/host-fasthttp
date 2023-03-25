package middleware

import (
	"github.com/Bofry/host"

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
	)

	registrar.EnableTracer(m.Enabled)
}
