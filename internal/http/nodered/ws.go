package nodered

import (
	"sync"
	"time"

	"github.com/fasthttp/websocket"
	"github.com/pkg/errors"
	"touchon-server/internal/g"
	"touchon-server/lib/helpers"
	"touchon-server/lib/interfaces"
)

func NewWS(conn *websocket.Conn, tickerInterval time.Duration) (interfaces.WebSocket, error) {
	if conn == nil {
		return nil, errors.Wrap(errors.New("conn is nil"), "newWS")
	}

	// Устанавливаем обработчик Ping сообщений
	conn.SetPingHandler(func(appData string) error {
		g.Logger.Debug("Send Pong to NodeRed")
		return conn.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(time.Second))
	})

	// Устанавливаем PongHandler
	conn.SetPongHandler(func(appData string) error {
		return nil
	})

	o := &WebSocketImpl{
		mu:   &sync.Mutex{},
		conn: conn,
	}

	o.pinger = helpers.NewTimer(tickerInterval, o.ping)
	o.pinger.Start()

	return o, nil
}

type WebSocketImpl struct {
	mu     *sync.Mutex
	conn   *websocket.Conn
	pinger *helpers.Timer
}

func (o *WebSocketImpl) ReadMessage() (messageType int, p []byte, err error) {
	return o.conn.ReadMessage()
}

func (o *WebSocketImpl) WriteJSON(v interface{}) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	return o.conn.WriteJSON(v)
}

func (o *WebSocketImpl) WriteControl(messageType int, data []byte, deadline time.Time) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	return o.conn.WriteControl(messageType, data, deadline)
}

func (o *WebSocketImpl) Close() error {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.pinger.Stop()

	return o.conn.Close()
}

func (o *WebSocketImpl) ping() {
	if err := o.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(time.Second)); err != nil {
		g.Logger.Error("Failed to send Ping: ", err)
		return
	}

	g.Logger.Debug("Send Ping to NodeRed")

	o.pinger.Reset()
}

func (o *WebSocketImpl) ReadMessages(handler interfaces.WebSocketMsgHandler) {
	for {
		mt, message, err := o.conn.ReadMessage()
		if err != nil || mt == websocket.CloseMessage {
			break // Выходим из цикла, если клиент пытается закрыть соединение или связь прервана
		}

		handler(o, mt, message)
	}
}
