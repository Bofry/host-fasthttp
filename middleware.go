package fasthttp

import (
	"github.com/Bofry/host"
	"github.com/Bofry/host-fasthttp/internal/middleware"
)

func UseErrorHandler(handler ErrorHandler) host.Middleware {
	if handler == nil {
		panic("argument 'handler' cannot be nil")
	}

	return &middleware.ErrorHandlerMiddleware{
		Handler: handler,
	}
}

func UseLogging(services ...LoggingService) host.Middleware {
	if len(services) == 0 {
		return &middleware.LoggingMiddleware{
			LoggingService: middleware.NoopLoggingServiceSingleton,
		}
	}

	return &middleware.LoggingMiddleware{
		LoggingService: middleware.NewCompositeLoggingService(services...),
	}
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

func UseTracing(enabled bool) host.Middleware {
	return &middleware.TracingMiddleware{
		Enabled: enabled,
	}
}

func UseXHttpMethodHeader(headers ...string) host.Middleware {
	return &middleware.XHttpMethodHeaderMiddleware{
		Headers: headers,
	}
}
