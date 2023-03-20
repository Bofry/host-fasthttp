package test

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Bofry/config"
	fasthttp "github.com/Bofry/host-fasthttp"
	"github.com/Bofry/host-fasthttp/response"
)

type RequestManager struct {
	*RootRequest     `url:"/"`
	*EchoRequest     `url:"/Echo"`
	*SettingRequest  `url:"/Setting"`
	*AccidentRequest `url:"/Accident"`
	*JsonRequest     `url:"/Json"`
	*TextRequest     `url:"/Text"`
	*TracingRequest  `url:"/Tracing"`
}

func TestApp_Sanity(t *testing.T) {
	/* like
	$ export REDIS_HOST=kubernate-redis:26379
	$ export REDIS_PASSWORD=1234
	$ export REDIS_POOL_SIZE=128
	*/
	initializeEnvironment()
	/* like
	$ go run app.go --address ":10094" --compress true --hostname "DemoService"
	*/
	initializeArgs()

	var errorBuffer bytes.Buffer

	var errorCount = 0

	app := App{}
	starter := fasthttp.Startup(&app).
		Middlewares(
			fasthttp.UseRequestManager(&RequestManager{}),
			fasthttp.UseXHttpMethodHeader(),
			fasthttp.UseErrorHandler(func(ctx *fasthttp.RequestCtx, err interface{}) {
				errorCount++
				v, ok := err.(error)
				if ok && v.Error() == "FAIL" {
					response.Failure(ctx, string(ctx.Response.Header.ContentType()), []byte("FAIL"), 400)
				}
				fmt.Fprintf(&errorBuffer, "err: %+v", err)
			}),
			fasthttp.UseLogging(&LoggingService{}),
			fasthttp.UseTracing(nil),
			fasthttp.UseRewriter(func(ctx *fasthttp.RequestCtx, path *fasthttp.RoutePath) *fasthttp.RoutePath {
				if strings.HasPrefix(path.Path, "/Echo/") {
					ctx.Request.URI().QueryArgs().Add("input", path.Path[6:])
					path.Path = "/Echo"
				}
				return path
			}),
			fasthttp.UseUnhandledRequestHandler(func(ctx *fasthttp.RequestCtx) {
				ctx.SetStatusCode(fasthttp.StatusNotFound)
			}),
		).
		ConfigureConfiguration(func(service *config.ConfigurationService) {
			service.
				LoadEnvironmentVariables("").
				LoadYamlFile("config.yaml").
				LoadCommandArguments()
		})

	runCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := starter.Start(runCtx); err != nil {
		t.Error(err)
	}

	client := &http.Client{}
	{
		req, err := http.NewRequest("POST", "http://127.0.0.1:10094", nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("X-Http-Method", "PING")
		resp, err := client.Do(req)
		if resp.StatusCode != 200 {
			t.Errorf("assert 'http.Response.StatusCode':: expected '%v', got '%v'", 200, resp.StatusCode)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if string(body) != "Pong" {
			t.Errorf("assert 'http.Response.Body':: expected '%v', got '%v'", "Pong", string(body))
		}
	}
	{
		req, err := http.NewRequest("PING", "http://127.0.0.1:10094", nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("If-None-Match", `W/"wyzzy"`)
		resp, err := client.Do(req)
		if resp.StatusCode != 200 {
			t.Errorf("assert 'http.Response.StatusCode':: expected '%v', got '%v'", 200, resp.StatusCode)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if string(body) != "Pong" {
			t.Errorf("assert 'http.Response.Body':: expected '%v', got '%v'", "Pong", string(body))
		}
	}
	{
		req, err := http.NewRequest("SEND", "http://127.0.0.1:10094/Echo?input=Can%20you%20hear%20me", nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("If-None-Match", `W/"wyzzy"`)
		resp, err := client.Do(req)
		if resp.StatusCode != 200 {
			t.Errorf("assert 'http.Response.StatusCode':: expected '%v', got '%v'", 200, resp.StatusCode)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if string(body) != "ECHO: Can you hear me" {
			t.Errorf("assert 'http.Response.Body':: expected '%v', got '%v'", "ECHO: Can you hear me", string(body))
		}
	}
	{
		req, err := http.NewRequest("SEND", "http://127.0.0.1:10094/Echo/Hello", nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("If-None-Match", `W/"wyzzy"`)
		resp, err := client.Do(req)
		if resp.StatusCode != 200 {
			t.Errorf("assert 'http.Response.StatusCode':: expected '%v', got '%v'", 200, resp.StatusCode)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if string(body) != "ECHO: Hello" {
			t.Errorf("assert 'http.Response.Body':: expected '%v', got '%v'", "ECHO: Hello", string(body))
		}
	}
	{
		req, err := http.NewRequest("PEEK", "http://127.0.0.1:10094/Setting", nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("If-None-Match", `W/"wyzzy"`)
		resp, err := client.Do(req)
		if resp.StatusCode != 200 {
			t.Errorf("assert 'http.Response.StatusCode':: expected '%v', got '%v'", 200, resp.StatusCode)
		}
		body, err := ioutil.ReadAll(resp.Body)
		expectedBody := strings.Join([]string{
			"Redis:",
			"    Host: kubernate-redis:26379",
			"    Password: 1234",
			"    DB: 3",
			"    PoolSize: 128",
			"From: SettingResource"}, "\n")
		if string(body) != expectedBody {
			t.Errorf("assert 'http.Response.Body':: expected '%v', got '%v'", expectedBody, string(body))
		}
	}
	{
		req, err := http.NewRequest("GET", "http://127.0.0.1:10094/unknown", nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Add("If-None-Match", `W/"wyzzy"`)
		resp, err := client.Do(req)
		if resp.StatusCode != 404 {
			t.Errorf("assert 'http.Response.StatusCode':: expected '%v', got '%v'", 404, resp.StatusCode)
		}
	}
	{
		req, err := http.NewRequest("OCCUR", "http://127.0.0.1:10094/Accident", nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Add("If-None-Match", `W/"wyzzy"`)
		_, _ = client.Do(req)
		expectedErrorString := "err: an error occurred"
		if errorBuffer.String() != expectedErrorString {
			t.Errorf("assert 'errorBuffer':: expected '%v', got '%v'", expectedErrorString, errorBuffer.String())
		}
		errorBuffer.Reset()
	}
	{
		req, err := http.NewRequest("FAIL", "http://127.0.0.1:10094/Accident", nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Add("If-None-Match", `W/"wyzzy"`)
		_, _ = client.Do(req)
		expectedErrorString := "err: FAIL"
		if errorBuffer.String() != expectedErrorString {
			t.Errorf("assert 'errorBuffer':: expected '%v', got '%v'", expectedErrorString, errorBuffer.String())
		}
		errorBuffer.Reset()
	}
	{
		req, err := http.NewRequest("PING", "http://127.0.0.1:10094/Json", nil)
		if err != nil {
			t.Error(err)
		}
		_, _ = client.Do(req)
		resp, err := client.Do(req)
		if resp.StatusCode != 200 {
			t.Errorf("assert 'http.Response.StatusCode':: expected '%v', got '%v'", 200, resp.StatusCode)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if string(body) != `{"message":"OK"}` {
			t.Errorf("assert 'http.Response.Body':: expected '%v', got '%v'", "ECHO: Hello", string(body))
		}
	}
	{
		req, err := http.NewRequest("FAIL", "http://127.0.0.1:10094/Json", nil)
		if err != nil {
			t.Error(err)
		}
		_, _ = client.Do(req)
		resp, err := client.Do(req)
		if resp.StatusCode != 400 {
			t.Errorf("assert 'http.Response.StatusCode':: expected '%v', got '%v'", 200, resp.StatusCode)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if string(body) != `{"message":"UNKNOWN_ERROR"}` {
			t.Errorf("assert 'http.Response.Body':: expected '%v', got '%v'", "ECHO: Hello", string(body))
		}
	}
	{
		req, err := http.NewRequest("PING", "http://127.0.0.1:10094/Text", nil)
		if err != nil {
			t.Error(err)
		}
		_, _ = client.Do(req)
		resp, err := client.Do(req)
		if resp.StatusCode != 200 {
			t.Errorf("assert 'http.Response.StatusCode':: expected '%v', got '%v'", 200, resp.StatusCode)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if string(body) != "OK" {
			t.Errorf("assert 'http.Response.Body':: expected '%v', got '%v'", "ECHO: Hello", string(body))
		}
	}
	{
		req, err := http.NewRequest("FAIL", "http://127.0.0.1:10094/Text", nil)
		if err != nil {
			t.Error(err)
		}
		_, _ = client.Do(req)
		resp, err := client.Do(req)
		if resp.StatusCode != 400 {
			t.Errorf("assert 'http.Response.StatusCode':: expected '%v', got '%v'", 200, resp.StatusCode)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if string(body) != "UNKNOWN_ERROR" {
			t.Errorf("assert 'http.Response.Body':: expected '%v', got '%v'", "UNKNOWN_ERROR", string(body))
		}
	}
	{
		req, err := http.NewRequest("PING", "http://127.0.0.1:10094/Tracing", nil)
		if err != nil {
			t.Error(err)
		}
		// just test tracing, without check response
		client.Do(req)
	}

	select {
	case <-runCtx.Done():
		if err := starter.Stop(context.Background()); err != nil {
			t.Error(err)
		}
	}

	var expectedErrorCount = 2
	if errorCount != expectedErrorCount {
		t.Errorf("assert 'errorCount':: expected '%v', got '%v'", expectedErrorCount, errorCount)
	}

	// assert app.Host
	{
		if app.Host == nil {
			t.Error("assert 'MockApp.Host':: should not be nil")
		}
		host := app.Host
		var expectedListenAddress = ":10094"
		if host.ListenAddress != expectedListenAddress {
			t.Errorf("assert 'Host.ListenAddress':: expected '%v', got '%v'", expectedListenAddress, host.ListenAddress)
		}
		var expectedEnableCompress = true
		if host.EnableCompress != true {
			t.Errorf("assert 'Host.EnableCompress':: expected '%v', got '%v'", expectedEnableCompress, host.EnableCompress)
		}
		var expectedServerName = "DemoService"
		if host.Server.Name != expectedServerName {
			t.Errorf("assert 'Host.EnableCompress':: expected '%v', got '%v'", expectedServerName, host.Server.Name)
		}
	}
	// assert app.Config
	{
		if app.Config == nil {
			t.Error("assert 'MockApp.Config':: should not be nil")
		}
		conf := app.Config
		if conf.ListenAddress != ":10094" {
			t.Errorf("assert 'Config.ListenAddress':: expected '%v', got '%v'", ":10094", conf.ListenAddress)
		}
		if conf.EnableCompress != true {
			t.Errorf("assert 'Config.EnableCompress':: expected '%v', got '%v'", true, conf.EnableCompress)
		}
		if conf.ServerName != "DemoService" {
			t.Errorf("assert 'Config.ServerName':: expected '%v', got '%v'", "DemoService", conf.ServerName)
		}
		if conf.RedisHost != "kubernate-redis:26379" {
			t.Errorf("assert 'Config.RedisHost':: expected '%v', got '%v'", "kubernate-redis:26379", conf.RedisHost)
		}
		if conf.RedisPassword != "1234" {
			t.Errorf("assert 'Config.RedisPassword':: expected '%v', got '%v'", "1234", conf.RedisPassword)
		}
		if conf.RedisDB != 3 {
			t.Errorf("assert 'Config.RedisDB':: expected '%v', got '%v'", 3, conf.RedisDB)
		}
		if conf.RedisPoolSize != 128 {
			t.Errorf("assert 'Config.RedisPoolSize':: expected '%v', got '%v'", 128, conf.RedisPoolSize)
		}
		if conf.Workspace != "demo_test" {
			t.Errorf("assert 'Config.Workspace':: expected '%v', got '%v'", "demo_test", conf.Workspace)
		}
	}
	// assert app.ServiceProvider
	{
		if app.ServiceProvider == nil {
			t.Error("assert 'MockApp.ServiceProvider':: should not be nil")
		}
		provider := app.ServiceProvider
		if provider.CacheClient == nil {
			t.Error("assert 'ServiceProvider.RedisClient':: should not be nil")
		}
		redisClient := provider.CacheClient
		if redisClient.Host != "kubernate-redis:26379" {
			t.Errorf("assert 'RedisClient.Host':: expected '%v', got '%v'", "kubernate-redis:26379", redisClient.Host)
		}
		if redisClient.Password != "1234" {
			t.Errorf("assert 'RedisClient.Password':: expected '%v', got '%v'", "1234", redisClient.Password)
		}
		if redisClient.DB != 3 {
			t.Errorf("assert 'RedisClient.DB':: expected '%v', got '%v'", 3, redisClient.DB)
		}
		if redisClient.PoolSize != 128 {
			t.Errorf("assert 'RedisClient.PoolSize':: expected '%v', got '%v'", 128, redisClient.PoolSize)
		}
	}
}

func initializeEnvironment() {
	os.Setenv("REDIS_HOST", "kubernate-redis:26379")
	os.Setenv("REDIS_PASSWORD", "1234")
	os.Setenv("REDIS_POOL_SIZE", "128")
}

func initializeArgs() {
	os.Args = []string{"example",
		"--address", ":10094",
		"--compress", "true",
		"--hostname", "DemoService"}

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}
