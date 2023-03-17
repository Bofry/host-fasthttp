package middleware

import (
	"github.com/Bofry/host"
	. "github.com/Bofry/host-fasthttp/internal"
)

type UnhandledRequestHandlerMiddleware struct {
	Handler RequestHandler
}

func (m *UnhandledRequestHandlerMiddleware) Init(appCtx *host.AppContext) {
	var (
		fasthttphost = asFasthttpHost(appCtx.Host())
		registrar    = NewFasthttpHostRegistrar(fasthttphost)
	)

	registrar.SetUnhandledRequestHandler(m.Handler)
}
