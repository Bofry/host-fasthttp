package response

import (
	"io"
	"testing"

	"github.com/valyala/fasthttp"
)

func TestSuccess(t *testing.T) {
	t.Parallel()

	s := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			Success(ctx, "text/plain", []byte("success message"))

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

			var expectedStatusCode = 200
			if resp.StatusCode() != expectedStatusCode {
				t.Errorf("status code: expect %v, got %v", expectedStatusCode, resp.StatusCode())
			}
			var (
				expectedContentType = "text/plain"
				actualContentType   = string(resp.Header.Peek("Content-Type"))
			)
			if actualContentType != expectedContentType {
				t.Errorf("content-type: expect %v, got %v", expectedContentType, actualContentType)
			}
			var (
				expectedBody = "success message"
				actualBody   = string(resp.Body())
			)
			if actualBody != expectedBody {
				t.Errorf("body: expect %v, got %v", expectedBody, actualBody)
			}
		})
}

func TestFailure(t *testing.T) {
	t.Parallel()

	s := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			Failure(ctx, "text/plain", []byte("error message"), fasthttp.StatusBadRequest)

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
		func(w io.Writer) {
			if _, err := w.Write([]byte("GET / HTTP/1.1\r\nHost: g.com\r\n\r\n")); err != nil {
				t.Fatal(err)
			}
		},
		func(resp *fasthttp.Response) {
			// t.Logf("result: %v", resp.StatusCode())
			// t.Logf("result: %v", resp.Header.String())
			// t.Logf("result: %v", string(resp.Body()))

			var expectedStatusCode = 400
			if resp.StatusCode() != expectedStatusCode {
				t.Errorf("status code: expect %v, got %v", expectedStatusCode, resp.StatusCode())
			}
			var (
				expectedContentType = "text/plain"
				actualContentType   = string(resp.Header.Peek("Content-Type"))
			)
			if actualContentType != expectedContentType {
				t.Errorf("content-type: expect %v, got %v", expectedContentType, actualContentType)
			}
			var (
				expectedBody = "error message"
				actualBody   = string(resp.Body())
			)
			if actualBody != expectedBody {
				t.Errorf("body: expect %v, got %v", expectedBody, actualBody)
			}
		})
}
