<script setup lang="ts">
import { ref } from 'vue'
import TavernShell from '@/components/wrappers/TavernShell.vue'
import UiButton from '@/components/form/UiButton.vue'
import {useAuthStore} from "@/stores/auth.ts";
import ShowPasswordIcon from "@/components/icons/ShowPasswordIcon.vue";

const auth = useAuthStore()

const username = ref('')
const password = ref('')
const loading = ref(false)
const confirmPassword = ref('')
const showPassword = ref(false)
const showConfirmPassword = ref(false)

const localError = ref('')

const submit = async () => {
  localError.value = ''

  if (password.value !== confirmPassword.value) {
    localError.value = 'Passwords do not match.'
    return
  }

  loading.value = true

  auth.register(username.value.trim(), password.value)
    .catch((err) => {
      localError.value = err.response.data.message
    })
    .finally(() => loading.value = false)
}
</script>

<template>
  <TavernShell>
    <div class="w-full flex justify-center">
      <div class="panel w-full max-w-2xl">
        <div class="font-display text-xl text-center">Register</div>
        <div class="mt-1 text-sm text-ink-900/70 text-center">
          Create an account to track wins, rating and leaderboard stats.
        </div>

            <div
              v-if="localError"
              class="mt-4 rounded-xl border border-danger-500/40 bg-danger-500/10 p-3 text-sm"
            >
              {{ localError }}
            </div>

        <form @submit.prevent="submit" class="flex flex-col gap-2 mt-3">
          <div class="text-md">
            <label class="mt-4 block font-semibold">Username</label>
            <input
              v-model="username"
              autocomplete="username"
              class="mt-1 w-full rounded-xl border border-wood-700/35 bg-parchment-50 px-4 py-2 text-ink-900 outline-none focus:ring-2 focus:ring-candle-400/60"
              placeholder="e.g. johndoe"
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
          <div class="text-md">
            <label class="mt-3 block font-semibold">Confirm password</label>
            <div
              class="mt-1 flex items-center rounded-xl border border-wood-700/35 bg-parchment-50 focus-within:ring-2 focus-within:ring-candle-400/60"
            >
              <input
                :type="showConfirmPassword ? 'text' : 'password'"
                v-model="confirmPassword"
                placeholder="••••••••"
                required
                class="w-full bg-transparent px-4 py-2 text-ink-900 outline-none"
              />

              <button
                type="button"
                @click="showConfirmPassword = !showConfirmPassword"
                class="flex shrink-0 items-center justify-center px-4 text-ink-900/55 transition hover:text-ink-900"
              >
                <ShowPasswordIcon :show-password="showConfirmPassword" />
              </button>
            </div>
          </div>
          <div class="mt-4 flex gap-2 sm:flex-row sm:items-center sm:justify-between">
            <div class="flex gap-2">
              <UiButton
                type="submit"
                variant="primary"
                :disabled="loading"
              >
                {{ 'Create account' }}
              </UiButton>

              <RouterLink to="/">
                <UiButton variant="ghost">Back</UiButton>
              </RouterLink>
            </div>

            <RouterLink to="/login" class="text-sm text-dark-candle-300 hover:text-dark-candle-200">
              Already have an account? Login →
            </RouterLink>
          </div>
        </form>
      </div>
    </div>
  </TavernShell>
</template>
