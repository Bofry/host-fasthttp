package test

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

type SettingResource struct {
	ServiceProvider *ServiceProvider
	name            string
}

var peekFormat = `Redis:
    Host: %s
    Password: %s
    DB: %d
    PoolSize: %d
From: %s`

func (r *SettingResource) Init() {
	fmt.Println("SettingResource.Init()")
	r.name = "SettingResource"
}

func (r *SettingResource) Peek(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, peekFormat,
		r.ServiceProvider.CacheClient.Host,
		r.ServiceProvider.CacheClient.Password,
		r.ServiceProvider.CacheClient.DB,
		r.ServiceProvider.CacheClient.PoolSize,
		r.name,
	)
}
