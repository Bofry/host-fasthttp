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
		preparer     = NewFasthttpHostPreparer(fasthttphost)
	)

	preparer.RegisterUnhandledRequestHandler(m.Handler)
}
