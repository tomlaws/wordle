<script lang="ts">
	import type { HTMLInputAttributes } from 'svelte/elements';

	let { children, value = $bindable(), class: inputClass, ...rest }: HTMLInputAttributes = $props();
</script>

<input
    {...rest}
    bind:value
    class={`rounded-md border border-[var(--color-input-border)] dark:border-[var(--color-input-border)] px-3 py-2 text-base focus:outline-none focus:ring-2 focus:ring-gray-500 dark:focus:ring-gray-400 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 transition-colors duration-200 ${inputClass}`}
    />