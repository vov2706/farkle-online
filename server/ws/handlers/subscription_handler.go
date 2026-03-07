package ws

import (
	"app/ws"
	"encoding/json"
)

type SubscribeInput struct {
	Channel string `json:"channel"`
}

func SubscribeHandler(authorizer *ws.Authorizer) ws.Handler {
	return func(ctx *ws.Ctx) error {
		var in SubscribeInput
		if err := json.Unmarshal(ctx.Env.Data, &in); err != nil || in.Channel == "" {
			ctx.Hub.Send(ctx.Conn, ws.Event("error", "", map[string]any{"message": "bad subscribe payload"}))
			return nil
		}

		ok, meta, err := authorizer.Authorize(ctx.User, in.Channel)

		if err != nil {
			ctx.Hub.Send(ctx.Conn, ws.Event("error", in.Channel, map[string]any{"message": "authorize error"}))
			return nil
		}
		if !ok {
			ctx.Hub.Send(ctx.Conn, ws.Event("subscription_error", in.Channel, map[string]any{"message": "forbidden"}))
			return nil
		}

		// presence: here/joining
		if ws.IsPresenceChannel(in.Channel) {
			here, _ := ctx.Hub.JoinPresence(in.Channel, ctx.Conn, meta)

			// send "here" only to this conn
			ctx.Hub.Send(ctx.Conn, ws.Event("presence.here", in.Channel, map[string]any{
				"users": here,
			}))

			// broadcast "joining" to others
			ctx.Hub.Broadcast(in.Channel, ws.Event("presence.joining", in.Channel, map[string]any{
				"user": meta,
			}))
		} else {
			ctx.Hub.Join(in.Channel, ctx.Conn)
		}

		ctx.Hub.Send(ctx.Conn, ws.Event("subscribed", in.Channel, map[string]any{"ok": true}))
		return nil
	}
}

func UnsubscribeHandler() ws.Handler {
	return func(ctx *ws.Ctx) error {
		var in SubscribeInput
		if err := json.Unmarshal(ctx.Env.Data, &in); err != nil || in.Channel == "" {
			return nil
		}
		ctx.Hub.Leave(in.Channel, ctx.Conn)
		ctx.Hub.Send(ctx.Conn, ws.Event("unsubscribed", in.Channel, map[string]any{"ok": true}))
		return nil
	}
}
