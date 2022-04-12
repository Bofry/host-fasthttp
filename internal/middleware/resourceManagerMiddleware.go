package middleware

import (
	"github.com/Bofry/host"
	. "github.com/Bofry/host-fasthttp/internal"
	"github.com/Bofry/structproto"
)

var _ host.Middleware = new(ResourceManagerMiddleware)

type ResourceManagerMiddleware struct {
	ResourceManager interface{}
}

func (m *ResourceManagerMiddleware) Init(appCtx *host.AppContext) {
	var (
		fasthttphost = asFasthttpHost(appCtx.Host())
		preparer     = NewFasthttpHostPreparer(fasthttphost)
	)

	binder := &ResourceManagerBinder{
		router:     preparer.Router(),
		appContext: appCtx,
	}

	err := m.performBindResourceManager(m.ResourceManager, binder)
	if err != nil {
		panic(err)
	}
}

func (m *ResourceManagerMiddleware) performBindResourceManager(target interface{}, binder *ResourceManagerBinder) error {
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
