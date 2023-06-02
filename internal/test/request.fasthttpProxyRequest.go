package test

import (
	fasthttp "github.com/Bofry/host-fasthttp"
	"github.com/Bofry/host-fasthttp/http"
	"github.com/Bofry/host-fasthttp/response"
	"github.com/Bofry/trace"
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
	req.Header.SetMethod("DO")

	http.Do(req, resp,
		http.WithTracePropagation(sp.Context(), r.ServiceProvider.TextMapPropagator()),
	)

	response.SendSuccess(ctx, resp)
}

func (r *FasthttpProxyRequest) DO(ctx *fasthttp.RequestCtx) {
	response.Success(ctx, "text/plain", []byte("OK"))
}
