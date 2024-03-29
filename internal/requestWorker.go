package internal

import (
	"context"
	"fmt"

	"github.com/Bofry/host-fasthttp/internal/requestutil"
	"github.com/Bofry/host-fasthttp/internal/responseutil"
	"github.com/Bofry/host-fasthttp/internal/tracingutil"
	"github.com/Bofry/trace"
)

type RequestWorker struct {
	RequestHandleService *RequestHandleService
	RequestTracerService *RequestTracerService
	RouteResolveService  *RouteResolveService
	Router               Router

	OnHostErrorProc OnHostErrorHandler

	ErrorHandler            ErrorHandler
	UnhandledRequestHandler RequestHandler
	RewriteHandler          RewriteHandler
}

func (w *RequestWorker) ProcessRequest(ctx *RequestCtx) {
	var (
		routePath = w.resolveRoutePath(ctx)
	)

	// rewriting url
	routePath = w.rewriteRequest(ctx, routePath)
	if routePath == nil {
		panic("invalid RoutePath. The RoutePath should not be nil.")
	}

	// start tracing
	var (
		componentID = w.Router.FindHandlerComponentID(routePath.Method, routePath.Path)
		carrier     = tracingutil.NewRequestHeaderCarrier(&ctx.Request.Header)

		spanName string = unhandledRequestSpanName
		tr       *trace.SeverityTracer
		sp       *trace.SeveritySpan
	)

	if w.Router.Has(*routePath) {
		spanName = routePath.String()
	}

	tr = w.RequestTracerService.Tracer(componentID)
	sp = tr.ExtractWithPropagator(
		ctx,
		w.RequestTracerService.TextMapPropagator(),
		carrier,
		spanName)
	defer sp.End()

	requestState := RequestState{
		RoutePath: routePath,
		Tracer:    tr,
		Span:      sp,
	}

	w.RequestHandleService.ProcessRequest(ctx, requestState, new(Recover))
}

func (w *RequestWorker) internalProcessRequest(ctx *RequestCtx, state RequestState, recover *Recover) {
	recover.
		Defer(func(err interface{}) {
			if err != nil {
				w.processError(ctx, err)
			}
		}).
		Do(func(finalizer Finalizer) {
			var (
				tr        *trace.SeverityTracer = state.Tracer
				sp        *trace.SeveritySpan   = state.Span
				routePath *RoutePath            = state.RoutePath
			)

			// set Tracer and Span
			requestutil.InjectTracer(ctx, tr)
			requestutil.InjectSpan(ctx, sp)

			finalizer.Add(func(err interface{}) {
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
					trace.HttpResponse(ctx.Response.String()),
					trace.HttpStatusCode(ctx.Response.StatusCode()),
				)

				var (
					resFlag = responseutil.ExtractResponseState(ctx)
				)

				if resFlag != nil {
					switch resFlag.Flag() {
					case responseutil.SUCCESS:
						sp.Reply(trace.PASS, nil)
					case responseutil.FAILURE:
						sp.Reply(trace.FAIL, nil)
					}
				}
			})

			sp.Tags(
				trace.HttpRequest(ctx.Request.String()),
				trace.HttpMethod(string(ctx.Request.Header.Method())),
				trace.HttpRequestPath(string(ctx.Request.URI().Path())),
				trace.HttpUserAgent(string(ctx.Request.Header.UserAgent())),
			)

			handler := w.Router.Get(routePath.Method, routePath.Path)
			if handler != nil {
				handler(ctx)
			} else {
				w.processUnhandledRequest(ctx)
			}
		})
}

func (w *RequestWorker) init() {
	// register the default RequestHandleModule
	stdRequestHandleModule := NewStdRequestHandleModule(w)
	w.RequestHandleService.Register(stdRequestHandleModule)
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
		ctx.NotFound()
	}
}

func (w *RequestWorker) start(ctx context.Context) {
	err := w.RequestHandleService.triggerStart(ctx)
	if err != nil {
		var disposed bool = false
		if w.OnHostErrorProc != nil {
			disposed = w.OnHostErrorProc(err)
		}
		if !disposed {
			FasthttpHostLogger.Fatalf("%+v", err)
		}
	}
}

func (w *RequestWorker) stop(ctx context.Context) {
	for err := range w.RequestHandleService.triggerStop(ctx) {
		if err != nil {
			FasthttpHostLogger.Printf("%+v", err)
		}
	}
}

func (w *RequestWorker) resolveRoutePath(ctx *RequestCtx) *RoutePath {
	var (
		method = w.RouteResolveService.ResolveHttpMethod(ctx)
		path   = w.RouteResolveService.ResolveHttpPath(ctx)
	)

	return &RoutePath{
		Method: method,
		Path:   path,
	}
}
