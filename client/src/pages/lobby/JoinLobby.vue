<script setup lang="ts">
import { onMounted, ref, watch } from "vue";
import { useRouter } from "vue-router";
import TavernShell from "../../components/TavernShell.vue";
import {type Game, getGames, joinGame} from "@/api/game";
import CurrencyIcon, {type CurrencyType} from "@/components/CurrencyIcon.vue";
import type {PaginationMeta} from "@/interfaces/pagination.ts";

const router = useRouter();

const search = ref<string|undefined>();
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

  joinGame(game.id)
    .then(() => {
      router.push(`/lobby/${game.code}`);
    })
    .catch(err => {
      error.value = err?.message ?? "Failed to join";
    })
    .finally(() => {joinLoading.value = false})
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

          <button
            v-for="r in rooms"
            :key="r.code"
            class="w-full text-left rounded-xl border px-4 py-3 transition"
            :class="[
              isFull(r)
                ? 'border-danger-500/40 bg-danger-500/5 cursor-not-allowed opacity-70'
                : 'border-wood-700/35 bg-parchment-50/60 hover:bg-parchment-50/80'
            ]"
            :disabled="joinLoading || isFull(r)"
            @click="join(r)"
          >
            <div class="flex justify-between">
              <div>
                <div class="font-semibold">{{ r.code }}</div>
                <div class="flex gap-1 text-sm text-ink-900/65">
                  <span>
                    Host:
                    <span class="font-semibold">{{ r.creator?.username ?? "—" }}</span>
                  </span>
                  <div class="flex">
                    • Bet:
                    <div class="flex ml-1">
                      <span class="font-semibold"> {{ r.bet }}</span>
                      <CurrencyIcon v-if="r.currency" :type="r.currency.slug as CurrencyType" />
                    </div>
                  </div>
                  <span>
                    • Target:
                    <span class="font-semibold">{{ r.winning_points }} pts</span>
                  </span>
                  <span>
                    • Players:
                    <span class="font-semibold">{{ r.players_count }} / 2</span>
                  </span>
                </div>
              </div>

              <div
                class="font-semibold"
                :class="isFull(r) ? 'text-danger-500' : 'text-candle-400'"
              >
                <template v-if="isFull(r)">
                  Full
                </template>
                <template v-else>
                  Join →
                </template>
              </div>
            </div>
          </button>
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
