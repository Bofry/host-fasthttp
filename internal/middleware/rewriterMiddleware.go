package middleware

import (
	"github.com/Bofry/host"
	. "github.com/Bofry/host-fasthttp/internal"
)

var _ host.Middleware = new(RewriterMiddleware)

type RewriterMiddleware struct {
	Handler RewriteHandler
}

func (m *RewriterMiddleware) Init(appCtx *host.AppContext) {
	var (
		fasthttphost = asFasthttpHost(appCtx.Host())
		preparer     = NewFasthttpHostPreparer(fasthttphost)
	)

	preparer.RegisterRewriteHandler(m.Handler)
}
