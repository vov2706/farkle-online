export type WsEnvelope<T = any> = {
  type: string
  channel?: string
  data?: T
}

export type SubscribePayload = { channel: string }

export type PresenceHere = { users: Array<{ id: number; username: string }> }
export type PresenceJoining = { user: { id: number; username: string } }
export type PresenceLeaving = { user: { id: number; username: string } }
