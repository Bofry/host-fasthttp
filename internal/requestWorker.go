package internal

import (
	"context"
)

type RequestWorker struct {
	requestHandleService *RequestHandleService
	routeResolveService  *RouteResolveService
	router               Router

	errorHandler            ErrorHandler
	unhandledRequestHandler RequestHandler
	rewriteHandler          RewriteHandler
}

func NewRequestWorker() *RequestWorker {
	return &RequestWorker{
		requestHandleService: NewRequestHandleService(),
		routeResolveService:  NewRouteResolveService(),
		router:               make(Router),
	}
}

func (w *RequestWorker) ProcessRequest(ctx *RequestCtx) {
	w.requestHandleService.ProcessRequest(ctx, new(RecoverService))
}

func (w *RequestWorker) internalProcessRequest(ctx *RequestCtx, recover *RecoverService) {
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

func (w *RequestWorker) init() {
	// register the default RequestHandleModule
	requestHandleModule := NewRequestWorkerHandleModule(w)
	w.requestHandleService.Register(requestHandleModule)
	// register the default RouteResolver
	w.routeResolveService.Register(RouteResolveModuleInstance)
}

func (w *RequestWorker) rewriteRequest(ctx *RequestCtx, path *RoutePath) *RoutePath {
	handler := w.rewriteHandler
	if handler != nil {
		return handler(ctx, path)
	}
	return path
}

func (h *RequestWorker) processError(ctx *RequestCtx, err interface{}) {
	if h.errorHandler != nil {
		h.errorHandler(ctx, err)
	}
}

func (w *RequestWorker) processUnhandledRequest(ctx *RequestCtx) {
	handler := w.unhandledRequestHandler
	if handler != nil {
		handler(ctx)
	} else {
		ctx.SetStatusCode(StatusNotFound)
	}
}

func (w *RequestWorker) start(ctx context.Context) {
	err := w.requestHandleService.triggerStart(ctx)
	if err != nil {
		FasthttpHostLogger.Fatalf("%+v", err)
	}
}

func (w *RequestWorker) stop(ctx context.Context) {
	for err := range w.requestHandleService.triggerStop(ctx) {
		if err != nil {
			FasthttpHostLogger.Printf("%+v", err)
		}
	}
}
