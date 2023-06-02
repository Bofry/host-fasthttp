package http

import (
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

const defaultMaxRedirectsCount = 16

var (
	strPostArgsContentType = []byte("application/x-www-form-urlencoded")

	clientURLResponseChPool sync.Pool
)

type clientURLResponse struct {
	statusCode int
	body       []byte
	err        error
}

// An extenstion method for github.com/valyala/fasthttp.Do()
func Do(req *fasthttp.Request, resp *fasthttp.Response, opts ...HttpClientOption) error {
	for _, opt := range opts {
		err := opt.apply(req, resp)
		if err != nil {
			return err
		}
	}
	return fasthttp.Do(req, resp)
}

// An extenstion method for github.com/valyala/fasthttp.DoTimeout()
func DoTimeout(req *fasthttp.Request, resp *fasthttp.Response, timeout time.Duration, opts ...HttpClientOption) error {
	for _, opt := range opts {
		err := opt.apply(req, resp)
		if err != nil {
			return err
		}
	}
	return fasthttp.DoTimeout(req, resp, timeout)
}

// An extenstion method for github.com/valyala/fasthttp.DoDeadline()
func DoDeadline(req *fasthttp.Request, resp *fasthttp.Response, deadline time.Time, opts ...HttpClientOption) error {
	for _, opt := range opts {
		err := opt.apply(req, resp)
		if err != nil {
			return err
		}
	}
	return fasthttp.DoDeadline(req, resp, deadline)
}

// An extenstion method for github.com/valyala/fasthttp.DoRedirects()
func DoRedirects(req *fasthttp.Request, resp *fasthttp.Response, maxRedirectsCount int, opts ...HttpClientOption) error {
	for _, opt := range opts {
		err := opt.apply(req, resp)
		if err != nil {
			return err
		}
	}
	return fasthttp.DoRedirects(req, resp, maxRedirectsCount)
}

// An extenstion method for github.com/valyala/fasthttp.Get()
func Get(dst []byte, url string, opts ...HttpClientOption) (statusCode int, body []byte, err error) {
	req := AcquireRequest()
	resp := AcquireResponse()
	defer ReleaseRequest(req)
	defer ReleaseResponse(resp)

	req.Header.SetMethod(MethodGet)
	req.Header.SetRequestURI(url)

	err = DoRedirects(req, resp, defaultMaxRedirectsCount, opts...)
	if err != nil {
		return
	}
	statusCode = resp.StatusCode()
	body = resp.Body()
	return
}

// An extenstion method for github.com/valyala/fasthttp.GetTimeout()
func GetTimeout(dst []byte, url string, timeout time.Duration, opts ...HttpClientOption) (statusCode int, body []byte, err error) {
	deadline := time.Now().Add(timeout)
	return GetDeadline(dst, url, deadline, opts...)
}

// An extenstion method for github.com/valyala/fasthttp.GetDeadline()
func GetDeadline(dst []byte, url string, deadline time.Time, opts ...HttpClientOption) (statusCode int, body []byte, err error) {
	timeout := -time.Since(deadline)
	if timeout <= 0 {
		return 0, dst, fasthttp.ErrTimeout
	}

	var ch chan clientURLResponse
	chv := clientURLResponseChPool.Get()
	if chv == nil {
		chv = make(chan clientURLResponse, 1)
	}
	ch = chv.(chan clientURLResponse)

	var mu sync.Mutex
	var timedout, responded bool

	go func() {
		req := AcquireRequest()
		resp := AcquireResponse()
		defer ReleaseRequest(req)
		defer ReleaseResponse(resp)

		req.Header.SetMethod(fasthttp.MethodGet)
		req.Header.SetRequestURI(url)

		errCopy := DoRedirects(req, resp, defaultMaxRedirectsCount, opts...)
		if err != nil {
			return
		}
		statusCodeCopy := resp.StatusCode()
		bodyCopy := resp.Body()
		mu.Lock()
		{
			if !timedout {
				ch <- clientURLResponse{
					statusCode: statusCodeCopy,
					body:       bodyCopy,
					err:        errCopy,
				}
				responded = true
			}
		}
		mu.Unlock()
	}()

	tc := fasthttp.AcquireTimer(timeout)
	select {
	case resp := <-ch:
		statusCode = resp.statusCode
		body = resp.body
		err = resp.err
	case <-tc.C:
		mu.Lock()
		{
			if responded {
				resp := <-ch
				statusCode = resp.statusCode
				body = resp.body
				err = resp.err
			} else {
				timedout = true
				err = fasthttp.ErrTimeout
				body = dst
			}
		}
		mu.Unlock()
	}
	fasthttp.ReleaseTimer(tc)

	clientURLResponseChPool.Put(chv)

	return statusCode, body, err
}

// An extenstion method for github.com/valyala/fasthttp.Post()
func Post(dst []byte, url string, postArgs *fasthttp.Args, opts ...HttpClientOption) (statusCode int, body []byte, err error) {
	req := AcquireRequest()
	resp := AcquireResponse()
	defer ReleaseRequest(req)
	defer ReleaseResponse(resp)

	req.Header.SetMethod(MethodPost)
	req.Header.SetRequestURI(url)
	req.Header.SetContentTypeBytes(strPostArgsContentType)
	if postArgs != nil {
		if _, err := postArgs.WriteTo(req.BodyWriter()); err != nil {
			return 0, nil, err
		}
	}

	err = DoRedirects(req, resp, defaultMaxRedirectsCount, opts...)
	if err != nil {
		return
	}
	statusCode = resp.StatusCode()
	body = resp.Body()
	return
}
