package interfaces

import (
	"time"

	"github.com/valyala/fasthttp"
)

type WebSocketMsgHandler func(ws WebSocket, messageType int, message []byte)

type WebSocket interface {
	ReadMessage() (messageType int, p []byte, err error)
	WriteJSON(v interface{}) error
	WriteControl(messageType int, data []byte, deadline time.Time) error
	Close() error
	ReadMessages(handler WebSocketMsgHandler)
}

type NodeRed interface {
	Handler(ctx *fasthttp.RequestCtx)
	Start() error
	Shutdown() error
}
