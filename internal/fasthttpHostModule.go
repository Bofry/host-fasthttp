package internal

import (
	"io"
	"reflect"

	"github.com/Bofry/host"
)

var _ host.HostModule = new(FasthttpHostModule)

type FasthttpHostModule struct{}

// ConfigureLogger implements host.HostService
func (s *FasthttpHostModule) ConfigureLogger(logflags int, w io.Writer) {
	FasthttpHostLogger.SetFlags(logflags)
	FasthttpHostLogger.SetOutput(w)
}

// Init implements host.HostService
func (s *FasthttpHostModule) Init(h host.Host, app *host.AppModule) {
	if v, ok := h.(*FasthttpHost); ok {
		v.preInit()
		v.TracerProvider = app.TracerProvider()
	}
}

// InitComplete implements host.HostService
func (s *FasthttpHostModule) InitComplete(h host.Host, app *host.AppModule) {
	if v, ok := h.(*FasthttpHost); ok {
		// TODO: 註冊 tracer 到 request handler
		v.init()
	}
}

// DescribeHostType implements host.HostService
func (s *FasthttpHostModule) DescribeHostType() reflect.Type {
	return typeOfHost
}
