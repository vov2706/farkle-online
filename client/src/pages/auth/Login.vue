<script setup lang="ts">
import {ref} from 'vue'
import TavernShell from '@/components/wrappers/TavernShell.vue'
import UiButton from '@/components/form/UiButton.vue'
import {useAuthStore} from "@/stores/auth.ts";
import ShowPasswordIcon from "@/components/icons/ShowPasswordIcon.vue";

const username = ref('')
const password = ref('')
const error = ref('')
const loading  = ref(false)
const showPassword = ref(false)

const submit = async () => {
  const auth = useAuthStore()

  loading.value = true

  auth.login(username.value, password.value)
    .catch((err) => {
      error.value = err.response.data.message
    })
    .finally(() => loading.value = false)

}
</script>

<template>
  <TavernShell>
    <div class="w-full flex justify-center">
      <div class="panel w-full max-w-2xl">
        <div class="font-display text-xl text-center">Login</div>

        <div v-if="error" class="mt-3 rounded-xl border border-danger-500/40 bg-danger-500/10 p-3 text-sm">
          {{ error }}
        </div>

        <form @submit.prevent="submit" class="flex flex-col gap-2 mt-3">
          <div class="text-md">
            <label class="mt-4 block font-semibold">Username</label>
            <input
              class="mt-1 text-md w-full rounded-xl border border-wood-700/35 bg-parchment-50 px-4 py-2 text-ink-900 outline-nonefocus:ring-2 focus:ring-candle-400/60"
              v-model="username"
              required
            />
          </div>
          <div class="text-md">
            <label class="mt-3 block font-semibold">Password</label>

            <div
              class="mt-1 flex items-center rounded-xl border border-wood-700/35 bg-parchment-50 focus-within:ring-2 focus-within:ring-candle-400/60"
            >
              <input
                :type="showPassword ? 'text' : 'password'"
                v-model="password"
                placeholder="••••••••"
                required
                class="w-full bg-transparent px-4 py-2 text-ink-900 outline-none"
              />

              <button
                type="button"
                @click="showPassword = !showPassword"
                class="flex shrink-0 items-center justify-center px-4 text-ink-900/55 transition hover:text-ink-900"
              >
                <ShowPasswordIcon :show-password="showPassword" />
              </button>
            </div>
          </div>
          <div class="mt-4 flex gap-2 sm:flex-row sm:items-center sm:justify-between">
            <div class="flex gap-2">
              <UiButton variant="primary" type="submit" >
                {{ 'Login' }}
              </UiButton>
              <RouterLink to="/">
                <UiButton variant="ghost">Back</UiButton>
              </RouterLink>
            </div>
            <RouterLink to="/register" class="text-sm text-dark-candle-300 hover:text-dark-candle-200">
              Don't have an account? Register now →
            </RouterLink>
          </div>
        </form>
      </div>
    </div>
  </TavernShell>
</template>
