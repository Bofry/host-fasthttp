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
	ctx *fasthttp.RequestCtx

	onCloseDelegate []func(app.MessageClient)

	message chan *Message
	stop    chan struct{}
	done    chan struct{}

	stopped bool
	closed  bool
	mutex   sync.Mutex
}

func NewMessageClient(ctx *fasthttp.RequestCtx) *MessageClient {
	return &MessageClient{
		ctx:     ctx,
		message: make(chan *Message),
		stop:    make(chan struct{}),
		done:    make(chan struct{}),
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
		for _, onClose := range client.onCloseDelegate {
			onClose(client)
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
	defer func() {
		err := recover()
		if err != nil {
			if verr, ok := err.(error); ok {
				pipe.Error(verr)
			} else {
				pipe.Error(fmt.Errorf("%v", err))
			}
		}
	}()

	var (
		upgrader = websocket.FastHTTPUpgrader{}
		ctx      = client.ctx
	)
	err := upgrader.Upgrade(ctx, func(ws *websocket.Conn) {
		defer func() {
			client.Close()
		}()
		defer ws.Close()

		var kontinue bool = true
		for kontinue {
			mt, p, err := ws.ReadMessage()
			if err != nil {
				if _, ok := err.(*websocket.CloseError); ok {
					pipe.Error(err)
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

			select {
			case v := <-client.message:
				err := ws.WriteMessage(v.Type, v.Payload)
				if err != nil {
					pipe.Error(err)
				}
			case <-client.stop:
				ws.SetReadLimit(0)
			case <-client.done:
				kontinue = false
				continue
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
