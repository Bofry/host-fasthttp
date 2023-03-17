package middleware

import (
	"github.com/Bofry/host"
	. "github.com/Bofry/host-fasthttp/internal"
)

var _ host.Middleware = new(XHttpMethodHeaderMiddleware)

type XHttpMethodHeaderMiddleware struct {
	Headers []string
}

// Init implements internal.Middleware
func (m *XHttpMethodHeaderMiddleware) Init(appCtx *host.AppContext) {
	var (
		fasthttphost = asFasthttpHost(appCtx.Host())
		registrar    = NewFasthttpHostRegistrar(fasthttphost)
	)

	routeResolveModule := &XHttpMethodHeaderRouteResolveModule{
		headers: m.Headers,
	}
	registrar.RegisterRouteResolveModule(routeResolveModule)
}
