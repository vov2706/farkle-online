<script setup lang="ts">
import { ref, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import UiButton from './UiButton.vue'
import { useAuthStore } from '@/stores/auth.ts'

const props = withDefaults(defineProps<{ withHeader?: boolean; withFooter?: boolean }>(), {
  withHeader: true,
  withFooter: true,
})

const router = useRouter()
const auth = useAuthStore()
const user = auth.user
const logout = auth.logout
const balance = user?.balance

const menuOpen = ref(false)
const menuWrapRef = ref<HTMLElement | null>(null)

const closeMenu = () => (menuOpen.value = false)
const toggleMenu = () => (menuOpen.value = !menuOpen.value)

const onDocPointerDown = (e: PointerEvent) => {
  if (!menuOpen.value) return
  const el = menuWrapRef.value
  if (!el) return
  if (e.target instanceof Node && !el.contains(e.target)) closeMenu()
}

const onDocKeyDown = (e: KeyboardEvent) => {
  if (!menuOpen.value) return
  if (e.key === 'Escape') closeMenu()
}

document.addEventListener('pointerdown', onDocPointerDown)
document.addEventListener('keydown', onDocKeyDown)

onBeforeUnmount(() => {
  document.removeEventListener('pointerdown', onDocPointerDown)
  document.removeEventListener('keydown', onDocKeyDown)
})

const goProfile = async () => {
  closeMenu()
  await router.push('/profile')
}

const doLogout = () => {
  closeMenu()
  logout()
}

</script>

<template>
  <div
    class="min-h-dvh w-full
           flex flex-col justify-between
           px-3 py-4 sm:px-6 sm:py-6
           bg-[radial-gradient(circle_at_50%_20%,rgba(212,145,57,.22),transparent_60%)]
           "
  >
    <header v-if="props.withHeader" class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
      <div @click="$router.push('/')" class="cursor-pointer">
        <div class="text-candle-200 font-display tracking-[0.06em] text-2xl">
          Farkle Online
        </div>
        <div class="text-parchment-50/70">
          Multiplayer dice in a tavern mode.
        </div>
      </div>

      <div class="flex items-center gap-2 flex-wrap justify-start sm:justify-end text-xl">
        <div v-if="auth.isLoggedIn" ref="menuWrapRef" class="relative">
          <button
            type="button"
            class="chip flex items-center gap-2 select-none hover:brightness-110 transition"
            @click="toggleMenu"
            :aria-expanded="menuOpen"
            aria-haspopup="menu"
          >
            <span>👤 {{ user?.username }}</span>
            <span class="mx-1 h-4 w-px bg-parchment-50/20"></span>
            <span class="ml-1 text-parchment-50/60">▾</span>
          </button>

          <!-- Popover -->
          <!-- Popover -->
          <transition name="fade-pop">
            <div
              v-if="menuOpen"
              class="
                absolute right-0 mt-2 w-48
                z-[9999]
                rounded-2xl
                border border-parchment-50/20
                bg-[#0b0906]/95
                backdrop-blur-md
                shadow-[0_18px_60px_rgba(0,0,0,.55)]
                ring-1 ring-black/30
                overflow-hidden
              "
              role="menu"
            >
              <button
                type="button"
                class="w-full text-left px-3 py-2.5 text-parchment-50/95 hover:bg-parchment-50/10 transition"
                role="menuitem"
                @click="goProfile"
              >
                Profile
              </button>

              <div class="h-px bg-parchment-50/10"></div>

              <button
                type="button"
                class="w-full text-left px-3 py-2.5 text-parchment-50/95 hover:bg-parchment-50/10 transition"
                role="menuitem"
                @click="doLogout"
              >
                Logout
              </button>
            </div>
          </transition>
        </div>

        <!-- GUEST -->
        <RouterLink v-if="!user" to="/login">
          <UiButton variant="ghost">Login</UiButton>
        </RouterLink>

        <RouterLink v-if="!user" to="/register">
          <UiButton variant="primary">Register</UiButton>
        </RouterLink>
      </div>
    </header>

    <main class="mt-6">
      <slot />
    </main>

    <footer v-if="props.withFooter" class="mt-8 text-parchment-50/50 text-sm text-center">
      <span>© {{ new Date().getFullYear() }} Farkle online. All Rights Reserved</span>
    </footer>
  </div>
</template>

<style scoped>
.fade-pop-enter-active,
.fade-pop-leave-active {
  transition: opacity 120ms ease, transform 120ms ease;
}
.fade-pop-enter-from,
.fade-pop-leave-to {
  opacity: 0;
  transform: translateY(-4px) scale(0.98);
}
</style>
