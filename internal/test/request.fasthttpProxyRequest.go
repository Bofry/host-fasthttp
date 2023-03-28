package test

import (
	fasthttp "github.com/Bofry/host-fasthttp"
	"github.com/Bofry/host-fasthttp/internal/tracingutil"
	"github.com/Bofry/host-fasthttp/response"
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
	carrier := tracingutil.NewRequestHeaderCarrier(&req.Header)
	sp.Inject(r.ServiceProvider.TextMapPropagator(), carrier)

	http.Do(req, resp)

	response.Success(ctx,
		string(resp.Header.ContentType()),
		resp.Body())
}

func (r *FasthttpProxyRequest) DO(ctx *fasthttp.RequestCtx) {
	response.Success(ctx, "text/plain", []byte("OK"))
}
