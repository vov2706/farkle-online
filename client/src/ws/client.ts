import type { WsEnvelope } from "./types"

type Listener = (msg: WsEnvelope) => void

export class WsClient {
  private ws: WebSocket | null = null
  private listeners = new Set<Listener>()
  private reconnectTimer: number | null = null

  isConnected = false

  private readonly url: string

  constructor(authToken: string) {
    const proto = location.protocol === "https:" ? "wss" : "ws"
    this.url = `${proto}://${import.meta.env.VITE_WS_HOST}/ws?token=${encodeURIComponent(authToken)}`
  }

  connect() {
    if (this.ws && (this.ws.readyState === WebSocket.OPEN || this.ws.readyState === WebSocket.CONNECTING)) return

    this.ws = new WebSocket(this.url)

    this.ws.onopen = () => {
      this.isConnected = true
      this.emit({ type: "client.connected" })
    }

    this.ws.onclose = () => {
      this.isConnected = false
      this.emit({ type: "client.disconnected" })
      this.scheduleReconnect()
    }

    this.ws.onerror = () => {
      // onclose теж буде
    }

    this.ws.onmessage = (e) => {
      try {
        const msg = JSON.parse(e.data) as WsEnvelope
        this.emit(msg)
      } catch {
        // ignore
      }
    }
  }

  disconnect() {
    if (this.reconnectTimer) {
      window.clearTimeout(this.reconnectTimer)
      this.reconnectTimer = null
    }
    this.ws?.close()
    this.ws = null
  }

  send<T>(payload: WsEnvelope<T>) {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) return
    this.ws.send(JSON.stringify(payload))
  }

  on(fn: Listener) {
    this.listeners.add(fn)
    return () => this.listeners.delete(fn)
  }

  private emit(msg: WsEnvelope) {
    for (const fn of this.listeners) fn(msg)
  }

  private scheduleReconnect() {
    if (this.reconnectTimer) return
    this.reconnectTimer = window.setTimeout(() => {
      this.reconnectTimer = null
      this.connect()
    }, 800)
  }
}
