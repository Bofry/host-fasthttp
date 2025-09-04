package middleware

import (
	"github.com/Bofry/host"
	. "github.com/Bofry/host-fasthttp/internal"
)

var _ host.Middleware = new(CorsMiddleware)

type CorsMiddleware struct {
}

// Init implements internal.Middleware.
func (c *CorsMiddleware) Init(app *host.AppModule) {
	var (
		fasthttphost = asFasthttpHost(app.Host())
		registrar    = NewFasthttpHostRegistrar(fasthttphost)
	)

	registrar.SetCorsHeader(true)
}
