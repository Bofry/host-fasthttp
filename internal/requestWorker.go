package internal

import (
	"context"
	"fmt"

	"github.com/Bofry/host-fasthttp/internal/requestutil"
	"github.com/Bofry/trace"
)

type RequestWorker struct {
	RequestHandleService *RequestHandleService
	RequestTracerService *RequestTracerService
	RouteResolveService  *RouteResolveService
	Router               Router

	ErrorHandler            ErrorHandler
	UnhandledRequestHandler RequestHandler
	RewriteHandler          RewriteHandler
}

func (w *RequestWorker) ProcessRequest(ctx *RequestCtx) {
	w.RequestHandleService.ProcessRequest(ctx, new(RecoverService))
}

func (w *RequestWorker) internalProcessRequest(ctx *RequestCtx, recover *RecoverService) {
	var (
		method = w.RouteResolveService.ResolveHttpMethod(ctx)
		path   = w.RouteResolveService.ResolveHttpPath(ctx)
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
		Do(func(finalizer Finalizer) {
			routePath = w.rewriteRequest(ctx, routePath)
			if routePath == nil {
				panic("invalid RoutePath. The RouttPath should not be nil.")
			}

			handler := w.Router.Get(routePath.Method, routePath.Path)

			moduleID := w.Router.FindRequestComponentID(routePath.Method, routePath.Path)
			tr := w.RequestTracerService.Tracer(moduleID)
			carrier := RequestHeaderCarrier{ctx: ctx}
			sp := tr.ExtractWithPropagator(
				ctx,
				w.RequestTracerService.TextMapPropagator,
				carrier,
				routePath.String(),
				trace.WithNewRoot())

			finalizer.Add(func(err interface{}) {
				defer sp.End()

				if err != nil {
					if e, ok := err.(error); ok {
						sp.Err(e)
					} else if e, ok := err.(string); ok {
						sp.Err(fmt.Errorf(e))
					} else if e, ok := err.(fmt.Stringer); ok {
						sp.Err(fmt.Errorf(e.String()))
					} else {
						sp.Err(fmt.Errorf("%+v", err))
					}
				}

				sp.Tags(
					trace.Stringer("http.response", &ctx.Response),
				)
			})

			sp.Tags(
				trace.Stringer("http.request", &ctx.Request),
			)

			requestutil.InjectTracer(ctx, tr)
			requestutil.InjectSpan(ctx, sp)

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
	w.RequestHandleService.Register(requestHandleModule)
	// register the default RouteResolver
	w.RouteResolveService.Register(RouteResolveModuleInstance)
}

func (w *RequestWorker) rewriteRequest(ctx *RequestCtx, path *RoutePath) *RoutePath {
	handler := w.RewriteHandler
	if handler != nil {
		return handler(ctx, path)
	}
	return path
}

func (h *RequestWorker) processError(ctx *RequestCtx, err interface{}) {
	if h.ErrorHandler != nil {
		h.ErrorHandler(ctx, err)
	}
}

func (w *RequestWorker) processUnhandledRequest(ctx *RequestCtx) {
	handler := w.UnhandledRequestHandler
	if handler != nil {
		handler(ctx)
	} else {
		ctx.SetStatusCode(StatusNotFound)
	}
}

func (w *RequestWorker) start(ctx context.Context) {
	err := w.RequestHandleService.triggerStart(ctx)
	if err != nil {
		FasthttpHostLogger.Fatalf("%+v", err)
	}
}

func (w *RequestWorker) stop(ctx context.Context) {
	for err := range w.RequestHandleService.triggerStop(ctx) {
		if err != nil {
			FasthttpHostLogger.Printf("%+v", err)
		}
	}
}
