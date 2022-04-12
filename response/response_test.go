package response

import (
	"bufio"
	"context"
	"testing"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

func TestSuccess(t *testing.T) {
	t.Parallel()

	s := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			Success(ctx, "text/plain", []byte("success message"))

			// check
			{
				obj := ctx.UserValue(RESPONSE_INVARIANT_NAME)
				v, ok := obj.(Response)
				if !ok {
					t.Errorf("assert 'ctx.UserValue(RESPONSE_INVARIANT_NAME)':: expected '%s', got '%T'", "Response", v)
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

	ln := fasthttputil.NewInmemoryListener()

	go func() {
		if err := s.Serve(ln); err != nil {
			t.Errorf("unexpected error: %s", err)
		}
	}()

	c, err := ln.Dial()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	_, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if _, err = c.Write([]byte("GET / HTTP/1.1\r\nHost: gle.com\r\n\r\n")); err != nil {
		t.Fatal(err)
	}

	br := bufio.NewReader(c)
	var resp fasthttp.Response
	if err := resp.Read(br); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if err := c.Close(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if err := ln.Close(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestError(t *testing.T) {
	t.Parallel()

	s := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			Failure(ctx, "text/plain", []byte("error message"), fasthttp.StatusBadRequest)

			// check
			{
				obj := ctx.UserValue(RESPONSE_INVARIANT_NAME)
				v, ok := obj.(Response)
				if !ok {
					t.Errorf("assert 'ctx.UserValue(RESPONSE_INVARIANT_NAME)':: expected '%s', got '%T'", "Response", v)
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

	ln := fasthttputil.NewInmemoryListener()

	go func() {
		if err := s.Serve(ln); err != nil {
			t.Errorf("unexpected error: %s", err)
		}
	}()

	c, err := ln.Dial()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	_, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if _, err = c.Write([]byte("GET / HTTP/1.1\r\nHost: gle.com\r\n\r\n")); err != nil {
		t.Fatal(err)
	}

	br := bufio.NewReader(c)
	var resp fasthttp.Response
	if err := resp.Read(br); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if err := c.Close(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if err := ln.Close(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}
