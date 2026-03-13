<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  variant?: 'primary' | 'danger' | 'ghost'
  block?: boolean
  type?: 'button' | 'submit' | 'reset'
  disabled?: boolean
}>(), {
  variant: 'primary',
  block: false,
  type: 'button',
  disabled: false,
})

const cls = computed(() => {
  const base =
    'inline-flex items-center justify-center gap-2 px-4 py-2 rounded-xl font-semibold select-none transition active:translate-y-[1px]'
  const block = props.block ? ' w-full' : ''
  const v =
    props.variant === 'primary'
      ? ' bg-candle-400 text-ink-900 hover:bg-candle-300'
      : props.variant === 'danger'
        ? ' bg-danger-500 text-parchment-50 hover:bg-danger-600'
        : ' bg-parchment-50/60\n' +
        '        text-ink-900/70\n' +
        '        hover:text-ink-900\n' +
        '        hover:bg-parchment-50/80\n' +
        '        active:bg-parchment-50\n' +
        '        border\n' +
        '        border-wood-700/20'
  return base + v + block
})
</script>

<template>
  <button :type="type" :class="cls" :disabled="disabled">
    <slot />
  </button>
</template>
