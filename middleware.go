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

func UseLogging(service LoggingService) host.Middleware {
	if service == nil {
		panic("argument 'service' cannot be nil")
	}

	return &middleware.LoggingMiddleware{
		LoggingService: service,
	}
}

func UseResourceManager(resourceManager interface{}) host.Middleware {
	if resourceManager == nil {
		panic("argument 'resourceManager' cannot be nil")
	}

	return &middleware.ResourceManagerMiddleware{
		ResourceManager: resourceManager,
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

func UseXHttpMethodHeader(headers ...string) host.Middleware {
	return &middleware.XHttpMethodHeaderMiddleware{
		Headers: headers,
	}
}
