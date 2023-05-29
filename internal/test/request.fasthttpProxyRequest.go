package test

import (
	fasthttp "github.com/Bofry/host-fasthttp"
	"github.com/Bofry/host-fasthttp/response"
	"github.com/Bofry/host-fasthttp/tracing"
	"github.com/Bofry/trace"
	http "github.com/valyala/fasthttp"
)

type FasthttpProxyRequest struct {
	ServiceProvider *ServiceProvider
}

func (r *FasthttpProxyRequest) POST(ctx *fasthttp.RequestCtx) {
	sp := trace.SpanFromContext(ctx)

	req := http.AcquireRequest()
	resp := http.AcquireResponse()
	defer http.ReleaseRequest(req)
	defer http.ReleaseResponse(resp)

	req.SetRequestURI(ctx.URI().String())
	req.Header.DisableNormalizing()
	req.Header.SetMethod("DO")
	carrier := tracing.NewRequestHeaderCarrier(&req.Header)
	sp.Inject(r.ServiceProvider.TextMapPropagator(), carrier)

	http.Do(req, resp)

	response.SendSuccess(ctx, resp)
}

func (r *FasthttpProxyRequest) DO(ctx *fasthttp.RequestCtx) {
	response.Success(ctx, "text/plain", []byte("OK"))
}
