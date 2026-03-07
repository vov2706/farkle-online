package ws

import (
	"encoding/json"
	"sync"

	"github.com/gofiber/contrib/v3/websocket"
)

type Hub struct {
	mu sync.RWMutex

	// channel -> conns
	channels map[string]map[*websocket.Conn]struct{}

	// conn -> channels
	connChannels map[*websocket.Conn]map[string]struct{}

	// presence: channel -> conn -> userInfo
	presence map[string]map[*websocket.Conn]map[string]any
}

func NewHub() *Hub {
	return &Hub{
		channels:     make(map[string]map[*websocket.Conn]struct{}),
		connChannels: make(map[*websocket.Conn]map[string]struct{}),
		presence:     make(map[string]map[*websocket.Conn]map[string]any),
	}
}

func (h *Hub) Join(channel string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.channels[channel]; !ok {
		h.channels[channel] = make(map[*websocket.Conn]struct{})
	}
	h.channels[channel][conn] = struct{}{}

	if _, ok := h.connChannels[conn]; !ok {
		h.connChannels[conn] = make(map[string]struct{})
	}
	h.connChannels[conn][channel] = struct{}{}
}

func (h *Hub) JoinPresence(channel string, conn *websocket.Conn, userInfo map[string]any) (here []map[string]any, isFirst bool) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// normal join
	if _, ok := h.channels[channel]; !ok {
		h.channels[channel] = make(map[*websocket.Conn]struct{})
	}
	h.channels[channel][conn] = struct{}{}

	if _, ok := h.connChannels[conn]; !ok {
		h.connChannels[conn] = make(map[string]struct{})
	}
	h.connChannels[conn][channel] = struct{}{}

	// presence map
	if _, ok := h.presence[channel]; !ok {
		h.presence[channel] = make(map[*websocket.Conn]map[string]any)
	}
	_, existed := h.presence[channel][conn]
	h.presence[channel][conn] = userInfo

	// build "here"
	users := make([]map[string]any, 0, len(h.presence[channel]))
	for _, u := range h.presence[channel] {
		users = append(users, u)
	}

	return users, !existed
}

func (h *Hub) Leave(channel string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if conns, ok := h.channels[channel]; ok {
		delete(conns, conn)
		if len(conns) == 0 {
			delete(h.channels, channel)
		}
	}

	if chans, ok := h.connChannels[conn]; ok {
		delete(chans, channel)
		if len(chans) == 0 {
			delete(h.connChannels, conn)
		}
	}

	if pres, ok := h.presence[channel]; ok {
		delete(pres, conn)
		if len(pres) == 0 {
			delete(h.presence, channel)
		}
	}
}

func (h *Hub) LeaveAll(conn *websocket.Conn) (left []string, presenceLeft map[string]map[string]any) {
	h.mu.Lock()
	defer h.mu.Unlock()

	presenceLeft = map[string]map[string]any{}

	chans := h.connChannels[conn]
	for ch := range chans {
		// presence info before delete
		if pres, ok := h.presence[ch]; ok {
			if u, ok2 := pres[conn]; ok2 {
				presenceLeft[ch] = u
				delete(pres, conn)
				if len(pres) == 0 {
					delete(h.presence, ch)
				}
			}
		}

		// remove from channel
		if conns, ok := h.channels[ch]; ok {
			delete(conns, conn)
			if len(conns) == 0 {
				delete(h.channels, ch)
			}
		}

		left = append(left, ch)
	}

	delete(h.connChannels, conn)
	return left, presenceLeft
}

func (h *Hub) Broadcast(channel string, msg []byte) {
	h.mu.RLock()
	conns := h.channels[channel]
	h.mu.RUnlock()

	for c := range conns {
		_ = c.WriteMessage(websocket.TextMessage, msg)
	}
}

func (h *Hub) Send(conn *websocket.Conn, msg []byte) {
	_ = conn.WriteMessage(websocket.TextMessage, msg)
}

func Event(eventType string, channel string, data any) []byte {
	out := map[string]any{"type": eventType}
	if channel != "" {
		out["channel"] = channel
	}
	if data != nil {
		out["data"] = data
	}
	b, _ := json.Marshal(out)
	return b
}

func IsPresenceChannel(channel string) bool {
	return len(channel) >= 9 && channel[:9] == "presence-"
}
