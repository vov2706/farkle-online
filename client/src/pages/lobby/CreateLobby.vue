<script setup lang="ts">
import {computed, onMounted, ref} from 'vue'
import TavernShell from '@/components/wrappers/TavernShell.vue'
import UiButton from '@/components/form/UiButton.vue'
import { useRouter } from 'vue-router'
import {type Currency, CurrencyType, getCurrencies} from "@/api/currency.ts";
import UiSelect from "@/components/form/UiSelect.vue";
import {createGame, type CreateGameInput, type Game, JoinType} from "@/api/game.ts";
import CurrencyIcon from "@/components/icons/CurrencyIcon.vue";

const types = [
  {slug: JoinType.ANYONE, title: 'Anyone can join', subtitle: 'Visible & open'},
  {slug: JoinType.FRIENDS, title: 'Only friends', subtitle: 'Friends list only'},
  {slug: JoinType.LINK, title: 'Join by link', subtitle: 'Invite URL'},
]

const router = useRouter()
const currencies = ref<Currency[]>([])

const currencyId = ref<number|undefined>(currencies.value[0]?.id)
const bet = ref<number>(10)
const winningPoints = ref<number>(3000)
const joinType = ref<JoinType>(JoinType.ANYONE)

const loading = ref(false)
const error = ref('')
const game = ref<Game | null>(null)

const betValid = computed(() => Number.isFinite(bet.value) && bet.value > 0)
const pointsValid = computed(() => Number.isFinite(winningPoints.value) && winningPoints.value > 0)
const currencyValid = computed(() => Number.isFinite(currencyId.value) && (currencyId.value ?? 0) > 0)
const joinTypeValid = computed(() => !!joinType.value)

const canCreate = computed(() => !loading.value && betValid.value && pointsValid.value && currencyValid.value && joinTypeValid.value)

const selectedCurrency = computed(() => currencies.value.find(c => c.id === currencyId.value))

const createRoom = async () => {
  if (!currencyId.value) {
    return;
  }

  error.value = ''
  loading.value = true
  try {
    const input: CreateGameInput = {
      currency_id: currencyId.value,
      bet: bet.value,
      winning_points: winningPoints.value,
      join_type: joinType.value,
    }

    game.value = await createGame(input)

    await router.push(`/lobby/${game.value.code}`)
  } catch (e: any) {
    error.value = e?.message ?? 'Failed to create game'
  } finally {
    loading.value = false
  }
}

const reset = () => {
  game.value = null
  error.value = ''
}

const copyCode = async () => {
  if (!game.value) return
  try {
    await navigator.clipboard.writeText(game.value.code)
  } catch {
    // ignore
  }
}

onMounted(async () => {
  currencies.value = await getCurrencies()
})

</script>

<template>
  <TavernShell>
    <div class="panel">
      <div class="flex flex-col gap-1 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <div class="font-display text-xl">Create game</div>
          <div class="text-sm text-ink-900/70">
            Set rules and invite players.
          </div>
        </div>

        <RouterLink to="/" class="text-sm text-ink-900/60 hover:text-ink-900">
          ← Back to menu
        </RouterLink>
      </div>

      <div v-if="error" class="mt-4 rounded-xl border border-danger-500/40 bg-danger-500/10 p-3 text-sm">
        {{ error }}
      </div>

      <!-- BEFORE CREATED -->
      <div v-if="!game" class="mt-6 grid grid-cols-1 gap-6 lg:grid-cols-12">
        <div class="lg:col-span-7">
          <div class="font-display text-lg">Room settings</div>

          <!-- bet + currency + winning points -->
          <div class="mt-6 grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div>
              <div class="text-base font-semibold">Currency</div>

              <div class="mt-2">
                <UiSelect
                  v-model="currencyId"
                  :options="currencies.map((c: Currency) => ({
                    label: c.name,
                    value: c.id,
                    type: c.slug,
                  }))"
                >
                  <template #option="{ option }">
                    <div class="flex items-center gap-2">
                      <CurrencyIcon :type="option.type as CurrencyType" :size="16" />
                      {{ option.label }}
                    </div>
                  </template>
<!--                  <slot name="option">-->
<!--                    <img v-if="option.icon" :src="option.icon" class="w-5 h-5" />-->
<!--                    <span class="truncate">{{ option.label }}</span>-->
<!--                  </slot>-->
                </UiSelect>
                <div v-if="!currencyValid" class="mt-1 text-base text-danger-600">
                  Select a currency.
                </div>
              </div>
            </div>

            <div>
              <div class="text-base font-semibold">Bet</div>

              <div class="mt-2">
                <input
                  v-model.number="bet"
                  type="number"
                  min="1"
                  step="1"
                  class="w-full rounded-xl border border-wood-700/35 bg-parchment-50/70 px-4 py-2 text-base outline-none transition focus:border-candle-400/60"
                  placeholder="e.g. 10"
                />
                <div v-if="!betValid" class="mt-1 text-xs text-danger-600">
                  Bet must be greater than 0.
                </div>
              </div>
            </div>

            <div class="sm:col-span-2">
              <div class="text-base font-semibold">Winning points</div>
              <div class="mt-2">
                <input
                  v-model.number="winningPoints"
                  type="number"
                  min="1"
                  step="50"
                  class="w-full rounded-xl border border-wood-700/35 bg-parchment-50/70 px-4 py-2 text-base outline-none transition focus:border-candle-400/60"
                  placeholder="e.g. 3000"
                />
                <div v-if="!pointsValid" class="mt-1 text-xs text-danger-600">
                  Winning points must be greater than 0.
                </div>
              </div>

              <div class="mt-2 text-sm text-ink-900/70">
                First player to reach <span class="font-semibold">{{ winningPoints }}</span> points wins the match.
              </div>
            </div>
          </div>

          <!-- join type -->
          <div class="mt-6">
            <div class="text-base font-semibold">Who can join</div>
            <div class="mt-2 grid grid-cols-1 gap-2 sm:grid-cols-3">
              <button
                v-for="jt in types"
                :key="jt.slug"
                class="rounded-xl border px-4 py-3 text-left transition"
                :class="joinType === jt.slug
                  ? 'border-candle-400/60 bg-candle-300/30'
                  : 'border-wood-700/35 bg-parchment-50/50 hover:bg-parchment-50/70'"
                @click="joinType = jt.slug"
              >
                <span class="font-semibold">{{ jt.title }}</span><br/>
                <span class="text-sm text-ink-900/70">{{ jt.subtitle }}</span>
              </button>
            </div>
          </div>

          <div class="mt-6 flex gap-2">
            <UiButton variant="primary" :disabled="!canCreate" @click="createRoom">
              {{ loading ? 'Creating…' : 'Create game' }}
            </UiButton>

            <RouterLink to="/">
              <UiButton variant="ghost">Back</UiButton>
            </RouterLink>
          </div>
        </div>

        <!-- preview -->
        <div class="lg:col-span-5">
          <div class="font-display text-lg">Preview</div>

          <div class="mt-3 rounded-tavern border border-wood-700/30 bg-tavern-900/50 p-4">
            <div class="text-xl text-parchment-50/70">Room will be created as:</div>

            <div class="mt-3 space-y-2 text-lg">
              <div class="chip flex gap-1">
                💰 Bet: <span class="font-bold">{{ bet }}</span>
                <CurrencyIcon v-if="selectedCurrency" :type="selectedCurrency.slug as CurrencyType" :size="20" />
              </div>
              <div class="chip">
                🏁 Winning points: <span class="font-bold">{{ winningPoints }}</span>
              </div>
              <div class="chip">
                🔒 Join:
                <span class="font-bold underline" v-if="joinType === JoinType.ANYONE">Anyone can join</span>
                <span class="font-bold underline" v-else-if="joinType === JoinType.FRIENDS">Only friends</span>
                <span class="font-bold underline" v-else>Join by link</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- AFTER CREATED -->
      <div v-else class="mt-6 grid grid-cols-1 gap-6 lg:grid-cols-12">
        <div class="lg:col-span-7">
          <div class="font-display text-lg">Room created</div>
          <div class="mt-1 text-sm text-ink-900/70">
            Share this code with players to join your game.
          </div>

          <div class="mt-4 rounded-[var(--radius-tavern)] border border-wood-700/35 bg-parchment-50/60 p-4">
            <div class="text-sm font-semibold">Room code</div>

            <div class="mt-2 flex flex-col gap-2 sm:flex-row sm:items-center">
              <div
                class="flex-1 rounded-xl border border-wood-700/35 bg-parchment-50 px-4 py-3 font-display text-xl tracking-[0.08em]"
              >
                {{ game.code }}
              </div>

              <UiButton variant="primary" @click="copyCode">Copy</UiButton>
            </div>

<!--            <div class="mt-3 text-sm text-ink-900/70">-->
<!--              Max players: <span class="font-semibold">{{ game.maxPlayers }}</span> •-->
<!--              Status: <span class="font-semibold">{{ game.status }}</span>-->
<!--            </div>-->
          </div>

          <div class="mt-5 flex flex-wrap gap-2">
            <RouterLink to="/join">
              <UiButton variant="ghost">Go to Join page</UiButton>
            </RouterLink>

            <UiButton variant="primary" @click="reset">Create another</UiButton>

            <RouterLink to="/">
              <UiButton variant="ghost">Back</UiButton>
            </RouterLink>
          </div>
        </div>

        <div class="lg:col-span-5">
          <div class="font-display text-lg">Next steps</div>

          <div class="mt-3 space-y-3">
            <div class="rounded-[var(--radius-tavern)] border border-wood-700/30 bg-tavern-900/50 p-4">
              <div class="text-parchment-50 font-semibold">Invite players</div>
              <div class="text-parchment-50/70 text-sm mt-1">
                Send the code or link (depending on join type).
              </div>
            </div>

            <div class="rounded-[var(--radius-tavern)] border border-wood-700/30 bg-tavern-900/50 p-4">
              <div class="text-parchment-50 font-semibold">Start match</div>
              <div class="text-parchment-50/70 text-sm mt-1">
                Route to <span class="font-semibold">/game/:code</span> and plug Three.js there.
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </TavernShell>
</template>
