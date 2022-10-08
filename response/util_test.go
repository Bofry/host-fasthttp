package response

import (
	"bufio"
	"io"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

func useInmemoryServer(
	s *fasthttp.Server,
	requestHandler func(w io.Writer),
	responseHandler func(resp *fasthttp.Response)) error {

	ln := fasthttputil.NewInmemoryListener()
	defer ln.Close()

	go func() {
		if err := s.Serve(ln); err != nil {
			panic(err)
		}
	}()

	c, err := ln.Dial()
	if err != nil {
		return err
	}
	defer c.Close()

	requestHandler(c)

	br := bufio.NewReader(c)
	var resp fasthttp.Response
	if err := resp.Read(br); err != nil {
		return err
	}

	responseHandler(&resp)
	return nil
}
