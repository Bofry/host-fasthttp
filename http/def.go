package http

import (
	_ "unsafe"

	"github.com/valyala/fasthttp"
)

const (
	MethodGet     = fasthttp.MethodGet     // RFC 7231, 4.3.1
	MethodHead    = fasthttp.MethodHead    // RFC 7231, 4.3.2
	MethodPost    = fasthttp.MethodPost    // RFC 7231, 4.3.3
	MethodPut     = fasthttp.MethodPut     // RFC 7231, 4.3.4
	MethodPatch   = fasthttp.MethodPatch   // RFC 5789
	MethodDelete  = fasthttp.MethodDelete  // RFC 7231, 4.3.5
	MethodConnect = fasthttp.MethodConnect // RFC 7231, 4.3.6
	MethodOptions = fasthttp.MethodOptions // RFC 7231, 4.3.7
	MethodTrace   = fasthttp.MethodTrace   // RFC 7231, 4.3.8
)

type (
	Args = fasthttp.Args

	HttpClientOption interface {
		apply(req *fasthttp.Request, resp *fasthttp.Response) error
	}
)

//go:linkname AcquireRequest github.com/valyala/fasthttp.AcquireRequest
func AcquireRequest() *fasthttp.Request

//go:linkname ReleaseRequest github.com/valyala/fasthttp.ReleaseRequest
func ReleaseRequest(req *fasthttp.Request)

//go:linkname AcquireResponse github.com/valyala/fasthttp.AcquireResponse
func AcquireResponse() *fasthttp.Response

//go:linkname ReleaseResponse github.com/valyala/fasthttp.ReleaseResponse
func ReleaseResponse(resp *fasthttp.Response)
