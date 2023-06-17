package internal

type FasthttpHostRegistrar struct {
	host *FasthttpHost
}

func NewFasthttpHostRegistrar(host *FasthttpHost) *FasthttpHostRegistrar {
	return &FasthttpHostRegistrar{
		host: host,
	}
}

func (r *FasthttpHostRegistrar) RegisterRequestHandleModule(module RequestHandleModule) {
	r.host.requestHandleService.Register(module)
}

func (r *FasthttpHostRegistrar) RegisterRouteResolveModule(module RouteResolveModule) {
	r.host.requestWorker.RouteResolveService.Register(module)
}

func (r *FasthttpHostRegistrar) EnableTracer(enabled bool) {
	r.host.requestTracerService.Enabled = enabled
}

func (r *FasthttpHostRegistrar) SetErrorHandler(handler ErrorHandler) {
	r.host.requestWorker.ErrorHandler = handler
}

func (r *FasthttpHostRegistrar) SetRewriteHandler(handler RewriteHandler) {
	r.host.requestWorker.RewriteHandler = handler
}

func (r *FasthttpHostRegistrar) SetUnhandledRequestHandler(handler RequestHandler) {
	r.host.requestWorker.UnhandledRequestHandler = handler
}

func (r *FasthttpHostRegistrar) SetRequestManager(requestManager interface{}) {
	r.host.requestManager = requestManager
}

func (r *FasthttpHostRegistrar) AddRoute(method string, path string, handler RequestHandler, handlerComponentID string) {
	r.host.requestWorker.Router.Add(method, path, handler, handlerComponentID)
}
