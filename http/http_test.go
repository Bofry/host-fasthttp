package http

import (
	"context"
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"
	"time"

	fasthttp "github.com/Bofry/host-fasthttp"
	"github.com/joho/godotenv"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

var (
	traceIDStr = "4bf92f3577b34da6a3ce929d0e0e4736"
	spanIDStr  = "00f067aa0ba902b7"

	__ENV_FILE        = "http_test.env"
	__ENV_FILE_SAMPLE = "http_test.env.sample"

	__TEST_HTTP_SERVER_ADDR string
	__TEST_HTTP_SERVER_URL  string

	__TEST_TRACE_ID = mustTraceIDFromHex(traceIDStr)
	__TEST_SPAN_ID  = mustSpanIDFromHex(spanIDStr)

	__TEST_PROPAGATOR = propagation.TraceContext{}
	__TEST_CONTEXT    = mustSpanContext()
)

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

func mustTraceIDFromHex(s string) (t trace.TraceID) {
	var err error
	t, err = trace.TraceIDFromHex(s)
	if err != nil {
		panic(err)
	}
	return
}

func mustSpanIDFromHex(s string) (t trace.SpanID) {
	var err error
	t, err = trace.SpanIDFromHex(s)
	if err != nil {
		panic(err)
	}
	return
}

func mustSpanContext() context.Context {
	sc := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    __TEST_TRACE_ID,
		SpanID:     __TEST_SPAN_ID,
		TraceFlags: 0,
	})
	return trace.ContextWithSpanContext(context.Background(), sc)
}

func TestMain(m *testing.M) {
	_, err := os.Stat(__ENV_FILE)
	if err != nil {
		if os.IsNotExist(err) {
			err = copyFile(__ENV_FILE_SAMPLE, __ENV_FILE)
			if err != nil {
				panic(err)
			}
		}
	}

	godotenv.Load(__ENV_FILE)
	{
		__TEST_HTTP_SERVER_ADDR = os.Getenv("HTTP_SERVER_ADDR")
		__TEST_HTTP_SERVER_URL = "http://" + __TEST_HTTP_SERVER_ADDR
	}
	m.Run()
}

func TestDo(t *testing.T) {
	s := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			traceparent := ctx.Request.Header.Peek("traceparent")
			if len(traceparent) == 0 {
				t.Error("missing request header 'traceparent'")
			}

			body := ctx.Request.Body()
			ctx.Write(body) //nolint:errcheck
		},
	}
	defer s.Shutdown()
	go s.ListenAndServe(__TEST_HTTP_SERVER_ADDR) //nolint:errcheck

	req, res := AcquireRequest(), AcquireResponse()
	defer func() {
		ReleaseRequest(req)
		ReleaseResponse(res)
	}()
	req.Header.SetMethod(MethodGet)
	req.SetRequestURI(__TEST_HTTP_SERVER_URL)
	req.SetBodyString("test")

	err := Do(req, res,
		WithTracePropagation(__TEST_CONTEXT, __TEST_PROPAGATOR),
	)
	if err != nil {
		t.Error(err)
	}
	if len(res.Body()) == 0 {
		t.Error("missing request body")
	}
	var expectedBody = []byte("test")
	if !reflect.DeepEqual(expectedBody, res.Body()) {
		t.Errorf("Response.Body() expect: %v, got: %v", expectedBody, res.Body())
	}
}

func TestDoTimeout(t *testing.T) {
	s := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			traceparent := ctx.Request.Header.Peek("traceparent")
			if len(traceparent) == 0 {
				t.Error("missing request header 'traceparent'")
			}

			body := ctx.Request.Body()
			ctx.Write(body) //nolint:errcheck
		},
	}
	defer s.Shutdown()
	go s.ListenAndServe(__TEST_HTTP_SERVER_ADDR) //nolint:errcheck

	req, res := AcquireRequest(), AcquireResponse()
	defer func() {
		ReleaseRequest(req)
		ReleaseResponse(res)
	}()
	req.Header.SetMethod(MethodGet)
	req.SetRequestURI(__TEST_HTTP_SERVER_URL)
	req.SetBodyString("test")

	err := DoTimeout(req, res, 2*time.Second,
		WithTracePropagation(__TEST_CONTEXT, __TEST_PROPAGATOR),
	)
	if err != nil {
		t.Error(err)
	}
	if len(res.Body()) == 0 {
		t.Error("missing request body")
	}
	var expectedBody = []byte("test")
	if !reflect.DeepEqual(expectedBody, res.Body()) {
		t.Errorf("Response.Body() expect: %v, got: %v", expectedBody, res.Body())
	}
}

func TestDoDeadline(t *testing.T) {
	s := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			traceparent := ctx.Request.Header.Peek("traceparent")
			if len(traceparent) == 0 {
				t.Error("missing request header 'traceparent'")
			}

			body := ctx.Request.Body()
			ctx.Write(body) //nolint:errcheck
		},
	}
	defer s.Shutdown()
	go s.ListenAndServe(__TEST_HTTP_SERVER_ADDR) //nolint:errcheck

	req, res := AcquireRequest(), AcquireResponse()
	defer func() {
		ReleaseRequest(req)
		ReleaseResponse(res)
	}()
	req.Header.SetMethod(MethodGet)
	req.SetRequestURI(__TEST_HTTP_SERVER_URL)
	req.SetBodyString("test")

	err := DoDeadline(req, res, time.Now().Add(2*time.Second),
		WithTracePropagation(__TEST_CONTEXT, __TEST_PROPAGATOR),
	)
	if err != nil {
		t.Error(err)
	}
	if len(res.Body()) == 0 {
		t.Error("missing request body")
	}
	var expectedBody = []byte("test")
	if !reflect.DeepEqual(expectedBody, res.Body()) {
		t.Errorf("Response.Body() expect: %v, got: %v", expectedBody, res.Body())
	}
}

func TestDoRedirects(t *testing.T) {
	s := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			traceparent := ctx.Request.Header.Peek("traceparent")
			if len(traceparent) == 0 {
				t.Error("missing request header 'traceparent'")
			}

			body := ctx.Request.Body()
			ctx.Write(body) //nolint:errcheck
		},
	}
	defer s.Shutdown()
	go s.ListenAndServe(__TEST_HTTP_SERVER_ADDR) //nolint:errcheck

	req, res := AcquireRequest(), AcquireResponse()
	defer func() {
		ReleaseRequest(req)
		ReleaseResponse(res)
	}()
	req.Header.SetMethod(MethodGet)
	req.SetRequestURI(__TEST_HTTP_SERVER_URL)
	req.SetBodyString("test")

	err := DoRedirects(req, res, 10,
		WithTracePropagation(__TEST_CONTEXT, __TEST_PROPAGATOR),
	)
	if err != nil {
		t.Error(err)
	}
	if len(res.Body()) == 0 {
		t.Error("missing request body")
	}
	var expectedBody = []byte("test")
	if !reflect.DeepEqual(expectedBody, res.Body()) {
		t.Errorf("Response.Body() expect: %v, got: %v", expectedBody, res.Body())
	}
}

func TestGet(t *testing.T) {
	s := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			traceparent := ctx.Request.Header.Peek("traceparent")
			if len(traceparent) == 0 {
				t.Error("missing request header 'traceparent'")
			}

			ctx.Write([]byte("test")) //nolint:errcheck
		},
	}
	defer s.Shutdown()
	go s.ListenAndServe(__TEST_HTTP_SERVER_ADDR) //nolint:errcheck

	req, res := AcquireRequest(), AcquireResponse()
	defer func() {
		ReleaseRequest(req)
		ReleaseResponse(res)
	}()

	code, body, err := Get(nil, __TEST_HTTP_SERVER_URL,
		WithTracePropagation(__TEST_CONTEXT, __TEST_PROPAGATOR),
	)
	if err != nil {
		t.Error(err)
	}
	var expectedCode = 200
	if expectedCode != code {
		t.Errorf("StatusCode expect: %v, got: %v", expectedCode, code)
	}
	if len(body) == 0 {
		t.Error("missing request body")
	}
	var expectedBody = []byte("test")
	if !reflect.DeepEqual(expectedBody, body) {
		t.Errorf("Response Body expect: %v, got: %v", expectedBody, body)
	}
}

func TestGetTimeout(t *testing.T) {
	s := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			traceparent := ctx.Request.Header.Peek("traceparent")
			if len(traceparent) == 0 {
				t.Error("missing request header 'traceparent'")
			}

			ctx.Write([]byte("test")) //nolint:errcheck
		},
	}
	defer s.Shutdown()
	go s.ListenAndServe(__TEST_HTTP_SERVER_ADDR) //nolint:errcheck

	req, res := AcquireRequest(), AcquireResponse()
	defer func() {
		ReleaseRequest(req)
		ReleaseResponse(res)
	}()

	code, body, err := GetTimeout(nil, __TEST_HTTP_SERVER_URL, 2*time.Second,
		WithTracePropagation(__TEST_CONTEXT, __TEST_PROPAGATOR),
	)
	if err != nil {
		t.Error(err)
	}
	var expectedCode = 200
	if expectedCode != code {
		t.Errorf("StatusCode expect: %v, got: %v", expectedCode, code)
	}
	if len(body) == 0 {
		t.Error("missing request body")
	}
	var expectedBody = []byte("test")
	if !reflect.DeepEqual(expectedBody, body) {
		t.Errorf("Response Body expect: %v, got: %v", expectedBody, body)
	}
}

func TestGetDeadline(t *testing.T) {
	s := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			traceparent := ctx.Request.Header.Peek("traceparent")
			if len(traceparent) == 0 {
				t.Error("missing request header 'traceparent'")
			}

			ctx.Write([]byte("test")) //nolint:errcheck
		},
	}
	defer s.Shutdown()
	go s.ListenAndServe(__TEST_HTTP_SERVER_ADDR) //nolint:errcheck

	req, res := AcquireRequest(), AcquireResponse()
	defer func() {
		ReleaseRequest(req)
		ReleaseResponse(res)
	}()

	code, body, err := GetDeadline(nil, __TEST_HTTP_SERVER_URL, time.Now().Add(2*time.Second),
		WithTracePropagation(__TEST_CONTEXT, __TEST_PROPAGATOR),
	)
	if err != nil {
		t.Error(err)
	}
	var expectedCode = 200
	if expectedCode != code {
		t.Errorf("StatusCode expect: %v, got: %v", expectedCode, code)
	}
	if len(body) == 0 {
		t.Error("missing request body")
	}
	var expectedBody = []byte("test")
	if !reflect.DeepEqual(expectedBody, body) {
		t.Errorf("Response Body expect: %v, got: %v", expectedBody, body)
	}
}

func TestPost(t *testing.T) {
	s := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			traceparent := ctx.Request.Header.Peek("traceparent")
			if len(traceparent) == 0 {
				t.Error("missing request header 'traceparent'")
			}

			args := ctx.PostArgs()
			args.Peek("foo")
			ctx.Write(args.Peek("foo")) //nolint:errcheck
		},
	}
	defer s.Shutdown()
	go s.ListenAndServe(__TEST_HTTP_SERVER_ADDR) //nolint:errcheck

	req, res := AcquireRequest(), AcquireResponse()
	defer func() {
		ReleaseRequest(req)
		ReleaseResponse(res)
	}()

	var args Args
	args.Add("foo", "test")

	code, body, err := Post(nil, __TEST_HTTP_SERVER_URL, &args,
		WithTracePropagation(__TEST_CONTEXT, __TEST_PROPAGATOR),
	)
	if err != nil {
		t.Error(err)
	}
	var expectedCode = 200
	if expectedCode != code {
		t.Errorf("StatusCode expect: %v, got: %v", expectedCode, code)
	}
	if len(body) == 0 {
		t.Error("missing request body")
	}
	var expectedBody = []byte("test")
	if !reflect.DeepEqual(expectedBody, body) {
		t.Errorf("Response Body expect: %v, got: %v", expectedBody, body)
	}
}
