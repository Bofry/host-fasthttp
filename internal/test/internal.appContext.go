package test

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/Bofry/host"
	fasthttp "github.com/Bofry/host-fasthttp"
	"github.com/Bofry/trace"
)

var (
	_ host.App                    = new(App)
	_ host.AppStaterConfigurator  = new(App)
	_ host.AppTracingConfigurator = new(App)
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

		// jaeger
		JaegerTraceUrl string `yaml:"jaegerTraceUrl"`
	}

	ServiceProvider struct {
		CacheClient *CacheServer
	}
)

func (app *App) Init() {
	fmt.Println("App.Init()")

	app.Component = &MockComponent{}
	app.ComponentRunner = &MockComponentRunner{prefix: "MockComponentRunner"}
}

func (app *App) OnInit() {
}

func (app *App) OnInitComplete() {
}

func (app *App) OnStart(ctx context.Context) {
}

func (app *App) OnStop(ctx context.Context) {
}

func (app *App) ConfigureLogger(logger *log.Logger) {
}

func (app *App) ConfigureTracerProvider() {
	tp, err := trace.JaegerProvider(app.Config.JaegerTraceUrl,
		trace.ServiceName("fasthttp-trace-demo"),
		trace.Environment("go-test"),
		trace.Pid(),
	)
	if err != nil {
		log.Fatal(err)
	}

	trace.SetTracerProvider(tp)
}

func (app *App) TracerProvider() *trace.SeverityTracerProvider {
	return trace.GetTracerProvider()
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
