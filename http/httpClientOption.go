package http

import (
	"context"

	"github.com/Bofry/host-fasthttp/tracing"
	"github.com/Bofry/trace"
	"github.com/valyala/fasthttp"
	"go.opentelemetry.io/otel/propagation"
)

type HttpClientOptionProc func(req *fasthttp.Request, resp *fasthttp.Response) error

func (proc HttpClientOptionProc) apply(req *fasthttp.Request, resp *fasthttp.Response) error {
	return proc(req, resp)
}

func WithTracePropagation(ctx context.Context, propagator propagation.TextMapPropagator) HttpClientOptionProc {
	return func(req *fasthttp.Request, resp *fasthttp.Response) error {
		carrier := tracing.NewRequestHeaderCarrier(&req.Header)
		if propagator == nil {
			propagator = trace.GetTextMapPropagator()
		}
		propagator.Inject(ctx, carrier)
		return nil
	}
}
