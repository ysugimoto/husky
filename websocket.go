package husky

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"time"
)

type WebSocketManager struct {
	clients map[string]*WebSocketDispatcher
}

func (wm *WebSocketManager) broadcast(message []byte, ignores ...string) {
	for connId, dispatcher := range wm.clients {
		canSend := true
		for _, id := range ignores {
			if dispatcher.Id == id {
				canSend = false
			}
		}
		if !canSend {
			continue
		}
		if err := dispatcher.Send(message); err != nil {
			log.Printf("[WebSocket] message send error to %s\n", connId)
		}
	}
}

type WebSocketDispatcher struct {
	conn      *websocket.Conn
	Id        string
	OnMessage func(msg string)
	OnOpen    func()
	OnClose   func()
	isClosed  bool
	manager   *WebSocketManager
}

func NewWebSocketDispatcher(ws *websocket.Conn) *WebSocketDispatcher {
	hash := sha1.Sum([]byte(fmt.Sprint(time.Now().UnixNano())))
	id := hex.EncodeToString(hash[:len(hash)])

	return &WebSocketDispatcher{
		Id:   id,
		conn: ws,
	}
}

func (w *WebSocketDispatcher) Broadcast(message []byte, ignores ...string) (err error) {
	if w.isClosed {
		return errors.New("Connection Closed")
	}

	w.manager.broadcast(message, ignores...)
	return nil
}

func (w *WebSocketDispatcher) Send(message []byte) (err error) {
	if w.isClosed {
		return errors.New("Connection Closed")
	}

	if err := websocket.Message.Send(w.conn, message); err != nil {
		return err
	}

	return nil
}

func (w *WebSocketDispatcher) Close() {
	if err := w.conn.Close(); err == nil {
		w.isClosed = true

		if w.OnClose != nil {
			w.OnClose()
		}
	}
}

func (w *WebSocketDispatcher) receive() (msg string, err error) {
	err = websocket.Message.Receive(w.conn, &msg)

	if err == nil && w.OnMessage != nil {
		w.OnMessage(msg)
	}

	return msg, err
}

func NewWebSocketHandler(handler HuskyWebSocketHandler) websocket.Handler {
	manager := &WebSocketManager{
		clients: make(map[string]*WebSocketDispatcher),
	}

	return websocket.Handler(func(ws *websocket.Conn) {
		dispatcher := NewWebSocketDispatcher(ws)
		dispatcher.manager = manager
		id := dispatcher.Id

		manager.clients[id] = dispatcher
		log.Printf("[WebSocket] id: %s connected.\n", id)

		defer func() {
			if err := ws.Close(); err == nil {
				if _, ok := manager.clients[id]; ok {
					if dispatcher.OnClose != nil {
						dispatcher.OnClose()
					}
					delete(manager.clients, id)
				}
			}
		}()

		handler(dispatcher)

		for {
			if dispatcher.isClosed {
				return
			}
			if msg, err := dispatcher.receive(); err != nil {
				log.Printf("[WebSocket] id: %s id disconnected.\n", id)
				delete(manager.clients, id)
				return
			} else {
				log.Printf("[WebSocket] message: %s\n", msg)
				manager.broadcast([]byte(msg))
			}
		}
	})
}
