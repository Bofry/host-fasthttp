package internal

import (
	"context"
	"log"
	"os"
	"reflect"

	"github.com/Bofry/host-fasthttp/internal/tracingutil"
	"github.com/Bofry/trace"
	"github.com/valyala/fasthttp"
	"go.opentelemetry.io/otel/propagation"
)

const (
	DEFAULT_HTTP_PORT   = "80"
	HEADER_XHTTP_METHOD = "X-Http-Method"

	StatusNotFound = 404

	LOGGER_PREFIX = "[host-fasthttp] "
)

var (
	typeOfHost               = reflect.TypeOf(FasthttpHost{})
	defaultTracerProvider    = createNoopTracerProvider()
	defaultTextMapPropagator = createNoopTextMapPropagator()
	defaultSpanExtractor     = tracingutil.RequestCtxSpanExtractor(0)
	unhandledRequestSpanName = "unknown path"

	FasthttpHostServiceInstance = new(FasthttpHostModule)

	FasthttpHostLogger *log.Logger = log.New(os.Stdout, LOGGER_PREFIX, log.LstdFlags|log.Lmsgprefix)
)

// import
type (
	Server         = fasthttp.Server
	RequestHandler = fasthttp.RequestHandler
	RequestCtx     = fasthttp.RequestCtx
)

// interface
type (
	RouteResolver interface {
		ResolveHttpMethod(ctx *RequestCtx) string
		ResolveHttpPath(ctx *RequestCtx) string
	}

	RouteResolveModule interface {
		RouteResolver

		CanSetSuccessor() bool
		SetSuccessor(successor RouteResolver)
	}

	RequestHandleModule interface {
		CanSetSuccessor() bool
		SetSuccessor(successor RequestHandleModule)
		ProcessRequest(ctx *RequestCtx, state RequestState, recover *RecoverService)
		OnInitComplete()
		OnStart(ctx context.Context) error
		OnStop(ctx context.Context) error
	}
)

// function
type (
	ErrorHandler   func(ctx *RequestCtx, err interface{})
	RewriteHandler func(ctx *RequestCtx, path *RoutePath) *RoutePath
)

func createNoopTracerProvider() *trace.SeverityTracerProvider {
	tp, err := trace.NoopProvider()
	if err != nil {
		FasthttpHostLogger.Fatalf("cannot create NoopProvider: %v", err)
	}
	return tp
}

func createNoopTextMapPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator()
}
