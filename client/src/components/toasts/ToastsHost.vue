<script setup lang="ts">
import { useToast } from '@/hooks/useToast'

const { toasts, dismiss } = useToast()

const kindClass = (kind?: 'info' | 'success' | 'error') => {
  if (kind === 'success') return 'border-candle-400/40 bg-tavern-900/70'
  if (kind === 'error') return 'border-danger-500/45 bg-tavern-900/70'
  return 'border-wood-700/35 bg-tavern-900/70'
}
</script>

<template>
  <div class="fixed right-4 top-4 z-50 space-y-2">
    <div
      v-for="t in toasts"
      :key="t.id"
      class="w-[320px] rounded-xl border p-3 shadow-[0_12px_30px_rgba(0,0,0,.45)] backdrop-blur"
      :class="kindClass(t.kind)"
    >
      <div class="flex items-start justify-between gap-3">
        <div class="min-w-0">
          <div v-if="t.title" class="font-semibold text-parchment-50">
            {{ t.title }}
          </div>
          <div class="text-sm text-parchment-50/75">
            {{ t.message }}
          </div>
        </div>

        <button
          class="text-parchment-50/60 hover:text-parchment-50"
          @click="dismiss(t.id)"
          aria-label="Dismiss"
        >
          ✕
        </button>
      </div>
    </div>
  </div>
</template>
