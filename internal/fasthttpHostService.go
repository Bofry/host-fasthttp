package internal

import (
	"reflect"

	"github.com/Bofry/host"
)

var _ host.HostService = new(FasthttpHostService)

type FasthttpHostService struct{}

func (s *FasthttpHostService) Init(h host.Host, app *host.AppContext) {
	if v, ok := h.(*FasthttpHost); ok {
		v.preInit()
	}
}

func (s *FasthttpHostService) InitComplete(h host.Host, app *host.AppContext) {
	if v, ok := h.(*FasthttpHost); ok {
		v.init()
	}
}

func (s *FasthttpHostService) GetHostType() reflect.Type {
	return typeOfHost
}
