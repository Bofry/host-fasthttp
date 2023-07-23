package websocket

import "github.com/Bofry/host/app"

var (
	_ app.StandardProtocol = pingProtocolImpl(0)

	__PONG_BYTES = []byte("pong")
)

type pingProtocolImpl int

// ConfigureProtocol implements app.StandardProtocol.
func (p pingProtocolImpl) ConfigureProtocol(registry *app.StandardProtocolRegistry) {
	var (
		body = []byte("ping")
	)

	for _, format := range app.StandardProtocolMessageFormats {
		registry.Add(app.Message{
			Format: format,
			Body:   body,
		}, p)
	}
}

// ReplyMessage implements app.StandardProtocol.
func (pingProtocolImpl) ReplyMessage(format app.MessageFormat, sender app.MessageSender) {
	sender.Send(format, __PONG_BYTES)
}
