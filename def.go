package fasthttp

import (
	"github.com/Bofry/host-fasthttp/internal"
	"github.com/Bofry/host-fasthttp/internal/middleware"
)

// import
type (
	Server         = internal.Server
	RequestHandler = internal.RequestHandler
	RequestCtx     = internal.RequestCtx
)

// interface
type (
	LoggingService = middleware.LoggingService
	EventLog       = middleware.EventLog
)

// struct
type (
	Host      = internal.FasthttpHost
	RoutePath = internal.RoutePath
)

// function
type (
	ErrorHandler   = internal.ErrorHandler
	RewriteHandler = internal.RewriteHandler
)
