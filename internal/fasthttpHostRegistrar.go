package internal

type FasthttpHostRegistrar struct {
	host *FasthttpHost
}

func NewFasthttpHostRegistrar(host *FasthttpHost) *FasthttpHostRegistrar {
	return &FasthttpHostRegistrar{
		host: host,
	}
}

func (p *FasthttpHostRegistrar) RegisterRequestHandleModule(module RequestHandleModule) {
	p.host.requestWorker.requestHandleService.Register(module)
}

func (p *FasthttpHostRegistrar) RegisterRequestHandlerReprocessModule(module RequestResourceProcessModule) {
	p.host.requestResourceProcessService.Register(module)
}

func (p *FasthttpHostRegistrar) RegisterRouteResolveModule(module RouteResolveModule) {
	p.host.requestWorker.routeResolveService.Register(module)
}

func (p *FasthttpHostRegistrar) SetErrorHandler(handler ErrorHandler) {
	p.host.requestWorker.errorHandler = handler
}

func (p *FasthttpHostRegistrar) SetRewriteHandler(handler RewriteHandler) {
	p.host.requestWorker.rewriteHandler = handler
}

func (p *FasthttpHostRegistrar) SetUnhandledRequestHandler(handler RequestHandler) {
	p.host.requestWorker.unhandledRequestHandler = handler
}

func (p *FasthttpHostRegistrar) SetRequestManager(requestManager interface{}) {
	p.host.requestManager = requestManager
}

func (p *FasthttpHostRegistrar) AddRoute(method string, path string, handler RequestHandler) {
	p.host.requestWorker.router.Add(method, path, handler)
}
