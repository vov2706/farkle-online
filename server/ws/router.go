package ws

import (
	"app/models"
	"encoding/json"

	"github.com/gofiber/contrib/v3/websocket"
)

type Envelope struct {
	Type    string          `json:"type"`
	Channel string          `json:"channel,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
}

type Ctx struct {
	Conn *websocket.Conn
	User *models.User
	Hub  *Hub
	Env  Envelope
}

type Handler func(ctx *Ctx) error

type Router struct {
	handlers map[string]Handler
	hub      *Hub
}

func NewRouter(hub *Hub) *Router {
	return &Router{
		handlers: make(map[string]Handler),
		hub:      hub,
	}
}

func (r *Router) On(event string, h Handler) {
	r.handlers[event] = h
}

func (r *Router) Handle(conn *websocket.Conn, user *models.User, raw []byte) {
	var env Envelope
	if err := json.Unmarshal(raw, &env); err != nil {
		r.hub.Send(conn, Event("error", "", map[string]any{"message": "bad json"}))
		return
	}

	h, ok := r.handlers[env.Type]
	if !ok {
		r.hub.Send(conn, Event("error", env.Channel, map[string]any{
			"message": "unknown event",
			"type":    env.Type,
		}))
		return
	}

	_ = h(&Ctx{
		Conn: conn,
		User: user,
		Hub:  r.hub,
		Env:  env,
	})
}
