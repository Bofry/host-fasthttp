package internal

import (
	"log"
	"reflect"

	"github.com/Bofry/host"
)

var _ host.HostService = new(FasthttpHostService)

type FasthttpHostService struct{}

// ConfigureLogger implements host.HostService
func (s *FasthttpHostService) ConfigureLogger(logger *log.Logger) {
	FasthttpHostLogger.SetFlags(logger.Flags())
	FasthttpHostLogger.SetOutput(logger.Writer())
}

// Init implements host.HostService
func (s *FasthttpHostService) Init(h host.Host, app *host.AppContext) {
	if v, ok := h.(*FasthttpHost); ok {
		v.preInit()
	}
}

// InitComplete implements host.HostService
func (s *FasthttpHostService) InitComplete(h host.Host, app *host.AppContext) {
	if v, ok := h.(*FasthttpHost); ok {
		// TODO: 註冊 tracer 到 request handler
		v.init()
	}
}

// DescribeHostType implements host.HostService
func (s *FasthttpHostService) DescribeHostType() reflect.Type {
	return typeOfHost
}
