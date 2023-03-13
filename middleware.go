package fasthttp

import (
	"github.com/Bofry/host"
	"github.com/Bofry/host-fasthttp/internal/middleware"
	"github.com/Bofry/trace"
)

func UseErrorHandler(handler ErrorHandler) host.Middleware {
	if handler == nil {
		panic("argument 'handler' cannot be nil")
	}

	return &middleware.ErrorHandlerMiddleware{
		Handler: handler,
	}
}

func UseLogging(service LoggingService) host.Middleware {
	if service == nil {
		panic("argument 'service' cannot be nil")
	}

	return &middleware.LoggingMiddleware{
		LoggingService: service,
	}
}

func UseResourceManager(resourceManager interface{}) host.Middleware {
	panic("method is obsolete. Use UseRequestManager() instead.")
}

func UseRequestManager(requestManager interface{}) host.Middleware {
	if requestManager == nil {
		panic("argument 'requestManager' cannot be nil")
	}

	return &middleware.RequestManagerMiddleware{
		RequestManager: requestManager,
	}
}

func UseRewriter(handler RewriteHandler) host.Middleware {
	if handler == nil {
		panic("argument 'handler' cannot be nil")
	}

	return &middleware.RewriterMiddleware{
		Handler: handler,
	}
}

func UseUnhandledRequestHandler(handler RequestHandler) host.Middleware {
	if handler == nil {
		panic("argument 'handler' cannot be nil")
	}

	return &middleware.UnhandledRequestHandlerMiddleware{
		Handler: handler,
	}
}

func UseTracing(tp *trace.SeverityTracerProvider) host.Middleware {
	if tp == nil {
		panic("argument 'tp' cannot be nil")
	}

	return &middleware.TracingMiddleware{
		TracerProvider: tp,
	}
}

func UseXHttpMethodHeader(headers ...string) host.Middleware {
	return &middleware.XHttpMethodHeaderMiddleware{
		Headers: headers,
	}
}
