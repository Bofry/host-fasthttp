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

func (m *RequestManagerMiddleware) Init(appCtx *host.AppContext) {
	var (
		fasthttphost = asFasthttpHost(appCtx.Host())
		preparer     = NewFasthttpHostPreparer(fasthttphost)
	)

	binder := &RequestManagerBinder{
		router:     preparer.Router(),
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
