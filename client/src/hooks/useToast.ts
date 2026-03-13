import { reactive } from 'vue'

export type Toast = {
  id: string
  title?: string
  message: string
  kind?: 'info' | 'success' | 'error'
  timeoutMs?: number
}

type ToastState = {
  items: Toast[]
}

const state = reactive<ToastState>({
  items: [],
})

const uid = () =>
  (globalThis.crypto?.randomUUID?.() ?? `${Date.now()}_${Math.random().toString(16).slice(2)}`)

export function useToast() {
  const push = (t: Omit<Toast, 'id'>) => {
    const id = uid()
    const toast: Toast = {
      id,
      kind: t.kind ?? 'info',
      timeoutMs: t.timeoutMs ?? 3000,
      title: t.title,
      message: t.message,
    }

    state.items.push(toast)

    if (toast.timeoutMs && toast.timeoutMs > 0) {
      window.setTimeout(() => dismiss(id), toast.timeoutMs)
    }

    return id
  }

  const dismiss = (id: string) => {
    const idx = state.items.findIndex((x) => x.id === id)
    if (idx >= 0) state.items.splice(idx, 1)
  }

  return {
    toasts: state.items,
    push,
    dismiss,
  }
}
