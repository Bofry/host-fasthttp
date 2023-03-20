package internal

import (
	"io"
	"reflect"

	"github.com/Bofry/host"
)

var _ host.HostService = new(FasthttpHostService)

type FasthttpHostService struct{}

// ConfigureLogger implements host.HostService
func (s *FasthttpHostService) ConfigureLogger(logflags int, w io.Writer) {
	FasthttpHostLogger.SetFlags(logflags)
	FasthttpHostLogger.SetOutput(w)
}

// Init implements host.HostService
func (s *FasthttpHostService) Init(h host.Host, app *host.AppModule) {
	if v, ok := h.(*FasthttpHost); ok {
		v.preInit()
	}
}

// InitComplete implements host.HostService
func (s *FasthttpHostService) InitComplete(h host.Host, app *host.AppModule) {
	if v, ok := h.(*FasthttpHost); ok {
		// TODO: 註冊 tracer 到 request handler
		v.init()
	}
}

// DescribeHostType implements host.HostService
func (s *FasthttpHostService) DescribeHostType() reflect.Type {
	return typeOfHost
}
