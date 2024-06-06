package testingutil

import (
	"bufio"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

type InmemoryRequestHandler fasthttp.RequestHandler

func (h InmemoryRequestHandler) Process(payload []byte) (*fasthttp.Response, error) {
	ln := fasthttputil.NewInmemoryListener()
	defer ln.Close()

	s := &fasthttp.Server{
		Handler: fasthttp.RequestHandler(h),
	}

	go func() {
		if err := s.Serve(ln); err != nil {
			panic(err)
		}
	}()

	c, err := ln.Dial()
	if err != nil {
		return nil, err
	}
	defer c.Close()

	_, err = c.Write(payload)
	if err != nil {
		return nil, err
	}

	var (
		resp fasthttp.Response
	)

	br := bufio.NewReader(c)
	if err := resp.Read(br); err != nil {
		return nil, err
	}

	return &resp, nil
}
