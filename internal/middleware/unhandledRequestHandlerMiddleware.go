package middleware

import (
	"github.com/Bofry/host"
	. "github.com/Bofry/host-fasthttp/internal"
)

type UnhandledRequestHandlerMiddleware struct {
	Handler RequestHandler
}

func (m *UnhandledRequestHandlerMiddleware) Init(app *host.AppModule) {
	var (
		fasthttphost = asFasthttpHost(app.Host())
		registrar    = NewFasthttpHostRegistrar(fasthttphost)
	)

	registrar.SetUnhandledRequestHandler(m.Handler)
}
