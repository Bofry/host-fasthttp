package middleware

import (
	"github.com/Bofry/host"
	. "github.com/Bofry/host-fasthttp/internal"
)

var _ host.Middleware = new(RewriterMiddleware)

type RewriterMiddleware struct {
	Handler RewriteHandler
}

// Init implements internal.Middleware
func (m *RewriterMiddleware) Init(app *host.AppModule) {
	var (
		fasthttphost = asFasthttpHost(app.Host())
		registrar    = NewFasthttpHostRegistrar(fasthttphost)
	)

	registrar.SetRewriteHandler(m.Handler)
}
