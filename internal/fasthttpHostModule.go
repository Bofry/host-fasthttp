package internal

import (
	"io"
	"reflect"

	"github.com/Bofry/host"
)

var _ host.HostModule = FasthttpHostModule{}

type FasthttpHostModule struct{}

// ConfigureLogger implements host.HostService
func (FasthttpHostModule) ConfigureLogger(logflags int, w io.Writer) {
	FasthttpHostLogger.SetFlags(logflags)
	FasthttpHostLogger.SetOutput(w)
}

// Init implements host.HostService
func (FasthttpHostModule) Init(h host.Host, app *host.AppModule) {
	if v, ok := h.(*FasthttpHost); ok {
		v.preInit()
		v.setTracerProvider(app.TracerProvider())
		v.setTextMapPropagator(app.TextMapPropagator())
		v.setLogger(app.Logger())
	}
}

// InitComplete implements host.HostService
func (FasthttpHostModule) InitComplete(h host.Host, app *host.AppModule) {
	if v, ok := h.(*FasthttpHost); ok {
		v.init()
	}
}

// DescribeHostType implements host.HostService
func (FasthttpHostModule) DescribeHostType() reflect.Type {
	return typeOfHost
}
