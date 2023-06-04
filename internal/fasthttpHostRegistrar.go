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
	p.host.requestHandleService.Register(module)
}

func (p *FasthttpHostRegistrar) RegisterRouteResolveModule(module RouteResolveModule) {
	p.host.routeResolveService.Register(module)
}

func (p *FasthttpHostRegistrar) EnableTracer(enabled bool) {
	p.host.requestTracerService.Enabled = enabled
}

func (p *FasthttpHostRegistrar) SetErrorHandler(handler ErrorHandler) {
	p.host.requestWorker.ErrorHandler = handler
}

func (p *FasthttpHostRegistrar) SetRewriteHandler(handler RewriteHandler) {
	p.host.requestWorker.RewriteHandler = handler
}

func (p *FasthttpHostRegistrar) SetUnhandledRequestHandler(handler RequestHandler) {
	p.host.requestWorker.UnhandledRequestHandler = handler
}

func (p *FasthttpHostRegistrar) SetRequestManager(requestManager interface{}) {
	p.host.requestManager = requestManager
}

func (p *FasthttpHostRegistrar) AddRoute(method string, path string, handler RequestHandler, handlerComponentID string) {
	p.host.router.Add(method, path, handler, handlerComponentID)
}
