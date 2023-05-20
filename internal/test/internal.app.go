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
	"go.opentelemetry.io/otel/propagation"
)

var (
	logger *log.Logger = log.New(log.Writer(), "[host-fasthttp-test] ", log.LstdFlags|log.Lmsgprefix|log.LUTC)
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
	{
		logger.Printf("stoping TracerProvider")
		tp := trace.GetTracerProvider()
		err := tp.Shutdown(ctx)
		if err != nil {
			logger.Printf("stoping TracerProvider error: %+v", err)
		}
	}
}

func (app *App) ConfigureLogger(l *log.Logger) {
	l.SetFlags(logger.Flags())
	l.SetOutput(logger.Writer())
}

func (app *App) Logger() *log.Logger {
	return logger
}

func (app *App) ConfigureTracerProvider() {
	if len(app.Config.JaegerTraceUrl) == 0 {
		tp, _ := trace.NoopProvider()
		trace.SetTracerProvider(tp)
		return
	}

	tp, err := trace.JaegerProvider(app.Config.JaegerTraceUrl,
		trace.ServiceName("fasthttp-trace-demo"),
		trace.Environment("go-test"),
		trace.Pid(),
	)
	if err != nil {
		logger.Fatal(err)
	}

	trace.SetTracerProvider(tp)
}

func (app *App) TracerProvider() *trace.SeverityTracerProvider {
	return trace.GetTracerProvider()
}

func (app *App) ConfigureTextMapPropagator() {
	trace.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
}

func (app *App) TextMapPropagator() propagation.TextMapPropagator {
	return trace.GetTextMapPropagator()
}

func (provider *ServiceProvider) Init(conf *Config) {
	provider.CacheClient = &CacheServer{
		Host:     conf.RedisHost,
		Password: conf.RedisPassword,
		DB:       conf.RedisDB,
		PoolSize: conf.RedisPoolSize,
	}
}

func (provider *ServiceProvider) TracerProvider() *trace.SeverityTracerProvider {
	return trace.GetTracerProvider()
}

func (provider *ServiceProvider) TextMapPropagator() propagation.TextMapPropagator {
	return trace.GetTextMapPropagator()
}

func (provider *ServiceProvider) Logger() *log.Logger {
	return logger
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
