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
func (m *RequestManagerMiddleware) Init(appCtx *host.AppContext) {
	var (
		fasthttphost = asFasthttpHost(appCtx.Host())
		registrar    = NewFasthttpHostRegistrar(fasthttphost)
	)

	// register RequestManager offer FasthttpHost processing later.
	registrar.SetRequestManager(m.RequestManager)

	// binding RequestManager
	binder := &RequestManagerBinder{
		registrar:  registrar,
		appContext: appCtx,
	}

	err := m.performBindRequestManager(m.RequestManager, binder)
	if err != nil {
		panic(err)
	}
}

func (m *RequestManagerMiddleware) performBindRequestManager(target interface{}, binder *RequestManagerBinder) error {
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
