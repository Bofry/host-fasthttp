package response

import (
	"testing"

	"github.com/Bofry/host-fasthttp/testingutil"
	"github.com/valyala/fasthttp"
)

func TestTextFormatterSuccess(t *testing.T) {
	t.Parallel()

	handler := func(ctx *fasthttp.RequestCtx) {
		Text.Success(ctx, "OK")

		// check
		{
			obj := ExtractResponseState(ctx)
			v, ok := obj.(ResponseState)
			if !ok {
				t.Errorf("assert 'ExtractResponseFlag(ctx)':: expected '%s', got '%T'", "Response", v)
			}
			if v.Flag() != SUCCESS {
				t.Errorf("assert 'Response.Flag()':: expected '%v', got '%v'", SUCCESS, v.Flag())
			}
			if v.StatusCode() != fasthttp.StatusOK {
				t.Errorf("assert 'Response.StatusCode()':: expected '%v', got '%v'", fasthttp.StatusOK, v.StatusCode())
			}
		}
	}

	resp, err := testingutil.RunRequestHandler(handler,
		[]byte("GET / HTTP/1.1\r\nHost: g.com\r\n\r\n"))

	if err != nil {
		t.Fatal(err)
	}

	// t.Logf("result: %v", resp.StatusCode())
	// t.Logf("result: %v", resp.Header.String())
	// t.Logf("result: %v", string(resp.Body()))

	var expectedStatusCode = 200
	if resp.StatusCode() != expectedStatusCode {
		t.Errorf("status code: expect %v, got %v", expectedStatusCode, resp.StatusCode())
	}
	var (
		expectedContentType = "text/plain; charset=utf-8"
		actualContentType   = string(resp.Header.Peek("Content-Type"))
	)
	if actualContentType != expectedContentType {
		t.Errorf("content-type: expect %v, got %v", expectedContentType, actualContentType)
	}
	var (
		expectedBody = "OK"
		actualBody   = string(resp.Body())
	)
	if actualBody != expectedBody {
		t.Errorf("body: expect %v, got %v", expectedBody, actualBody)
	}
}

func TestTextFormatterFailure(t *testing.T) {
	t.Parallel()

	handler := func(ctx *fasthttp.RequestCtx) {
		Text.Failure(ctx, "UNKNOWN_ERROR", fasthttp.StatusBadRequest)

		// check
		{
			obj := ExtractResponseState(ctx)
			v, ok := obj.(ResponseState)
			if !ok {
				t.Errorf("assert 'ExtractResponseFlag(ctx)':: expected '%s', got '%T'", "Response", v)
			}
			if v.Flag() != FAILURE {
				t.Errorf("assert 'Response.Flag()':: expected '%v', got '%v'", SUCCESS, v.Flag())
			}
			if v.StatusCode() != fasthttp.StatusBadRequest {
				t.Errorf("assert 'Response.StatusCode()':: expected '%v', got '%v'", fasthttp.StatusOK, v.StatusCode())
			}
		}
	}

	resp, err := testingutil.RunRequestHandler(handler,
		[]byte("GET / HTTP/1.1\r\nHost: g.com\r\n\r\n"))

	if err != nil {
		t.Fatal(err)
	}

	// t.Logf("result: %v", resp.StatusCode())
	// t.Logf("result: %v", resp.Header.String())
	// t.Logf("result: %v", string(resp.Body()))

	var expectedStatusCode = 400
	if resp.StatusCode() != expectedStatusCode {
		t.Errorf("status code: expect %v, got %v", expectedStatusCode, resp.StatusCode())
	}
	var (
		expectedContentType = "text/plain; charset=utf-8"
		actualContentType   = string(resp.Header.Peek("Content-Type"))
	)
	if actualContentType != expectedContentType {
		t.Errorf("content-type: expect %v, got %v", expectedContentType, actualContentType)
	}
	var (
		expectedBody = "UNKNOWN_ERROR"
		actualBody   = string(resp.Body())
	)
	if actualBody != expectedBody {
		t.Errorf("body: expect %v, got %v", expectedBody, actualBody)
	}
}
