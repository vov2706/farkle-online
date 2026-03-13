<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'

const props = withDefaults(
  defineProps<{
    /** Показувати/ховати оверлей */
    show: boolean
    /** UTC ms timestamp, коли гра має стартувати (start_at) */
    startAt?: number | null
    /** Резервний варіант: просто показати 3..2..1 без startAt */
    seconds?: number
    /** Текст зверху */
    title?: string
    /** Дозволити закривати (на майбутнє під "Cancel") */
    closable?: boolean
  }>(),
  {
    startAt: null,
    seconds: 3,
    title: 'Game will start in',
    closable: false,
  }
)

const emit = defineEmits<{
  (e: 'done'): void
  (e: 'close'): void
}>()

const remaining = ref<number>(props.seconds)
let raf = 0
let interval: number | null = null

const nowMs = () => Date.now()

const tick = () => {
  if (!props.show) return

  if (props.startAt) {
    const msLeft = props.startAt - nowMs()
    const sLeft = Math.ceil(msLeft / 1000)
    remaining.value = Math.max(0, sLeft)
  } else {
    // fallback: просто відраховуємо секундоміром
    // (interval робить це)
  }

  if (remaining.value <= 0) {
    stop()
    emit('done')
    return
  }

  raf = requestAnimationFrame(tick)
}

const start = () => {
  stop()

  remaining.value = props.startAt
    ? Math.max(0, Math.ceil((props.startAt - nowMs()) / 1000))
    : props.seconds

  // якщо без startAt — секундомір
  if (!props.startAt) {
    interval = window.setInterval(() => {
      remaining.value = Math.max(0, remaining.value - 1)
      if (remaining.value <= 0) {
        stop()
        emit('done')
      }
    }, 1000)
  }

  raf = requestAnimationFrame(tick)
}

const stop = () => {
  if (raf) cancelAnimationFrame(raf)
  raf = 0
  if (interval) window.clearInterval(interval)
  interval = null
}

watch(
  () => [props.show, props.startAt] as const,
  ([show]) => {
    if (show) start()
    else stop()
  },
  { immediate: true }
)

onBeforeUnmount(stop)

const bigNumber = computed(() => {
  // красиво: 3 2 1 і потім "GO"
  if (remaining.value <= 0) return 'GO'
  return String(remaining.value)
})

const subtitle = computed(() => {
  // дрібний текст під цифрою
  if (remaining.value <= 0) return 'Good luck!'
  return 'Get ready…'
})

const onKeyDown = (e: KeyboardEvent) => {
  if (!props.show) return
  if (!props.closable) return
  if (e.key === 'Escape') emit('close')
}

onMounted(() => window.addEventListener('keydown', onKeyDown))
onBeforeUnmount(() => window.removeEventListener('keydown', onKeyDown))
</script>

<template>
  <Transition name="td-fade">
    <div v-if="show" class="td-overlay" role="dialog" aria-modal="true">
      <div class="td-backdrop" />

      <div class="td-card">
        <div class="td-top">
          <div class="td-title">{{ title }}</div>

          <button
            v-if="closable"
            class="td-x"
            type="button"
            aria-label="Close"
            @click="$emit('close')"
          >
            ✕
          </button>
        </div>

        <div class="td-number-wrap" aria-live="polite">
          <div class="td-glow" />
          <div class="td-number">{{ bigNumber }}</div>
        </div>

        <div class="td-sub">{{ subtitle }}</div>

        <div class="td-hint">
          <span class="td-dot" />
          <span>Starting…</span>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
/* overlay */
.td-overlay {
  position: fixed;
  inset: 0;
  z-index: 80;
  display: grid;
  place-items: center;
  padding: 24px;
}

.td-backdrop {
  position: absolute;
  inset: 0;
  background:
    radial-gradient(1200px 600px at 50% 35%, rgba(255, 220, 160, 0.14), transparent 65%),
    rgba(10, 7, 5, 0.78);
  backdrop-filter: blur(6px);
}

.td-card {
  position: relative;
  width: min(560px, 92vw);
  border-radius: 22px;
  border: 1px solid rgba(245, 235, 220, 0.18);
  background:
    radial-gradient(900px 420px at 50% 10%, rgba(255, 210, 150, 0.14), transparent 65%),
    linear-gradient(180deg, rgba(35, 24, 16, 0.92), rgba(18, 12, 8, 0.92));
  box-shadow:
    0 24px 70px rgba(0, 0, 0, 0.55),
    inset 0 1px 0 rgba(255, 255, 255, 0.08);
  padding: 22px 22px 18px;
  text-align: center;
}

.td-top {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  margin-bottom: 8px;
}

.td-title {
  font-family: ui-serif, Georgia, "Times New Roman", serif;
  font-size: 18px;
  letter-spacing: 0.04em;
  color: rgba(245, 235, 220, 0.92);
}

.td-x {
  position: absolute;
  top: 14px;
  right: 14px;
  height: 36px;
  width: 36px;
  border-radius: 12px;
  border: 1px solid rgba(245, 235, 220, 0.16);
  background: rgba(0, 0, 0, 0.18);
  color: rgba(245, 235, 220, 0.9);
  cursor: pointer;
  transition: transform 0.12s ease, background 0.12s ease, border-color 0.12s ease;
}
.td-x:hover {
  transform: translateY(-1px);
  background: rgba(0, 0, 0, 0.28);
  border-color: rgba(255, 210, 150, 0.22);
}

.td-number-wrap {
  position: relative;
  margin: 12px auto 4px;
  display: grid;
  place-items: center;
  height: 160px;
}

.td-glow {
  position: absolute;
  width: 240px;
  height: 240px;
  border-radius: 999px;
  background: radial-gradient(circle at 50% 50%, rgba(255, 190, 110, 0.25), transparent 62%);
  filter: blur(2px);
  animation: td-pulse 1.15s ease-in-out infinite;
}

.td-number {
  position: relative;
  font-family: ui-serif, Georgia, "Times New Roman", serif;
  font-size: 92px;
  line-height: 1;
  letter-spacing: 0.02em;
  color: rgba(255, 240, 220, 0.98);
  text-shadow:
    0 2px 0 rgba(0, 0, 0, 0.35),
    0 12px 32px rgba(0, 0, 0, 0.45);
  transform-origin: 50% 60%;
  animation: td-pop 1s ease-in-out infinite;
}

.td-sub {
  margin-top: 4px;
  font-size: 13px;
  color: rgba(245, 235, 220, 0.72);
}

.td-hint {
  margin-top: 14px;
  display: inline-flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  border-radius: 999px;
  border: 1px solid rgba(245, 235, 220, 0.12);
  background: rgba(0, 0, 0, 0.16);
  color: rgba(245, 235, 220, 0.82);
  font-size: 12px;
}

.td-dot {
  width: 8px;
  height: 8px;
  border-radius: 999px;
  background: rgba(255, 190, 110, 0.9);
  box-shadow: 0 0 0 3px rgba(255, 190, 110, 0.15);
  animation: td-dot 1.1s ease-in-out infinite;
}

/* animations */
@keyframes td-pulse {
  0%, 100% { transform: scale(0.94); opacity: 0.75; }
  50% { transform: scale(1.02); opacity: 1; }
}

@keyframes td-pop {
  0%, 100% { transform: translateY(0) scale(1); }
  35% { transform: translateY(-2px) scale(1.02); }
  70% { transform: translateY(0) scale(1); }
}

@keyframes td-dot {
  0%, 100% { transform: scale(0.9); opacity: 0.65; }
  50% { transform: scale(1.15); opacity: 1; }
}

/* transition */
.td-fade-enter-active,
.td-fade-leave-active {
  transition: opacity 0.18s ease;
}
.td-fade-enter-from,
.td-fade-leave-to {
  opacity: 0;
}
</style>
