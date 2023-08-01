package websocket

import (
	"time"

	"github.com/fasthttp/websocket"
	"github.com/valyala/fasthttp"
)

var _ MessageClientOption = MessageClientOptionFunc(nil)

type MessageClientOptionFunc func(websocket.FastHTTPUpgrader)

func (f MessageClientOptionFunc) apply(u websocket.FastHTTPUpgrader) {
	f(u)
}

// -----------------------------------------------

func WithHandshakeTimeout(timeout time.Duration) MessageClientOption {
	return MessageClientOptionFunc(func(u websocket.FastHTTPUpgrader) {
		u.HandshakeTimeout = timeout
	})
}

func WithReadBufferSize(size int) MessageClientOption {
	return MessageClientOptionFunc(func(u websocket.FastHTTPUpgrader) {
		u.ReadBufferSize = size
	})
}

func WithWriteBufferSize(size int) MessageClientOption {
	return MessageClientOptionFunc(func(u websocket.FastHTTPUpgrader) {
		u.WriteBufferSize = size
	})
}

func WithWriteBufferPool(pool websocket.BufferPool) MessageClientOption {
	return MessageClientOptionFunc(func(u websocket.FastHTTPUpgrader) {
		u.WriteBufferPool = pool
	})
}

func WithSubprotocols(v []string) MessageClientOption {
	return MessageClientOptionFunc(func(u websocket.FastHTTPUpgrader) {
		u.Subprotocols = v
	})
}

func WithError(action func(ctx *fasthttp.RequestCtx, status int, reason error)) MessageClientOption {
	return MessageClientOptionFunc(func(u websocket.FastHTTPUpgrader) {
		u.Error = action
	})
}

func WithCheckOrigin(predicate func(ctx *fasthttp.RequestCtx) bool) MessageClientOption {
	return MessageClientOptionFunc(func(u websocket.FastHTTPUpgrader) {
		u.CheckOrigin = predicate
	})
}

func WithEnableCompression(enabled bool) MessageClientOption {
	return MessageClientOptionFunc(func(u websocket.FastHTTPUpgrader) {
		u.EnableCompression = enabled
	})
}
