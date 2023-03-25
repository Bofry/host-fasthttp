package internal

import (
	http "github.com/valyala/fasthttp"
	"go.opentelemetry.io/otel/propagation"
)

var (
	_ propagation.TextMapCarrier = RequestHeaderCarrier{}
)

// RequestHeaderCarrier adapts fasthttp.RequestHeader to satisfy the TextMapCarrier interface.
type RequestHeaderCarrier struct {
	ctx *http.RequestCtx
}

// Get returns the value associated with the passed key.
func (hc RequestHeaderCarrier) Get(key string) string {
	req := hc.ctx
	return string(req.Request.Header.Peek(key))
}

// Set stores the key-value pair.
func (hc RequestHeaderCarrier) Set(key string, value string) {
	req := hc.ctx
	req.Request.Header.Set(key, value)
}

// Keys lists the keys stored in this carrier.
func (hc RequestHeaderCarrier) Keys() []string {
	req := hc.ctx
	keys := make([]string, 0, req.Request.Header.Len())
	req.Request.Header.VisitAll(func(key, value []byte) {
		keys = append(keys, string(key))
	})
	return keys
}
