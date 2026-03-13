<script setup lang="ts">
  defineProps<{
    showPassword: boolean
  }>()
</script>

<template>
  <svg
    v-if="!showPassword"
    xmlns="http://www.w3.org/2000/svg"
    fill="none"
    viewBox="0 0 24 24"
    stroke="currentColor"
    class="h-5 w-5"
  >
    <path
      stroke-linecap="round"
      stroke-linejoin="round"
      stroke-width="1.8"
      d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7Z"
    />
    <path
      stroke-linecap="round"
      stroke-linejoin="round"
      stroke-width="1.8"
      d="M12 15.5A3.5 3.5 0 1 0 12 8.5a3.5 3.5 0 0 0 0 7Z"
    />
  </svg>

  <svg
    v-else
    xmlns="http://www.w3.org/2000/svg"
    fill="none"
    viewBox="0 0 24 24"
    stroke="currentColor"
    class="h-5 w-5"
  >
    <path
      stroke-linecap="round"
      stroke-linejoin="round"
      stroke-width="1.8"
      d="M3 3l18 18"
    />
    <path
      stroke-linecap="round"
      stroke-linejoin="round"
      stroke-width="1.8"
      d="M10.58 10.58A2 2 0 0 0 13.42 13.42"
    />
    <path
      stroke-linecap="round"
      stroke-linejoin="round"
      stroke-width="1.8"
      d="M9.88 5.09A10.94 10.94 0 0 1 12 5c4.48 0 8.27 2.94 9.54 7a10.96 10.96 0 0 1-4.16 5.17"
    />
    <path
      stroke-linecap="round"
      stroke-linejoin="round"
      stroke-width="1.8"
      d="M6.23 6.23A10.95 10.95 0 0 0 2.46 12c1.27 4.06 5.06 7 9.54 7 1.61 0 3.14-.38 4.5-1.05"
    />
  </svg>
</template>

<style scoped>

</style>
