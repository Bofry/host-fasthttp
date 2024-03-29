package websocket

import (
	"fmt"
	"sync"

	"github.com/Bofry/host/app"
	"github.com/fasthttp/websocket"
	"github.com/valyala/fasthttp"
)

var (
	__MessageTypeMap = map[int]app.MessageFormat{
		websocket.TextMessage:   app.TEXT_MESSAGE,
		websocket.BinaryMessage: app.BINARY_MESSAGE,
		websocket.CloseMessage:  app.CLOSE_MESSAGE,
		websocket.PingMessage:   app.PING_MESSAGE,
		websocket.PongMessage:   app.PONG_MESSAGE,
	}

	__MessageFormatMap = map[app.MessageFormat]int{
		app.TEXT_MESSAGE:   websocket.TextMessage,
		app.BINARY_MESSAGE: websocket.BinaryMessage,
	}
)

var _ app.MessageClient = new(MessageClient)

type MessageClient struct {
	ctx     *fasthttp.RequestCtx
	options []MessageClientOption

	onCloseDelegate []func(app.MessageClient)

	message chan *Message
	stop    chan struct{}
	done    chan struct{}

	stopped bool
	closed  bool
	mutex   sync.Mutex

	*app.MessageClientInfo
}

func NewMessageClient(ctx *fasthttp.RequestCtx, opts ...MessageClientOption) *MessageClient {
	return &MessageClient{
		ctx:     ctx,
		message: make(chan *Message),
		stop:    make(chan struct{}),
		done:    make(chan struct{}),
		options: opts,

		MessageClientInfo: app.NewMessageClientInfo(),
	}
}

// RegisterCloseHandler implements app.MessageClient.
func (client *MessageClient) RegisterCloseHandler(proc func(app.MessageClient)) {
	client.onCloseDelegate = append(client.onCloseDelegate, proc)
}

// Close implements app.MessageSource.
func (client *MessageClient) Close() error {
	if client.closed {
		return nil
	}

	client.mutex.Lock()
	defer client.mutex.Unlock()

	if !client.closed {
		client.closed = true

		restricted := app.NewRestrictedMessageClient(client)
		for _, onClose := range client.onCloseDelegate {
			onClose(restricted)
		}
		close(client.done)
	}
	return nil
}

// Stop implements app.MessageSource.
func (client *MessageClient) Stop() {
	if client.stopped {
		return
	}

	client.mutex.Lock()
	defer client.mutex.Unlock()

	if !client.stopped {
		client.stopped = true
		close(client.stop)
	}
}

// Start implements app.MessageSource.
func (client *MessageClient) Start(pipe *app.MessagePipe) {
	var (
		upgrader = websocket.FastHTTPUpgrader{}
		ctx      = client.ctx
	)

	// setup WebSocketOption
	for _, opt := range client.options {
		opt.apply(&upgrader)
	}

	err := upgrader.Upgrade(ctx, func(ws *websocket.Conn) {
		defer func() {
			client.Close()
		}()
		defer ws.Close()

		var kontinue bool = true
		go func() {
			for kontinue {
				select {
				case v, ok := <-client.message:
					if ok {
						err := ws.WriteMessage(v.Type, v.Payload)
						if err != nil {
							pipe.Error(err)
						}
					}
				case <-client.stop:
					ws.SetReadLimit(0)
				case <-client.done:
					kontinue = false
					break
				}
			}
		}()
		for kontinue {
			mt, p, err := ws.ReadMessage()
			if err != nil {
				if _, ok := err.(*websocket.CloseError); ok {
					client.MessageClientInfo.Logger().Println(err)
					kontinue = false
					break
				}
				pipe.Error(err)
				continue
			}

			switch mt {
			case websocket.CloseMessage:
				message := &app.Message{
					Format: app.CLOSE_MESSAGE,
					Body:   p,
				}
				pipe.Forward(client, message)
				kontinue = false
				continue
			default:
				var message *app.Message
				format, ok := __MessageTypeMap[mt]
				if ok {
					message = &app.Message{
						Format: format,
						Body:   p,
					}
				} else {
					message = &app.Message{
						Format: app.UNKNOWN_MESSAGE,
						Body:   p,
					}
				}
				pipe.Forward(client, message)
			}
		}
	})

	if err != nil {
		pipe.Error(err)
	}
}

// Send implements app.MessageSource.
func (client *MessageClient) Send(format app.MessageFormat, payload []byte) error {
	mt, ok := __MessageFormatMap[format]
	if !ok {
		return fmt.Errorf("unsupported 'MessageFormat(%v)", format)
	}
	client.message <- &Message{
		Type:    mt,
		Payload: payload,
	}
	return nil
}
