package middleware

import (
	. "github.com/Bofry/host-fasthttp/internal"
)

var _ RouteResolveModule = new(XHttpMethodHeaderRouteResolveModule)

type XHttpMethodHeaderRouteResolveModule struct {
	headers []string

	successor RouteResolver
}

func (m *XHttpMethodHeaderRouteResolveModule) CanSetSuccessor() bool {
	return true
}

func (m *XHttpMethodHeaderRouteResolveModule) SetSuccessor(successor RouteResolver) {
	m.successor = successor
}

func (m *XHttpMethodHeaderRouteResolveModule) ResolveHttpMethod(ctx *RequestCtx) string {
	var method = ctx.Request.Header.Peek(HEADER_XHTTP_METHOD)
	if method != nil {
		return string(method)
	}

	// read available http method headers and find method
	for _, h := range m.headers {
		method = ctx.Request.Header.Peek(h)
		if method != nil {
			return string(method)
		}
	}
	// pass to successor, if does not find available http method from header
	return m.successor.ResolveHttpMethod(ctx)
}

func (m *XHttpMethodHeaderRouteResolveModule) ResolveHttpPath(ctx *RequestCtx) string {
	return m.successor.ResolveHttpPath(ctx)
}
