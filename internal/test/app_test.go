package test

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Bofry/config"
	fasthttp "github.com/Bofry/host-fasthttp"
	"github.com/Bofry/host-fasthttp/response"
	"github.com/Bofry/host-fasthttp/response/failure"
)

var (
	__CONFIG_YAML_FILE        = "config.yaml"
	__CONFIG_YAML_FILE_SAMPLE = "config.yaml.sample"
)

type RequestManager struct {
	*RootRequest          `url:"/"`
	*EchoRequest          `url:"/Echo"`
	*SettingRequest       `url:"/Setting"`
	*AccidentRequest      `url:"/Accident"`
	*JsonRequest          `url:"/Json"`
	*TextRequest          `url:"/Text"`
	*TracingRequest       `url:"/Tracing"`
	*HttpProxyRequest     `url:"/Proxy/http"`
	*FasthttpProxyRequest `url:"/Proxy/fasthttp"`
}

func copyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}

func TestMain(m *testing.M) {
	_, err := os.Stat(__CONFIG_YAML_FILE)
	if err != nil {
		if os.IsNotExist(err) {
			err = copyFile(__CONFIG_YAML_FILE_SAMPLE, __CONFIG_YAML_FILE)
			if err != nil {
				panic(err)
			}
		}
	}
	m.Run()
}

func TestStartup(t *testing.T) {
	/* like
	 * $ export REDIS_HOST=kubernate-redis:26379
	 * $ export REDIS_PASSWORD=1234
	 * $ export REDIS_POOL_SIZE=128
	 */
	t.Setenv("REDIS_HOST", "kubernate-redis:26379")
	t.Setenv("REDIS_PASSWORD", "1234")
	t.Setenv("REDIS_POOL_SIZE", "128")

	/* like
	 * $ go run app.go --address ":10094" --compress true --hostname "DemoService"
	 */
	os.Args = []string{"example",
		"--address", ":10094",
		"--compress", "true",
		"--hostname", "DemoService"}

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	var (
		errorBuffer bytes.Buffer
	)

	app := App{}
	starter := fasthttp.Startup(&app).
		Middlewares(
			fasthttp.UseRequestManager(&RequestManager{}),
			fasthttp.UseXHttpMethodHeader(),
			fasthttp.UseErrorHandler(func(ctx *fasthttp.RequestCtx, err interface{}) {
				if fail, ok := err.(*failure.Failure); ok {
					if fail != nil {
						response.Json.Failure(ctx, fail, fasthttp.StatusBadRequest)
					}
				}
				if v, ok := err.(error); ok && v.Error() == "FAIL" {
					response.Json.Failure(ctx, "FAIL", 400)
				}
				fmt.Fprintf(&errorBuffer, "err: %+v", err)
			}),
			fasthttp.UseTracing(false),
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

	runCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
		req, err := http.NewRequest("FAIL2", "http://127.0.0.1:10094/Accident", nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Add("If-None-Match", `W/"wyzzy"`)
		_, _ = client.Do(req)
		expectedErrorString := "err: UNKNOWN_ERROR - an error occurred"
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
	{
		req, err := http.NewRequest("POST", "http://127.0.0.1:10094/Proxy/http", nil)
		if err != nil {
			t.Error(err)
		}
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
		req, err := http.NewRequest("POST", "http://127.0.0.1:10094/Proxy/fasthttp", nil)
		if err != nil {
			t.Error(err)
		}
		resp, err := client.Do(req)
		if resp.StatusCode != 200 {
			t.Errorf("assert 'http.Response.StatusCode':: expected '%v', got '%v'", 200, resp.StatusCode)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if string(body) != "OK" {
			t.Errorf("assert 'http.Response.Body':: expected '%v', got '%v'", "ECHO: Hello", string(body))
		}
	}

	select {
	case <-runCtx.Done():
		if err := starter.Stop(context.Background()); err != nil {
			t.Error(err)
		}
	}
}

func TestStartup_UseTracing(t *testing.T) {
	/* like
	 * $ export REDIS_HOST=kubernate-redis:26379
	 * $ export REDIS_PASSWORD=1234
	 * $ export REDIS_POOL_SIZE=128
	 */
	t.Setenv("REDIS_HOST", "kubernate-redis:26379")
	t.Setenv("REDIS_PASSWORD", "1234")
	t.Setenv("REDIS_POOL_SIZE", "128")

	/* like
	 * $ go run app.go --address ":10094" --compress true --hostname "DemoService"
	 */
	os.Args = []string{"example",
		"--address", ":10094",
		"--compress", "true",
		"--hostname", "DemoService"}

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	var (
		errorBuffer  bytes.Buffer
		testStartAt  time.Time
		requestCount int = 0
	)

	app := App{}
	starter := fasthttp.Startup(&app).
		Middlewares(
			fasthttp.UseRequestManager(&RequestManager{}),
			fasthttp.UseXHttpMethodHeader(),
			fasthttp.UseErrorHandler(func(ctx *fasthttp.RequestCtx, err interface{}) {
				if fail, ok := err.(*failure.Failure); ok {
					if fail != nil {
						response.Json.Failure(ctx, fail, fasthttp.StatusBadRequest)
					}
				}
				if v, ok := err.(error); ok && v.Error() == "FAIL" {
					response.Json.Failure(ctx, "FAIL", 400)
				}
				fmt.Fprintf(&errorBuffer, "err: %+v", err)
			}),
			fasthttp.UseTracing(true),
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

	runCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := starter.Start(runCtx); err != nil {
		t.Error(err)
	}

	testStartAt = time.Now()

	client := &http.Client{}
	{
		req, err := http.NewRequest("POST", "http://127.0.0.1:10094", nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("X-Http-Method", "PING")
		client.Do(req)
		requestCount++
	}
	{
		req, err := http.NewRequest("PING", "http://127.0.0.1:10094", nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("If-None-Match", `W/"wyzzy"`)
		client.Do(req)
		requestCount++
	}
	{
		req, err := http.NewRequest("SEND", "http://127.0.0.1:10094/Echo?input=Can%20you%20hear%20me", nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("If-None-Match", `W/"wyzzy"`)
		client.Do(req)
		requestCount++
	}
	{
		req, err := http.NewRequest("SEND", "http://127.0.0.1:10094/Echo/Hello", nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("If-None-Match", `W/"wyzzy"`)
		client.Do(req)
		requestCount++
	}
	{
		req, err := http.NewRequest("PEEK", "http://127.0.0.1:10094/Setting", nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("If-None-Match", `W/"wyzzy"`)
		client.Do(req)
		requestCount++
	}
	{
		req, err := http.NewRequest("GET", "http://127.0.0.1:10094/unknown", nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Add("If-None-Match", `W/"wyzzy"`)
		client.Do(req)
		requestCount++
	}
	{
		req, err := http.NewRequest("OCCUR", "http://127.0.0.1:10094/Accident", nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Add("If-None-Match", `W/"wyzzy"`)
		client.Do(req)
		requestCount++
	}
	{
		req, err := http.NewRequest("FAIL", "http://127.0.0.1:10094/Accident", nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Add("If-None-Match", `W/"wyzzy"`)
		client.Do(req)
		requestCount++
	}
	{
		req, err := http.NewRequest("FAIL2", "http://127.0.0.1:10094/Accident", nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Add("If-None-Match", `W/"wyzzy"`)
		client.Do(req)
		requestCount++
	}
	{
		req, err := http.NewRequest("PING", "http://127.0.0.1:10094/Json", nil)
		if err != nil {
			t.Error(err)
		}
		client.Do(req)
		requestCount++
	}
	{
		req, err := http.NewRequest("FAIL", "http://127.0.0.1:10094/Json", nil)
		if err != nil {
			t.Error(err)
		}
		client.Do(req)
		requestCount++
	}
	{
		req, err := http.NewRequest("PING", "http://127.0.0.1:10094/Text", nil)
		if err != nil {
			t.Error(err)
		}
		client.Do(req)
		requestCount++
	}
	{
		req, err := http.NewRequest("FAIL", "http://127.0.0.1:10094/Text", nil)
		if err != nil {
			t.Error(err)
		}
		client.Do(req)
		requestCount++
	}
	{
		req, err := http.NewRequest("PING", "http://127.0.0.1:10094/Tracing", nil)
		if err != nil {
			t.Error(err)
		}
		// just test tracing, without check response
		client.Do(req)
		requestCount++
	}
	{
		req, err := http.NewRequest("POST", "http://127.0.0.1:10094/Proxy/http", nil)
		if err != nil {
			t.Error(err)
		}
		client.Do(req)
		requestCount++
	}
	{
		req, err := http.NewRequest("POST", "http://127.0.0.1:10094/Proxy/fasthttp", nil)
		if err != nil {
			t.Error(err)
		}
		client.Do(req)
		requestCount++
	}

	select {
	case <-runCtx.Done():
		if err := starter.Stop(context.Background()); err != nil {
			t.Error(err)
		}

		testEndAt := time.Now()
		var queryUrl = fmt.Sprintf(
			"%s?end=%d&limit=50&lookback=1h&&service=fasthttp-trace-demo&start=%d",
			app.Config.JaegerQueryUrl,
			testEndAt.UnixMicro(),
			testStartAt.UnixMicro())
		req, err := http.NewRequest("GET", queryUrl, nil)
		if err != nil {
			t.Error(err)
		}
		resp, err := client.Do(req)
		if resp.StatusCode != 200 {
			t.Errorf("assert query 'Jeager Query Url StatusCode':: expected '%v', got '%v'", 200, resp.StatusCode)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
		}
		// t.Logf("%v", string(body))
		// parse content
		{
			var reply map[string]interface{}
			dec := json.NewDecoder(bytes.NewBuffer(body))
			dec.UseNumber()
			if err := dec.Decode(&reply); err != nil {
				t.Error(err)
			}

			data := reply["data"].([]interface{})
			if data == nil {
				t.Errorf("missing data section")
			}
			var expectedDataLength int = requestCount
			if expectedDataLength != len(data) {
				t.Errorf("assert 'Jaeger Query size of replies':: expected '%v', got '%v'", expectedDataLength, len(data))
			}
		}
	}
}

func TestStartup_UseLogging_And_UseTracing(t *testing.T) {
	/* like
	 * $ export REDIS_HOST=kubernate-redis:26379
	 * $ export REDIS_PASSWORD=1234
	 * $ export REDIS_POOL_SIZE=128
	 */
	t.Setenv("REDIS_HOST", "kubernate-redis:26379")
	t.Setenv("REDIS_PASSWORD", "1234")
	t.Setenv("REDIS_POOL_SIZE", "128")

	/* like
	 * $ go run app.go --address ":10094" --compress true --hostname "DemoService"
	 */
	os.Args = []string{"example",
		"--address", ":10094",
		"--compress", "true",
		"--hostname", "DemoService"}

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	var (
		loggingBuffer bytes.Buffer
		errorBuffer   bytes.Buffer
		testStartAt   time.Time
	)

	app := App{}
	starter := fasthttp.Startup(&app).
		Middlewares(
			fasthttp.UseRequestManager(&RequestManager{}),
			fasthttp.UseXHttpMethodHeader(),
			fasthttp.UseErrorHandler(func(ctx *fasthttp.RequestCtx, err interface{}) {
				if fail, ok := err.(*failure.Failure); ok {
					if fail != nil {
						response.Json.Failure(ctx, fail, fasthttp.StatusBadRequest)
					}
				}
				if v, ok := err.(error); ok && v.Error() == "FAIL" {
					response.Json.Failure(ctx, "FAIL", 400)
				}
				fmt.Fprintf(&errorBuffer, "err: %+v", err)
			}),
			fasthttp.UseLogging(&MultiLoggerService{
				LoggingServices: []fasthttp.LoggingService{
					&LoggingService{},
					&BlackholeLoggerService{
						Buffer: &loggingBuffer,
					},
				},
			}),
			fasthttp.UseTracing(true),
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

	runCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := starter.Start(runCtx); err != nil {
		t.Error(err)
	}

	testStartAt = time.Now()

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

	select {
	case <-runCtx.Done():
		if err := starter.Stop(context.Background()); err != nil {
			t.Error(err)
		}

		testEndAt := time.Now()
		var queryUrl = fmt.Sprintf(
			"%s?end=%d&limit=21&lookback=1h&&service=fasthttp-trace-demo&start=%d",
			app.Config.JaegerQueryUrl,
			testEndAt.UnixMicro(),
			testStartAt.UnixMicro())
		req, err := http.NewRequest("GET", queryUrl, nil)
		if err != nil {
			t.Error(err)
		}
		resp, err := client.Do(req)
		if resp.StatusCode != 200 {
			t.Errorf("assert query 'Jeager Query Url StatusCode':: expected '%v', got '%v'", 200, resp.StatusCode)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
		}
		// t.Logf("%v", string(body))
		// parse content
		{
			var reply map[string]interface{}
			dec := json.NewDecoder(bytes.NewBuffer(body))
			dec.UseNumber()
			if err := dec.Decode(&reply); err != nil {
				t.Error(err)
			}

			data := reply["data"].([]interface{})
			if data == nil {
				t.Errorf("missing data section")
			}
			var expectedDataLength int = 2
			if expectedDataLength != len(data) {
				t.Errorf("assert 'Jaeger Query size of replies':: expected '%v', got '%v'", expectedDataLength, len(data))
			}
		}

		// test loggingBuffer
		var expectedLoggingBuffer string = strings.Join([]string{
			"CreateEventLog()\n",
			"OnProcessRequest()\n",
			"OnProcessRequestComplete()\n",
			"Flush()\n",
			"CreateEventLog()\n",
			"OnProcessRequest()\n",
			"OnError()\n",
			"Flush()\n",
		}, "")
		if expectedLoggingBuffer != loggingBuffer.String() {
			t.Errorf("assert loggingBuffer:: expected '%v', got '%v'", expectedLoggingBuffer, loggingBuffer.String())
		}
	}
}
