<script setup lang="ts">
import { useRouter } from 'vue-router'
import TavernShell from '@/components/TavernShell.vue'
import UiButton from '@/components/UiButton.vue'

const router = useRouter()

const goHome = () => {
  router.push('/')
}
</script>

<template>
  <TavernShell :with-header="false" :with-footer="false">
    <div class="relative flex min-h-[80vh] flex-col items-center justify-center text-center px-6">
      <!-- Background glow -->
      <div class="absolute inset-0 bg-[radial-gradient(circle_at_center,_rgba(255,200,120,0.08),_transparent_70%)]" />

      <!-- 404 -->
      <div class="relative z-10">

        <div class="font-display text-[120px] leading-none text-amber-400 drop-shadow-[0_0_25px_rgba(255,180,80,0.5)]">
          404
        </div>

        <div class="mt-4 text-2xl font-display text-parchment-200">
          This path leads nowhere, traveler.
        </div>

        <p class="mt-4 max-w-xl text-sm text-parchment-300/70">
          The tavern keeper swears this road existed yesterday...
          Perhaps too much ale was involved.
        </p>

        <!-- Dice decoration -->
        <div class="mt-10 flex justify-center gap-4 text-4xl opacity-70">
          🎲 🎲 🎲
        </div>

        <div class="mt-10">
          <UiButton @click="goHome" size="lg">
            Return to the Tavern
          </UiButton>
        </div>

      </div>
    </div>
  </TavernShell>
</template>

<style scoped>
.font-display {
  font-family: 'Cinzel', serif;
}
</style>
