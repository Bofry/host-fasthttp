package tracing

import (
	_ "unsafe"

	"github.com/Bofry/host-fasthttp/internal/tracingutil"
	http "github.com/valyala/fasthttp"
)

type (
	RequestHeaderCarrier = tracingutil.RequestHeaderCarrier
)

//go:linkname NewRequestHeaderCarrier github.com/Bofry/host-fasthttp/internal/tracingutil.NewRequestHeaderCarrier
func NewRequestHeaderCarrier(header *http.RequestHeader) *RequestHeaderCarrier
