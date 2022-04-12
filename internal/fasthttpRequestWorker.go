package internal

type FasthttpRequestWorker struct {
	requestHandleService *RequestHandleService
	routeResolveService  *RouteResolveService
	router               Router

	errorHandler            ErrorHandler
	unhandledRequestHandler RequestHandler
	rewriteHandler          RewriteHandler
}

func NewFasthttpRequestWorker() *FasthttpRequestWorker {
	return &FasthttpRequestWorker{
		requestHandleService: new(RequestHandleService),
		routeResolveService:  new(RouteResolveService),
		router:               make(Router),
	}
}

func (w *FasthttpRequestWorker) ProcessRequest(ctx *RequestCtx) {
	w.requestHandleService.ProcessRequest(ctx, new(RecoverService))
}

func (w *FasthttpRequestWorker) internalProcessRequest(ctx *RequestCtx, recover *RecoverService) {
	var (
		method = w.routeResolveService.ResolveHttpMethod(ctx)
		path   = w.routeResolveService.ResolveHttpPath(ctx)
	)

	routePath := &RoutePath{
		Method: method,
		Path:   path,
	}

	recover.
		Defer(func(err interface{}) {
			if err != nil {
				w.processError(ctx, err)
			}
		}).
		Do(func() {
			routePath = w.rewriteRequest(ctx, routePath)
			if routePath == nil {
				panic("invalid RoutePath. The RouttPath should not be nil.")
			}

			handler := w.router.Get(routePath.Method, routePath.Path)
			if handler != nil {
				handler(ctx)
			} else {
				w.processUnhandledRequest(ctx)
			}
		})
}

func (w *FasthttpRequestWorker) init() {
	// register the default RequestHandleModule
	requestHandleModule := NewRequestHandleModuleImpl(w)
	w.requestHandleService.Register(requestHandleModule)
	// register the default RouteResolver
	w.routeResolveService.Register(RouteResolveModuleInstance)
}

func (w *FasthttpRequestWorker) rewriteRequest(ctx *RequestCtx, path *RoutePath) *RoutePath {
	handler := w.rewriteHandler
	if handler != nil {
		return handler(ctx, path)
	}
	return path
}

func (h *FasthttpRequestWorker) processError(ctx *RequestCtx, err interface{}) {
	if h.errorHandler != nil {
		h.errorHandler(ctx, err)
	}
}

func (w *FasthttpRequestWorker) processUnhandledRequest(ctx *RequestCtx) {
	handler := w.unhandledRequestHandler
	if handler != nil {
		handler(ctx)
	} else {
		ctx.SetStatusCode(StatusNotFound)
	}
}
