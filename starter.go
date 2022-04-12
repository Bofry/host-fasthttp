package fasthttp

import (
	"github.com/Bofry/host"
	"github.com/Bofry/host-fasthttp/internal"
)

func Startup(app interface{}) *host.Starter {
	var (
		starter = host.Startup(app)
	)

	host.RegisterHostService(starter, internal.FasthttpHostServiceInstance)

	return starter
}
