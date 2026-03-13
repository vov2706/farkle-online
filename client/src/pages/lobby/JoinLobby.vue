<script setup lang="ts">
import {onMounted, ref, watch} from "vue";
import {useRouter} from "vue-router";
import TavernShell from "@/components/wrappers/TavernShell.vue";
import {type Game, getGames, joinGame} from "@/api/game";
import CurrencyIcon from "@/components/icons/CurrencyIcon.vue";
import type {PaginationMeta} from "@/interfaces/pagination.ts";
import type {CurrencyType} from "@/api/currency.ts";

const router = useRouter();

const search = ref<string | undefined>();
const loading = ref(false);
const joinLoading = ref(false);
const error = ref("");
const maxPlayers = 2

const rooms = ref<Game[]>([]);
const meta = ref<PaginationMeta | null>(null);
const page = ref(1);

const isFull = (game: Game) => {
  return game.players_count >= maxPlayers;
};

const loadRooms = async () => {
  loading.value = true;
  try {
    const res = await getGames({page: page.value, search: search.value});
    rooms.value = res.data;
    meta.value = res.meta;
  } finally {
    loading.value = false;
  }
};

let debounceTimer: any;

watch(search, () => {
  clearTimeout(debounceTimer);

  debounceTimer = setTimeout(() => {
    page.value = 1;
    loadRooms();
  }, 400);
});

const nextPage = async () => {
  if (!meta.value?.has_more_pages) return;
  page.value++;
  await loadRooms();
};

const prevPage = async () => {
  if (page.value <= 1) return;
  page.value--;
  await loadRooms();
};

const join = async (game: Game) => {
  error.value = "";
  joinLoading.value = true;

  joinGame(game.code)
    .then(() => {
      router.push(`/lobby/${game.code}`);
    })
    .catch(err => {
      error.value = err?.message ?? "Failed to join";
    })
    .finally(() => {
      joinLoading.value = false
    })
};

onMounted(loadRooms);
</script>

<template>
  <TavernShell>
    <div class="mx-auto w-full max-w-7xl">
      <div class="panel">

        <div class="flex items-center justify-between">
          <div>
            <div class="font-display text-xl">Join room</div>
            <div class="text-sm text-ink-900/60">
              Search public rooms by enter a lobby code.
            </div>
          </div>

          <RouterLink
            to="/"
            class="text-sm font-semibold text-ink-900/70 hover:text-ink-900 transition"
          >
            Back to menu →
          </RouterLink>
        </div>

        <div
          v-if="error"
          class="mt-4 rounded-xl border border-danger-500/40 bg-danger-500/10 p-3 text-sm"
        >
          {{ error }}
        </div>

        <input
          v-model="search"
          class="mt-6 w-full rounded-xl border border-wood-700/35 bg-parchment-50 px-4 py-3 text-sm text-ink-900 outline-none
                 focus:ring-2 focus:ring-candle-400/60"
          placeholder="Search by room code…"
        />

        <!-- ROOMS LIST -->
        <div class="mt-6 space-y-2">
          <div v-if="loading" class="text-ink-900/60">Loading…</div>

          <div v-else-if="rooms.length === 0" class="text-ink-900/60">
            No rooms found.
          </div>

          <div class="space-y-3">
            <div
              v-for="r in rooms"
              :key="r.code"
              role="button"
              tabindex="0"
              class="group relative w-full overflow-hidden rounded-2xl border transition-all duration-200 hover:-translate-y-[2px] hover:shadow-lg"
              :class="[
                isFull(r)
                  ? 'border-danger-500/35 bg-danger-950/10 opacity-75 cursor-not-allowed'
                  : 'border-wood-700/35 bg-parchment-50/70 hover:-translate-y-0.5 hover:border-candle-400/40 hover:bg-parchment-50/90 hover:shadow-[0_12px_30px_rgba(0,0,0,0.18)] cursor-pointer'
              ]"
              @click="!joinLoading && !isFull(r) && join(r)"
              @keydown.enter="!joinLoading && !isFull(r) && join(r)"
              @keydown.space.prevent="!joinLoading && !isFull(r) && join(r)"
            >
              <div
                class="pointer-events-none absolute inset-0 opacity-0 transition-opacity duration-200"
                :class="isFull(r) ? '' : 'group-hover:opacity-100'"
                style="background: linear-gradient(90deg, rgba(255,255,255,0.00) 0%, rgba(255,244,214,0.18) 50%, rgba(255,255,255,0.00) 100%);"
              />

              <div class="relative flex items-center justify-between gap-4 px-5 py-4">
                <div class="min-w-0 flex-1">
                  <div class="flex flex-wrap items-center gap-2">
                    <div
                      class="flex h-10 w-10 items-center justify-center rounded-xl border border-wood-700/30 bg-wood-900/90 text-lg text-candle-300 shadow-inner">
                      🎲
                    </div>

                    <div class="min-w-0">
                      <div class="flex items-center gap-2">
                        <h3 class="truncate font-display text-lg tracking-wide text-ink-950">
                          Room {{ r.code }}
                        </h3>

                        <span
                          class="inline-flex items-center rounded-full border px-2 py-0.5 text-[11px] font-semibold uppercase tracking-[0.14em]"
                          :class="
                  isFull(r)
                    ? 'border-danger-500/30 bg-danger-500/10 text-danger-500'
                    : 'border-candle-500/30 bg-candle-400/10 text-candle-700'
                "
                        >
                {{ isFull(r) ? 'Full' : 'Open' }}
              </span>
                      </div>

                      <div class="mt-0.5 text-sm text-ink-900/65">
                        Host:
                        <span class="font-semibold text-ink-950">
                {{ r.creator?.username ?? '—' }}
              </span>
                      </div>
                    </div>
                  </div>

                  <div class="mt-4 grid grid-cols-2 gap-2 md:grid-cols-4">
                    <div class="rounded-xl border border-wood-700/20 bg-white/35 px-3 py-2">
                      <div class="text-[11px] uppercase tracking-[0.14em] text-ink-900/45">
                        Bet
                      </div>

                      <div class="mt-1 flex items-center gap-1.5 font-semibold text-ink-950">
                        <span>{{ r.bet }}</span>
                        <CurrencyIcon
                          v-if="r.currency"
                          :type="r.currency.slug as CurrencyType"
                        />
                      </div>
                    </div>

                    <div class="rounded-xl border border-wood-700/20 bg-white/35 px-3 py-2">
                      <div class="text-[11px] uppercase tracking-[0.14em] text-ink-900/45">
                        Target
                      </div>

                      <div class="mt-1 font-semibold text-ink-950">
                        {{ r.winning_points }} pts
                      </div>
                    </div>

                    <div class="rounded-xl border border-wood-700/20 bg-white/35 px-3 py-2">
                      <div class="text-[11px] uppercase tracking-[0.14em] text-ink-900/45">
                        Players
                      </div>

                      <div class="mt-1 font-semibold text-ink-950">
                        {{ r.players_count }} / 2
                      </div>
                    </div>

                    <div class="rounded-xl border border-wood-700/20 bg-white/35 px-3 py-2">
                      <div class="text-[11px] uppercase tracking-[0.14em] text-ink-900/45">
                        Entry
                      </div>

                      <div
                        class="mt-1 font-semibold"
                        :class="isFull(r) ? 'text-danger-500' : 'text-emerald-700'"
                      >
                        {{ isFull(r) ? 'Closed' : 'Available' }}
                      </div>
                    </div>
                  </div>
                </div>

                <div class="shrink-0 self-stretch">
                  <div
                    class="flex h-[72px] min-w-[140px] items-center justify-center rounded-xl border px-4 text-sm font-semibold transition-all duration-200"
                    :class="[
                      isFull(r)
                        ? 'border-danger-500/25 bg-danger-500/8 text-danger-500'
                        : 'border-candle-500/30 bg-candle-400/12 text-candle-700 group-hover:border-candle-400/50 group-hover:bg-candle-400/18'
                    ]"
                  >
                    <template v-if="joinLoading">
                      Joining...
                    </template>
                    <template v-else-if="isFull(r)">
                      No seats
                    </template>
                    <template v-else>
                      Join table →
                    </template>
                  </div>
                </div>
              </div>
            </div>
          </div>

        </div>

        <!-- PAGINATION -->
        <div
          v-if="meta"
          class="mt-6 flex items-center justify-between text-sm text-ink-900/70"
        >
          <div>
            Page {{ meta.current_page }} / {{ meta.last_page }}
            • Total: {{ meta.total }}
          </div>

          <div class="flex gap-3">
            <button
              class="hover:text-ink-900"
              :disabled="page === 1"
              @click="prevPage"
            >
              ← Prev
            </button>

            <button
              class="hover:text-ink-900"
              :disabled="!meta.has_more_pages"
              @click="nextPage"
            >
              Next →
            </button>
          </div>
        </div>

      </div>
    </div>
  </TavernShell>
</template>
