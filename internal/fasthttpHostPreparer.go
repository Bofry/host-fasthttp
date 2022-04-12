package internal

type FasthttpHostPreparer struct {
	subject *FasthttpHost
}

func NewFasthttpHostPreparer(subject *FasthttpHost) *FasthttpHostPreparer {
	return &FasthttpHostPreparer{
		subject: subject,
	}
}

func (p *FasthttpHostPreparer) RegisterRequestHandleModule(successor RequestHandleModule) {
	p.subject.requestWorker.requestHandleService.Register(successor)
}

func (p *FasthttpHostPreparer) RegisterRouteResolveModule(successor RouteResolveModule) {
	p.subject.requestWorker.routeResolveService.Register(successor)
}

func (p *FasthttpHostPreparer) RegisterErrorHandler(handler ErrorHandler) {
	p.subject.requestWorker.errorHandler = handler
}

func (p *FasthttpHostPreparer) RegisterRewriteHandler(handler RewriteHandler) {
	p.subject.requestWorker.rewriteHandler = handler
}

func (p *FasthttpHostPreparer) RegisterUnhandledRequestHandler(handler RequestHandler) {
	p.subject.requestWorker.unhandledRequestHandler = handler
}

func (p *FasthttpHostPreparer) Router() Router {
	return p.subject.requestWorker.router
}
