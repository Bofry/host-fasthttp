package test

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

type SettingRequest struct {
	ServiceProvider *ServiceProvider
	name            string
}

var peekFormat = `Redis:
    Host: %s
    Password: %s
    DB: %d
    PoolSize: %d
From: %s`

func (r *SettingRequest) Init() {
	fmt.Println("SettingResource.Init()")
	r.name = "SettingResource"
}

func (r *SettingRequest) Peek(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, peekFormat,
		r.ServiceProvider.CacheClient.Host,
		r.ServiceProvider.CacheClient.Password,
		r.ServiceProvider.CacheClient.DB,
		r.ServiceProvider.CacheClient.PoolSize,
		r.name,
	)
}
