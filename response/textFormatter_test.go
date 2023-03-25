package response

import (
	"io"
	"testing"

	"github.com/valyala/fasthttp"
)

func TestTextFormatterSuccess(t *testing.T) {
	t.Parallel()

	s := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
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
		},
	}

	useInmemoryServer(s,
		// request handler
		func(w io.Writer) {
			if _, err := w.Write([]byte("GET / HTTP/1.1\r\nHost: g.com\r\n\r\n")); err != nil {
				t.Fatal(err)
			}
		},
		// response handler
		func(resp *fasthttp.Response) {
			// t.Logf("result: %v", resp.StatusCode())
			// t.Logf("result: %v", resp.Header.String())
			// t.Logf("result: %v", string(resp.Body()))

			var exceptedStatusCode = 200
			if resp.StatusCode() != exceptedStatusCode {
				t.Errorf("status code: except %v, got %v", exceptedStatusCode, resp.StatusCode())
			}
			var (
				exceptedContentType = "text/plain; charset=utf-8"
				actualContentType   = string(resp.Header.Peek("Content-Type"))
			)
			if actualContentType != exceptedContentType {
				t.Errorf("content-type: except %v, got %v", exceptedContentType, actualContentType)
			}
			var (
				exceptedBody = "OK"
				actualBody   = string(resp.Body())
			)
			if actualBody != exceptedBody {
				t.Errorf("body: except %v, got %v", exceptedBody, actualBody)
			}
		})
}

func TestTextFormatterFailure(t *testing.T) {
	t.Parallel()

	s := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
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
		},
	}

	useInmemoryServer(s,
		// request handler
		func(w io.Writer) {
			if _, err := w.Write([]byte("GET / HTTP/1.1\r\nHost: g.com\r\n\r\n")); err != nil {
				t.Fatal(err)
			}
		},
		// response handler
		func(resp *fasthttp.Response) {
			// t.Logf("result: %v", resp.StatusCode())
			// t.Logf("result: %v", resp.Header.String())
			// t.Logf("result: %v", string(resp.Body()))

			var exceptedStatusCode = 400
			if resp.StatusCode() != exceptedStatusCode {
				t.Errorf("status code: except %v, got %v", exceptedStatusCode, resp.StatusCode())
			}
			var (
				exceptedContentType = "text/plain; charset=utf-8"
				actualContentType   = string(resp.Header.Peek("Content-Type"))
			)
			if actualContentType != exceptedContentType {
				t.Errorf("content-type: except %v, got %v", exceptedContentType, actualContentType)
			}
			var (
				exceptedBody = "UNKNOWN_ERROR"
				actualBody   = string(resp.Body())
			)
			if actualBody != exceptedBody {
				t.Errorf("body: except %v, got %v", exceptedBody, actualBody)
			}
		})
}
