package websocket

import "github.com/fasthttp/websocket"

type (
	MessageClientOption interface {
		apply(*websocket.FastHTTPUpgrader)
	}
)
