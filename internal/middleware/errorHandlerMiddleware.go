package middleware

import (
	"github.com/Bofry/host"
	. "github.com/Bofry/host-fasthttp/internal"
)

var _ host.Middleware = new(ErrorHandlerMiddleware)

type ErrorHandlerMiddleware struct {
	Handler ErrorHandler
}

// Init implements internal.Middleware
func (m *ErrorHandlerMiddleware) Init(app *host.AppModule) {
	var (
		fasthttphost = asFasthttpHost(app.Host())
		registrar    = NewFasthttpHostRegistrar(fasthttphost)
	)

	registrar.SetErrorHandler(m.Handler)
}
