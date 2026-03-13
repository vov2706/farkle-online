<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from "vue"
import { useRoute, useRouter } from "vue-router"
import TavernShell from "@/components/wrappers/TavernShell.vue"
import UiButton from "@/components/form/UiButton.vue"
import CurrencyIcon from "@/components/icons/CurrencyIcon.vue"

import { useToast } from "@/hooks/useToast.ts"
import { useConfirm } from "@/hooks/useConfirm.ts"
import { useAuthStore } from "@/stores/auth.ts"
import { getGame, type Player } from "@/api/game.ts"

import { WsClient } from "@/ws/client"
import { WsChannels } from "@/ws/channels"
import type { WsEnvelope } from "@/ws/types"
import GameStartOverlay from "@/components/overlays/GameStartOverlay.vue";

type LobbyPlayer = {
  name: string
  isHost: boolean
  isReady: boolean
}

const route = useRoute()
const router = useRouter()
const toast = useToast()
const confirm = useConfirm()

const auth = useAuthStore()
const me = computed(() => auth.user)

const maxPlayers = 2

const game = ref<any>(null)
const players = ref<LobbyPlayer[]>([])

const showStart = ref(false)
const startAt = ref<number>(3)

const roomCode = computed(() => String(route.params.code ?? "").toUpperCase())
const filled = computed(() => players.value.length)
const emptySlots = computed(() => Math.max(0, maxPlayers - players.value.length))

const myPlayer = computed(() => {
  if (!me.value?.username) return null
  return players.value.find((p) => p.name === me.value?.username) ?? null
})

const myReady = computed(() => !!myPlayer.value?.isReady)

const isHost = computed(() => {
  if (!me.value?.username) return false
  return players.value.some((p) => p.name === me.value?.username && p.isHost)
})

const allReady = computed(() => filled.value === maxPlayers && players.value.every((p) => p.isReady === true))
const canStart = computed(() => isHost.value && allReady.value)

const bet = computed(() => game.value?.bet ?? 0)
const currency = computed(() => game.value?.currency)
const winningPoints = computed(() => game.value?.winning_points ?? null)
const joinType = computed(() => game.value?.join_type ?? null)

const joinTypeLabel = computed(() => {
  switch (joinType.value) {
    case "anyone":
      return "Anyone can join"
    case "friends":
      return "Only friends"
    case "link":
      return "Join by link"
    default:
      return "—"
  }
})

const copyInviteLink = async () => {
  const link = `${window.location.origin}/lobby/${roomCode.value}`
  try {
    await navigator.clipboard.writeText(link)
    toast.push({ kind: "success", title: "Copied", message: "Invite link copied." })
  } catch {
    toast.push({ kind: "error", title: "Copy failed", message: "Clipboard is not available." })
  }
}

const onLeaveRoom = async () => {
  if (!game.value) return

  const ok = await confirm.open({
    title: "Leave the lobby?",
    message: "This action will reset your lobby settings.",
    confirmText: "Leave",
    cancelText: "Stay",
    loadingText: "Leaving...",
    tone: "danger",
  })

  if (!ok || !channels) return

  channels.emit(
    lobbyChannel.value,
    "lobby.player_left",
    { code: roomCode.value }
  )
}

const emitStartGame = () => {
  if (!canStart.value) {
    toast.push({ kind: "error", title: "Cannot start", message: "Every player must be ready." })
    return
  }

  channels?.emit(lobbyChannel.value, "start_game", {code: roomCode.value})
}

// Websockets
const wsConnected = ref(false)
let client: WsClient | null = null
let channels: WsChannels | null = null

const lobbyChannel = computed(() => `presence-lobby:${roomCode.value}`)
const listeners: Array<() => void> = []

const ensureMeInPlayers = () => {
  if (!me.value) return
  if (players.value.some((p) => p.name === me.value?.username)) return
  players.value.push({
    name: me.value?.username ?? "You",
    isHost: false,
    isReady: false,
  })
}

const initWs = () => {
  if (!auth.isLoggedIn) return
  if (client && channels) return

  client = new WsClient(auth.token)
  channels = new WsChannels(client)

  // connected/disconnected
  listeners.push(
    client.on((msg: WsEnvelope) => {
      console.log(msg)
      if (msg.type === "client.connected") wsConnected.value = true
      if (msg.type === "client.disconnected") wsConnected.value = false
    })
  )

  // subscribe lobby presence
  channels.subscribe(lobbyChannel.value)

  // ---- Presence базові івенти ----
  listeners.push(
    channels.on(lobbyChannel.value, 'presence.here', (msg) => {
      const users = (msg.data as any)?.users ?? []
      console.log("here", users)
    })
  )

  listeners.push(
    channels.on(lobbyChannel.value, 'presence.joining', (_) => {
      // if user is online
    })
  )

  listeners.push(
    channels.on(lobbyChannel.value, 'presence.leaving', (_) => {
      // if user is not online
    })
  )

  listeners.push(
    channels.on(lobbyChannel.value, 'lobby.player_ready', (msg) => {
      const data: any = msg.data ?? {}
      const username = data.player
      const isReady = !!data.is_ready

      if (!username) return
      const p = players.value.find((x) => x.name === username)
      if (p) p.isReady = isReady
    })
  )

  listeners.push(
    channels.on(lobbyChannel.value, 'lobby.player_left', async (msg) => {
      const data: any = msg.data ?? {}
      const username = data.player

      if (!username) return
      const idx = players.value.findIndex((x) => x.name === username)
      players.value.splice(idx, 1)

      await auth.fetchProfile()

      if (auth.user?.username !== username) {
        toast.push({ kind: "success", title: 'Notification', message: `Player ${username} left the lobby` })
      } else {
        await router.push('/')
      }
    })
  )

  listeners.push(
    channels.on(lobbyChannel.value, 'lobby.player_joined', async (msg) => {
      const data: any = msg.data ?? {}
      const username = data.player

      if (!username) return

      players.value.push({
        name: username,
        isHost: false,
        isReady: false,
      })
    })
  )

  listeners.push(
    channels.on(lobbyChannel.value, 'lobby.game_deleted', (msg) => {
      const data: any = msg.data ?? {}

      router.push("/")

      toast.push({ kind: "success", title: 'Notification', message: `Lobby ${data.code} was deleted` })
    })
  )

  listeners.push(
    channels.on(lobbyChannel.value, 'game.started', (msg) => {
      showStart.value = true
      startAt.value = msg.data.start_time
    })
  )

  client.connect()
}

const toggleReady = () => {
  if (!me.value?.username) return
  if (!channels) return

  const p = players.value.find((x) => x.name === me.value!.username)
  if (!p) return

  channels.emit(
    lobbyChannel.value,
    "lobby.player_ready",
    { code: roomCode.value, is_ready: !p.isReady }
  )
}

const startGame = () => {
  showStart.value = false

  router.push(`/game/${roomCode.value}`)
}

onMounted(async () => {
  game.value = await getGame(String(route.params.code))

  players.value = (game.value.players as Player[]).map((p: Player) => ({
    name: p.username,
    isHost: p.is_host,
    isReady: false,
  }))

  ensureMeInPlayers()
  initWs()
})

onBeforeUnmount(() => {
  try {
    channels?.unsubscribe(lobbyChannel.value)
  } catch {}

  for (const off of listeners) {
    try {
      off()
    } catch {}
  }
  listeners.length = 0

  client?.disconnect()
  client = null
  channels = null
})
</script>


<template>
  <TavernShell>
    <div class="panel">
      <div class="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
        <div>
          <div class="font-display text-2xl">Lobby</div>
          <div class="text-lg text-ink-900/70 flex flex-wrap items-center gap-3 mt-2">
            <span>
              Room <span class="font-semibold">{{ roomCode }}</span>
            </span>
          </div>
        </div>
        <div class="flex flex-wrap items-center gap-2">
          <UiButton variant="ghost" @click="copyInviteLink">Copy invite link</UiButton>
          <UiButton variant="ghost" @click="onLeaveRoom">Leave</UiButton>
        </div>
      </div>

      <div class="mt-6 grid grid-cols-1 gap-6 lg:grid-cols-12">
        <section class="lg:col-span-7">
          <div class="font-display text-lg">Players</div>

          <div class="mt-3 grid grid-cols-1 gap-2 sm:grid-cols-2">
            <div
              v-for="p in players"
              :key="p.name"
              class="flex items-center justify-between rounded-xl border border-wood-700/35 bg-tavern-900/50 px-4 py-3"
            >
              <div class="flex items-center gap-3 min-w-0">
                <div class="h-9 w-9 rounded-full bg-candle-300/35 flex items-center justify-center shrink-0">
                  🎲
                </div>

                <div class="min-w-0">
                  <div class="font-semibold text-parchment-50 truncate">
                    {{ p.name }}
                    <span v-if="p.isHost" class="ml-1 text-xs text-candle-300">(host)</span>
                    <span v-if="p.name === me?.username" class="ml-1 text-xs text-parchment-50/60">(you)</span>
                  </div>
                  <div class="text-xs text-parchment-50/60">
                    {{ p.isReady ? "Ready" : "Not ready" }}
                  </div>
                </div>
              </div>

              <div
                class="text-xs px-2 py-1 rounded-full shrink-0"
                :class="p.isReady
                  ? 'bg-candle-300/30 text-parchment-50'
                  : 'bg-tavern-800/50 text-parchment-50/60'"
              >
                {{ p.isReady ? "Ready" : "Waiting" }}
              </div>
            </div>

            <div
              v-for="i in emptySlots"
              :key="'empty-' + i"
              class="flex items-center justify-center rounded-xl border border-dashed border-wood-700/35
                     bg-tavern-900/30 px-4 py-3 text-parchment-50/35"
            >
              Empty slot
            </div>
          </div>

          <div class="mt-4 text-sm text-ink-900/60">
            To start the game: everyone player must be <span class="font-semibold">Ready</span>.
          </div>
          <div class="mt-5 rounded-tavern border border-wood-700/30 bg-parchment-50/60 p-4">
            <div class="flex items-center justify-between">
              <div class="font-display text-base">Room rules</div>
              <div class="text-xs text-ink-900/60">Summary</div>
            </div>

            <div class="mt-3 grid grid-cols-1 gap-2 sm:grid-cols-2">
              <div class="flex items-center justify-between rounded-xl border border-wood-700/25 bg-parchment-50/40 px-3 py-2">
                <div class="flex items-center gap-2 text-ink-900/70">
                  <span class="text-base">🎯</span>
                  <span class="text-sm">Winning</span>
                </div>
                <div class="font-semibold">
                  <span v-if="winningPoints">{{ winningPoints }} pts</span>
                  <span v-else class="text-ink-900/50">—</span>
                </div>
              </div>

              <div class="flex items-center justify-between rounded-xl border border-wood-700/25 bg-parchment-50/40 px-3 py-2">
                <div class="flex items-center gap-2 text-ink-900/70">
                  <span class="text-base">🔒</span>
                  <span class="text-sm">Join</span>
                </div>
                <div class="font-semibold">{{ joinTypeLabel }}</div>
              </div>

              <div class="flex items-center justify-between rounded-xl border border-wood-700/25 bg-parchment-50/40 px-3 py-2">
                <div class="flex items-center gap-2 text-ink-900/70">
                  <span class="text-base">💰</span>
                  <span class="text-sm">Bet</span>
                </div>
                <div class="flex items-center gap-1 font-semibold">
                  <span>{{ bet }}</span>
                  <CurrencyIcon v-if="currency" :type="currency.slug" :size="16" />
                </div>
              </div>

              <div class="flex items-center justify-between rounded-xl border border-wood-700/25 bg-parchment-50/40 px-3 py-2">
                <div class="flex items-center gap-2 text-ink-900/70">
                  <span class="text-base">👥</span>
                  <span class="text-sm">Players</span>
                </div>
                <div class="font-semibold">{{ filled }}/{{ maxPlayers }}</div>
              </div>
            </div>

            <div class="mt-3 text-xs text-ink-900/60">
              Tip: Ready resets on refresh. The host can start when everyone is ready.
            </div>
          </div>

        </section>

        <!-- RIGHT: CONTROL -->
        <aside class="lg:col-span-5">
          <div class="lg:h-full">
            <div class="font-display text-lg">Controls</div>

            <div class="mt-3 rounded-tavern border border-wood-700/30 bg-parchment-50/60 p-6">
              <div class="text-base font-semibold">Room status</div>

              <div class="mt-3 space-y-2">
                <div class="flex items-center justify-between">
                  <span class="text-base text-ink-900/70">All ready</span>
                  <span class="text-base font-semibold">{{ allReady ? "Yes" : "No" }}</span>
                </div>

                <div class="flex items-center justify-between">
                  <span class="text-base text-ink-900/70">You</span>
                  <span class="text-base font-semibold">{{ myReady ? "Ready" : "Not ready" }}</span>
                </div>

                <div class="flex items-center justify-between">
                  <span class="text-base text-ink-900/70">Host</span>
                  <span class="text-base font-semibold">{{ isHost ? "You" : "Other" }}</span>
                </div>
              </div>

              <div class="mt-4 flex flex-col gap-2">
                <UiButton
                  class="py-3 text-base"
                  :variant="myReady ? 'ghost' : 'primary'"
                  @click="toggleReady"
                >
                  {{ myReady ? "Unready" : "Ready" }}
                </UiButton>

                <UiButton
                  class="py-3 text-base"
                  v-if="isHost && allReady"
                  variant="primary"
                  @click="emitStartGame"
                >
                  Start game
                </UiButton>

                <UiButton
                  v-if="isHost && maxPlayers > players.length"
                  class="py-3 text-base"
                  variant="ghost"
                  @click="copyInviteLink"
                >
                  Invite
                </UiButton>
              </div>
            </div>
          </div>
        </aside>
      </div>
    </div>
    <GameStartOverlay
      :show="showStart"
      :start-at="startAt"
      @done="startGame"
    />
  </TavernShell>
</template>
