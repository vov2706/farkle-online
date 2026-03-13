<script setup lang="ts">
defineProps<{
  title: string
  description: string
  to: string
  badge?: string
}>()
</script>

<template>
  <RouterLink
    :to="to"
    class="group block rounded-[var(--radius-tavern)] border border-wood-700/35 bg-tavern-900/55 p-5
           transition hover:-translate-y-[1px] hover:border-candle-400/45 hover:bg-tavern-900/70"
  >
    <div class="flex items-start justify-between gap-4">
      <div>
        <div class="font-display text-lg text-parchment-50">
          {{ title }}
        </div>
        <div class="mt-1 text-parchment-50/70">
          {{ description }}
        </div>
      </div>

      <span v-if="badge" class="chip whitespace-nowrap">
        {{ badge }}
      </span>
    </div>

    <div class="mt-4 text-candle-200/90 group-hover:text-candle-200">
      Open →
    </div>
  </RouterLink>
</template>
