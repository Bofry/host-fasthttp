package tracingutil

import (
	"net/textproto"

	http "github.com/valyala/fasthttp"
	"go.opentelemetry.io/otel/propagation"
)

var (
	_ propagation.TextMapCarrier = RequestHeaderCarrier{}
)

// RequestHeaderCarrier adapts fasthttp.RequestHeader to satisfy the TextMapCarrier interface.
type RequestHeaderCarrier struct {
	header *http.RequestHeader
}

func NewRequestHeaderCarrier(header *http.RequestHeader) *RequestHeaderCarrier {
	return &RequestHeaderCarrier{
		header: header,
	}
}

// Get returns the value associated with the passed key.
func (hc RequestHeaderCarrier) Get(key string) string {
	var value = string(hc.header.Peek(key))
	if len(value) == 0 {
		// NOTE: patch compatibility
		// https://www.w3.org/TR/trace-context/#header-name
		canonicalHeaderKey := textproto.CanonicalMIMEHeaderKey(key)
		value = string(hc.header.Peek(canonicalHeaderKey))
	}
	return value
}

// Set stores the key-value pair.
func (hc RequestHeaderCarrier) Set(key string, value string) {
	hc.header.Set(key, value)
}

// Keys lists the keys stored in this carrier.
func (hc RequestHeaderCarrier) Keys() []string {
	header := hc.header
	keys := make([]string, 0, header.Len())
	header.VisitAll(func(key, value []byte) {
		keys = append(keys, string(key))
	})
	return keys
}
