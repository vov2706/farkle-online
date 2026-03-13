<script setup lang="ts">
defineProps<{
  title: string
  description?: string
  to: string
}>()
</script>

<template>
  <RouterLink
    :to="to"
    class="group flex items-center justify-between rounded-[var(--radius-tavern)] border border-wood-700/25
           bg-tavern-900/40 px-5 py-4 transition
           hover:bg-tavern-900/55 hover:border-candle-400/35"
  >
    <div>
      <div class="font-display text-xl text-parchment-50">
        {{ title }}
      </div>
      <div v-if="description" class="mt-1 text-parchment-50/65">
        {{ description }}
      </div>
    </div>

    <div class="text-candle-200/80 group-hover:text-candle-200">
      →
    </div>
  </RouterLink>
</template>
