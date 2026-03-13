<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from "vue";

type Tone = "default" | "danger" | "success";

const props = withDefaults(
  defineProps<{
    modelValue: boolean;

    title?: string;
    message?: string;

    confirmText?: string;
    cancelText?: string;
    loadingText?: string;

    tone?: Tone;

    loading?: boolean;
    disableClose?: boolean;

    // опційно: показати маленький текст під кнопками (Enter/Esc)
    hint?: boolean;
  }>(),
  {
    title: "Confirm action",
    message: "Are you sure you want to proceed?",
    confirmText: "Confirm",
    cancelText: "Cancel",
    loadingText: "Working...",
    tone: "default",
    loading: false,
    disableClose: false,
    hint: true,
  }
);

const emit = defineEmits<{
  (e: "update:modelValue", v: boolean): void;
  (e: "confirm"): void;
  (e: "cancel"): void;
}>();

const open = computed(() => props.modelValue);

// const panelRef = ref<HTMLElement | null>(null);
const confirmBtnRef = ref<HTMLButtonElement | null>(null);

const palette = computed(() => {
  switch (props.tone) {
    case "danger":
      return {
        glow: "from-danger-500/35 to-transparent",
        ring: "ring-danger-500/30",
        seal: "bg-danger-500/80",
        confirm:
          "bg-danger-500/90 hover:bg-danger-500 text-white shadow-[0_10px_30px_-12px_rgba(239,68,68,.8)]",
      };
    case "success":
      return {
        glow: "from-emerald-500/35 to-transparent",
        ring: "ring-emerald-500/30",
        seal: "bg-emerald-500/80",
        confirm:
          "bg-emerald-500/90 hover:bg-emerald-500 text-white shadow-[0_10px_30px_-12px_rgba(16,185,129,.8)]",
      };
    default:
      return {
        glow: "from-candle-400/35 to-transparent",
        ring: "ring-candle-400/30",
        seal: "bg-candle-400/90",
        confirm:
          "bg-candle-400/90 hover:bg-candle-400 text-ink-950 shadow-[0_10px_30px_-12px_rgba(251,191,36,.75)]",
      };
  }
});

const close = () => {
  if (props.disableClose || props.loading) return;
  emit("update:modelValue", false);
  emit("cancel");
};

const confirm = () => {
  if (props.loading) return;
  emit("confirm");
};

const onOverlayClick = () => close();

const onKeydown = (e: KeyboardEvent) => {
  if (!open.value) return;

  if (e.key === "Escape") {
    e.preventDefault();
    close();
  }

  if (e.key === "Enter" && !e.shiftKey && !e.altKey && !e.metaKey && !e.ctrlKey) {
    e.preventDefault();
    confirm();
  }
};

watch(
  () => open.value,
  async (v) => {
    if (!v) return;
    await nextTick();
    // фокус на confirm, щоб Enter працював природно
    confirmBtnRef.value?.focus();
  }
);

onMounted(() => window.addEventListener("keydown", onKeydown));
onUnmounted(() => window.removeEventListener("keydown", onKeydown));
</script>

<template>
  <Teleport to="body">
    <Transition name="tavern-pop">
      <div
        v-if="open"
        class="fixed inset-0 z-[90] flex items-center justify-center p-4"
        role="dialog"
        aria-modal="true"
      >
        <!-- Overlay -->
        <div
          class="absolute inset-0 bg-ink-950/75 backdrop-blur-[2px]"
          @click="onOverlayClick"
        />

        <!-- Panel -->
        <div ref="panelRef" class="relative w-full max-w-lg">
          <!-- Outer wood frame -->
          <div
            class="rounded-3xl border border-wood-700/55 bg-wood-900/35 p-2
                   shadow-[0_35px_110px_-40px_rgba(0,0,0,.9),inset_0_1px_0_rgba(255,255,255,.06)]"
          >
            <!-- Inner parchment -->
            <div
              class="relative overflow-hidden rounded-[22px] border border-wood-700/35
                     bg-parchment-50/90 p-6 text-ink-900 ring-1"
              :class="palette.ring"
            >
              <!-- Candle glow -->
              <div
                class="pointer-events-none absolute -top-16 -right-16 h-64 w-64 rounded-full
                       bg-gradient-to-br blur-[2px] animate-[pulse_2.6s_ease-in-out_infinite]"
                :class="palette.glow"
              />
              <div
                class="pointer-events-none absolute -bottom-24 -left-24 h-64 w-64 rounded-full
                       bg-gradient-to-tr from-wood-800/25 to-transparent blur-[1px]"
              />

              <!-- Parchment noise texture (no images) -->
              <div
                class="pointer-events-none absolute inset-0 opacity-[0.07]"
                style="
                  background-image:
                    radial-gradient(circle at 18% 12%, rgba(0,0,0,.9) 0.6px, transparent 0.7px),
                    radial-gradient(circle at 84% 28%, rgba(0,0,0,.9) 0.6px, transparent 0.7px),
                    radial-gradient(circle at 30% 80%, rgba(0,0,0,.9) 0.6px, transparent 0.7px);
                  background-size: 22px 22px;
                "
              />

              <!-- Header row -->
              <div class="flex items-start gap-4">
                <!-- Wax seal icon -->
                <div class="shrink-0">
                  <div
                    class="relative h-11 w-11 rounded-2xl border border-wood-700/30 bg-wood-900/10 grid place-items-center"
                  >
                    <div
                      class="h-7 w-7 rounded-full border border-ink-950/25 shadow-[inset_0_1px_0_rgba(255,255,255,.25)]"
                      :class="palette.seal"
                    />
                    <div
                      class="pointer-events-none absolute inset-0 rounded-2xl"
                      style="box-shadow: inset 0 0 0 1px rgba(0,0,0,.06);"
                    />
                  </div>
                </div>

                <!-- Text -->
                <div class="min-w-0">
                  <div class="font-display text-xl leading-tight">
                    {{ title }}
                  </div>
                  <div class="mt-2 text-sm leading-relaxed text-ink-700/85">
                    {{ message }}
                  </div>
                </div>
              </div>

              <!-- Divider -->
              <div class="mt-5 h-px w-full bg-wood-700/20" />

              <!-- Actions -->
              <div class="mt-5 flex flex-col-reverse gap-3 sm:flex-row sm:justify-end">
                <button
                  type="button"
                  class="w-full sm:w-auto rounded-2xl border border-wood-700/35 bg-wood-900/5
                         px-5 py-3 text-sm text-ink-900 transition
                         hover:bg-wood-900/10 active:translate-y-[1px]
                         disabled:opacity-60 disabled:active:translate-y-0"
                  :disabled="loading"
                  @click="close"
                >
                  {{ cancelText }}
                </button>

                <button
                  ref="confirmBtnRef"
                  type="button"
                  class="w-full sm:w-auto rounded-2xl px-5 py-3 text-sm font-semibold transition
                         active:translate-y-[1px] disabled:opacity-60 disabled:active:translate-y-0"
                  :class="palette.confirm"
                  :disabled="loading"
                  @click="confirm"
                >
                  <span v-if="!loading">{{ confirmText }}</span>
                  <span v-else class="inline-flex items-center gap-2">
                    <span
                      class="h-4 w-4 animate-spin rounded-full border-2 border-white/40 border-t-white"
                    />
                    {{ loadingText }}
                  </span>
                </button>
              </div>

              <div v-if="hint" class="mt-4 text-xs text-ink-700/60">
                <span class="font-semibold">Enter</span> — confirm,
                <span class="font-semibold">Esc</span> — close
              </div>
            </div>
          </div>

          <!-- Small “metal studs” outside -->
          <div class="pointer-events-none absolute -left-2 top-8 h-3 w-3 rounded-full bg-ink-950/35" />
          <div class="pointer-events-none absolute -right-2 top-12 h-3 w-3 rounded-full bg-ink-950/35" />
          <div class="pointer-events-none absolute -left-2 bottom-12 h-3 w-3 rounded-full bg-ink-950/35" />
          <div class="pointer-events-none absolute -right-2 bottom-8 h-3 w-3 rounded-full bg-ink-950/35" />
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.tavern-pop-enter-active {
  transition: opacity 170ms ease, transform 200ms ease;
}
.tavern-pop-leave-active {
  transition: opacity 140ms ease, transform 160ms ease;
}
.tavern-pop-enter-from,
.tavern-pop-leave-to {
  opacity: 0;
  transform: translateY(10px) scale(0.98);
}
</style>
