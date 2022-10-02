package test

import (
	"fmt"
	"runtime"
	"strings"

	fasthttp "github.com/Bofry/host-fasthttp"
)

type (
	App struct {
		Host            *Host
		Config          *Config
		ServiceProvider *ServiceProvider

		Component       *MockComponent
		ComponentRunner *MockComponentRunner
	}

	Host fasthttp.Host

	Config struct {
		// fasthttp server
		ListenAddress  string `arg:"address"`
		EnableCompress bool   `arg:"compress"`
		ServerName     string `arg:"hostname"`

		// redis
		RedisHost     string `env:"*REDIS_HOST"       yaml:"redisHost"`
		RedisPassword string `env:"*REDIS_PASSWORD"   yaml:"redisPassword"`
		RedisDB       int    `env:"REDIS_DB"          yaml:"redisDB"`
		RedisPoolSize int    `env:"REDIS_POOL_SIZE"   yaml:"redisPoolSize"`
		Workspace     string `env:"-"                 yaml:"workspace"`
	}

	ServiceProvider struct {
		CacheClient *CacheServer
	}
)

func (app *App) Init(conf *Config) {
	fmt.Println("App.Init()")

	app.Component = &MockComponent{}
	app.ComponentRunner = &MockComponentRunner{prefix: "MockComponentRunner"}
}

func (provider *ServiceProvider) Init(conf *Config) {
	provider.CacheClient = &CacheServer{
		Host:     conf.RedisHost,
		Password: conf.RedisPassword,
		DB:       conf.RedisDB,
		PoolSize: conf.RedisPoolSize,
	}
}

func (h *Host) Init(conf *Config) {
	h.Server = &fasthttp.Server{
		Name:                          conf.ServerName,
		DisableKeepalive:              true,
		DisableHeaderNamesNormalizing: true,
	}
	h.ListenAddress = conf.ListenAddress
	h.EnableCompress = conf.EnableCompress
	h.Version = strings.Join([]string{
		"v201206",
		runtime.Version(),
	}, " ")
}
