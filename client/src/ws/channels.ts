import type { WsEnvelope } from "./types"
import { WsClient } from "./client"

type ChannelListener = (msg: WsEnvelope) => void
type Key = string

export class WsChannels {
  private subscribed = new Set<string>()
  private listeners = new Map<Key, Set<ChannelListener>>()

  constructor(private client: WsClient) {
    this.client.on((msg) => {
      if (!msg.channel) return
      const key = this.key(msg.channel, msg.type)
      const set = this.listeners.get(key)
      if (!set) return
      for (const fn of set) fn(msg)
    })

    this.client.on((msg) => {
      if (msg.type !== "client.connected") return
      for (const ch of this.subscribed) {
        this.client.send({ type: "subscribe", data: { channel: ch } })
      }
    })
  }

  subscribe(channel: string) {
    if (this.subscribed.has(channel)) return
    this.subscribed.add(channel)
    this.client.send({ type: "subscribe", data: { channel } })
  }

  unsubscribe(channel: string) {
    if (!this.subscribed.has(channel)) return
    this.subscribed.delete(channel)
    this.client.send({ type: "unsubscribe", data: { channel } })

    for (const k of [...this.listeners.keys()]) {
      if (k.startsWith(channel + ":")) this.listeners.delete(k)
    }
  }

  emit<T>(channel: string, type: string, data: T) {
    this.client.send({ type, channel, data })
  }

  on(channel: string, type: string, fn: ChannelListener) {
    const key = this.key(channel, type)
    if (!this.listeners.has(key)) this.listeners.set(key, new Set())
    this.listeners.get(key)!.add(fn)
    return () => this.listeners.get(key)?.delete(fn)
  }

  private key(channel: string, type: string): Key {
    return `${channel}:${type}`
  }
}
