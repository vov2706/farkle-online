<script setup lang="ts">
import ToastsHost from "./components/ToastsHost.vue";
import { useConfirm } from "@/composables/useConfirm";
import UiConfirmModal from "@/components/UiConfirmModal.vue";

const confirm = useConfirm();
</script>

<template>
  <div class="app-root">
    <RouterView />
    <ToastsHost />
    <UiConfirmModal
      v-model="confirm.isOpen.value"
      :title="confirm.options.value.title"
      :message="confirm.options.value.message"
      :confirm-text="confirm.options.value.confirmText"
      :cancel-text="confirm.options.value.cancelText"
      :tone="confirm.options.value.tone"
      :loading="confirm.loading.value"
      @confirm="confirm.onConfirm()"
      @cancel="confirm.onCancel"
    />
  </div>
</template>

<style scoped>
.logo {
  height: 6em;
  padding: 1.5em;
  will-change: filter;
  transition: filter 300ms;
}
.logo:hover {
  filter: drop-shadow(0 0 2em #646cffaa);
}
.logo.vue:hover {
  filter: drop-shadow(0 0 2em #42b883aa);
}
.app-root {
  min-height: 100dvh;
}
</style>
