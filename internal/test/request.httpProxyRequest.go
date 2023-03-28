package test

import (
	"io/ioutil"
	"net/http"

	fasthttp "github.com/Bofry/host-fasthttp"
	"github.com/Bofry/host-fasthttp/response"
	"github.com/Bofry/trace"
	"go.opentelemetry.io/otel/propagation"
)

type HttpProxyRequest struct {
	ServiceProvider *ServiceProvider
}

func (r *HttpProxyRequest) POST(ctx *fasthttp.RequestCtx) {
	sp := trace.SpanFromContext(ctx)

	client := &http.Client{}
	req, err := http.NewRequest("DO", ctx.URI().String(), nil)
	if err != nil {
		panic(err)
	}
	carrier := propagation.HeaderCarrier(req.Header)
	sp.Inject(r.ServiceProvider.TextMapPropagator(), carrier)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	response.Success(ctx, resp.Header.Get("Content-Type"), body)
}

func (r *HttpProxyRequest) DO(ctx *fasthttp.RequestCtx) {
	response.Success(ctx, "text/plain", []byte("OK"))
}
