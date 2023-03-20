package internal

import (
	"context"
	"log"
	"os"
	"reflect"

	"github.com/valyala/fasthttp"
)

const (
	DEFAULT_HTTP_PORT   = "80"
	HEADER_XHTTP_METHOD = "X-Http-Method"

	StatusNotFound = 404

	LOGGER_PREFIX = "[host-fasthttp] "
)

var (
	typeOfHost = reflect.TypeOf(FasthttpHost{})

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
		ProcessRequest(ctx *RequestCtx, recover *RecoverService)
		OnInitComplete()
		OnStart(ctx context.Context) error
		OnStop(ctx context.Context) error
	}

	RequestResourceProcessModule interface {
		ProcessRequestResource(rv reflect.Value)
	}
)

// function
type (
	ErrorHandler   func(ctx *RequestCtx, err interface{})
	RewriteHandler func(ctx *RequestCtx, path *RoutePath) *RoutePath
)
