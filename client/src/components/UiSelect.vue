<script setup lang="ts">
import { computed, ref } from "vue"
import { onClickOutside } from "@vueuse/core"

export type Option = {
  label: string
  value: string | number
  type?: string
}

const props = defineProps<{
  options: Option[]
  placeholder?: string
}>()

// v-model для компонента
const model = defineModel<string | number | null>({ default: null })

const isOpen = ref(false)
const root = ref<HTMLElement | null>(null)

const isSame = (a: string | number | null, b: string | number) =>
  a !== null && String(a) === String(b)

const selected = computed(() =>
  props.options.find(o => isSame(model.value, o.value))
)

const selectOption = (option: Option) => {
  model.value = option.value
  isOpen.value = false
}

onClickOutside(root, () => {
  isOpen.value = false
})
</script>

<template>
  <div ref="root" class="relative w-full">
    <button
      type="button"
      @click.stop="isOpen = !isOpen"
      class="w-full rounded-xl border border-wood-700/35 bg-parchment-50/70 px-4 py-2 text-base text-left flex items-center justify-between transition hover:border-candle-400/60"
    >
      <slot v-if="selected" name="option" :option="selected">
        {{ selected.label }}
      </slot>
      <span v-else class="truncate">Select an option</span>
      <span class="transition-transform duration-200" :class="{ 'rotate-180': isOpen }">
        ▾
      </span>
    </button>

    <transition
      enter-active-class="transition duration-150 ease-out"
      enter-from-class="opacity-0 -translate-y-1"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition duration-100 ease-in"
      leave-from-class="opacity-100 translate-y-0"
      leave-to-class="opacity-0 -translate-y-1"
    >
      <div
        v-if="isOpen"
        class="absolute z-50 mt-2 w-full rounded-xl border border-wood-700/35 bg-parchment-50 shadow-xl backdrop-blur"
        @click.stop
      >
        <button
          v-for="option in options"
          :key="String(option.value)"
          type="button"
          class="w-full text-left px-4 py-2 hover:bg-wood-700/10 text-base cursor-pointer flex items-center gap-2 transition"
          @mousedown.prevent
          @click.prevent="selectOption(option)"
        >
          <slot name="option" :option="option">
            {{ option.label }}
          </slot>
        </button>
      </div>
    </transition>
  </div>
</template>
