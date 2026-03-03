<script setup lang="ts">
import { computed } from "vue"

export type CurrencyType = "bronze" | "silver" | "gold"

const props = defineProps<{
  type: CurrencyType
  size?: number
}>()

const size = computed(() => props.size ?? 20)

const color = computed(() => {
  switch (props.type) {
    case "gold":
      return "#D4A437"
    case "silver":
      return "#C0C0C0"
    case "bronze":
      return "#B87333"
    default:
      return "#999"
  }
})
</script>

<template>
  <svg
    :width="size"
    :height="size"
    viewBox="0 0 24 24"
    fill="none"
  >
    <circle
      cx="12"
      cy="12"
      r="9"
      :fill="color"
      stroke="#000"
      stroke-opacity="0.2"
      stroke-width="1"
    />
    <circle
      cx="9"
      cy="9"
      r="3"
      fill="white"
      opacity="0.15"
    />
  </svg>
</template>
