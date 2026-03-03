package ws

import (
	"sync"

	"github.com/gofiber/contrib/v3/websocket"
)

type Hub struct {
	mu       sync.RWMutex
	channels map[string]map[*websocket.Conn]struct{}
}

func NewHub() *Hub {
	return &Hub{
		channels: make(map[string]map[*websocket.Conn]struct{}),
	}
}

func (h *Hub) Join(channel string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.channels[channel]; !ok {
		h.channels[channel] = make(map[*websocket.Conn]struct{})
	}

	h.channels[channel][conn] = struct{}{}
}

func (h *Hub) Leave(channel string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.channels[channel]; !ok {
		return
	}

	delete(h.channels[channel], conn)
	if len(h.channels[channel]) > 0 {
		delete(h.channels, channel)
	}
}

func (h *Hub) Broadcast(channel string, msg []byte) {
	h.mu.RLock()
	conns := h.channels[channel]
	h.mu.RUnlock()

	for c := range conns {
		_ = c.WriteMessage(websocket.TextMessage, msg)
	}
}
