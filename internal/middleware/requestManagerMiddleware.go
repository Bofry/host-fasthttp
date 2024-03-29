package middleware

import (
	"github.com/Bofry/host"
	. "github.com/Bofry/host-fasthttp/internal"
	"github.com/Bofry/structproto"
)

var _ host.Middleware = new(RequestManagerMiddleware)

type RequestManagerMiddleware struct {
	RequestManager interface{}
}

// Init implements internal.Middleware
func (m *RequestManagerMiddleware) Init(app *host.AppModule) {
	var (
		fasthttphost = asFasthttpHost(app.Host())
		registrar    = NewFasthttpHostRegistrar(fasthttphost)
	)

	// register RequestManager offer FasthttpHost processing later.
	registrar.SetRequestManager(m.RequestManager)

	// binding RequestManager
	binder := &RequestManagerBinder{
		registrar: registrar,
		app:       app,
	}

	err := m.bindRequestManager(m.RequestManager, binder)
	if err != nil {
		panic(err)
	}
}

func (m *RequestManagerMiddleware) bindRequestManager(target interface{}, binder *RequestManagerBinder) error {
	prototype, err := structproto.Prototypify(target,
		&structproto.StructProtoResolveOption{
			TagName:     TAG_URL,
			TagResolver: UrlTagResolver,
		},
	)
	if err != nil {
		return err
	}

	return prototype.Bind(binder)
}
