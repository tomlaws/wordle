<script lang="ts">
	import { TOAST_KEY, type ToastAPI } from '$lib/context/toast-context';
  import { setContext } from 'svelte';
  import { toast } from 'svoast';

  const toastAPI: ToastAPI = {
    success: (msg) => toast.success(msg),
    error: (msg) => toast.error(msg),
    info: (msg) => toast.info(msg),
    warning: (msg) => toast.warning(msg)
  };

  setContext<ToastAPI>(TOAST_KEY, toastAPI);
</script>

<slot />
