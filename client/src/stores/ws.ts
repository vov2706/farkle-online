import { defineStore } from "pinia"
import { WsClient } from "@/ws/client"
import { WsChannels } from "@/ws/channels"
import type { PresenceHere, PresenceJoining, PresenceLeaving, WsEnvelope } from "@/ws/types"

export const useWsStore = defineStore("ws", {
  state: () => ({
    connected: false,
    client: null as WsClient | null,
    channels: null as WsChannels | null,

    // простий приклад presence-стану по channel
    presence: {} as Record<string, Array<{ id: number; username: string }>>,
  }),

  actions: {
    init() {
      if (this.client && this.channels) return

      // IMPORTANT:
      // якщо твій WsProtected middleware використовує Cookie/JWT в cookie — ок.
      // якщо токен в header — WebSocket header так просто не засунеш -> краще токен в query (?token=)
      const url = `${location.protocol === "https:" ? "wss" : "ws"}://${location.host}/ws`

      const client = new WsClient(url)
      const channels = new WsChannels(client)

      client.on((msg: WsEnvelope) => {
        if (msg.type === "client.connected") this.connected = true
        if (msg.type === "client.disconnected") this.connected = false
      })

      this.client = client
      this.channels = channels
      client.connect()
    },

    joinPresenceLobby(code: string) {
      this.init()
      const ch = `presence-lobby:${code}`

      this.channels!.subscribe(ch)

      // here
      this.channels!.on(ch, "presence.here", (msg) => {
        const data = msg.data as PresenceHere
        this.presence[ch] = data.users
      })

      // joining
      this.channels!.on(ch, "presence.joining", (msg) => {
        const data = msg.data as PresenceJoining
        const list = this.presence[ch] ?? []
        if (!list.find((u) => u.id === data.user.id)) list.push(data.user)
        this.presence[ch] = list
      })

      // leaving
      this.channels!.on(ch, "presence.leaving", (msg) => {
        const data = msg.data as PresenceLeaving
        this.presence[ch] = (this.presence[ch] ?? []).filter((u) => u.id !== data.user.id)
      })

      return ch
    },

    leave(channel: string) {
      if (!this.channels) return
      this.channels.unsubscribe(channel)
      delete this.presence[channel]
    },

    lobbyReady(channel: string, isReady: boolean) {
      this.channels?.emit(channel, "lobby.player_ready", { is_ready: isReady })
    },
  },
})
