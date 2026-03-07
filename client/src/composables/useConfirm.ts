import { ref } from "vue";

type ConfirmTone = "danger" | "default" | "success";

type ConfirmOptions = {
  title?: string;
  message?: string;
  confirmText?: string;
  cancelText?: string;
  loadingText?: string;
  tone?: ConfirmTone;
};

const isOpen = ref(false);
const loading = ref(false);
const options = ref<ConfirmOptions>({
  title: "Are you sure?",
  message: "This action cannot be undone.",
  confirmText: "Confirm",
  cancelText: "Cancel",
  loadingText: "Working...",
  tone: "default",
});

let resolver: ((v: boolean) => void) | null = null;

export function useConfirm() {
  const open = (opts: ConfirmOptions = {}) => {
    options.value = {
      title: "Are you sure?",
      message: "This action cannot be undone.",
      confirmText: "Confirm",
      cancelText: "Cancel",
      loadingText: "Working...",
      tone: "default",
      ...opts,
    };

    isOpen.value = true;
    loading.value = false;

    return new Promise<boolean>((resolve) => {
      resolver = resolve;
    });
  };

  const onConfirm = async (cb?: () => Promise<void> | void) => {
    try {
      loading.value = true;
      await cb?.();
      isOpen.value = false;
      resolver?.(true);
    } finally {
      loading.value = false;
      resolver = null;
    }
  };

  const onCancel = () => {
    isOpen.value = false;
    resolver?.(false);
    resolver = null;
  };

  return { isOpen, loading, options, open, onConfirm, onCancel };
}
