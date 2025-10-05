<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import favicon from '$lib/assets/favicon.svg';
	import ToastProvider from '$lib/components/ToastProvider.svelte';
	import IconSun from '$lib/components/icons/IconSun.svelte';
	import IconMoon from '$lib/components/icons/IconMoon.svelte';

	let { children } = $props();

	let dark = $state(false);
	onMount(() => {
		const saved = localStorage.getItem('theme');
		if (saved) {
			dark = saved === 'dark';
		} else {
			dark = window.matchMedia('(prefers-color-scheme: dark)').matches;
		}
		setDarkClass(dark);
	});

	function toggleDark() {
		dark = !dark;
		localStorage.setItem('theme', dark ? 'dark' : 'light');
		setDarkClass(dark);
	}

	function setDarkClass(d: boolean) {
		document.documentElement.classList.toggle('dark', d);
	}
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<div
	class="min-h-screen bg-white text-gray-900 transition-colors duration-300 dark:bg-[#191e25] dark:text-gray-100"
>
	<header class="fixed left-0 right-0 flex items-center justify-between py-4 px-8">
		<h1 class="text-2xl font-bold uppercase">Wordle</h1>
		<button
			class="rounded-md bg-gray-200 p-2 transition-colors hover:bg-gray-300 dark:bg-gray-700 dark:hover:bg-gray-600 cursor-pointer"
			onclick={toggleDark}
			aria-label="Toggle Dark Mode"
		>
			{#if dark}
				<IconSun class="h-6 w-6" />
			{:else}
				<IconMoon class="h-6 w-6" />
			{/if}
		</button>
	</header>
	<main>
		<ToastProvider>
			{@render children?.()}
		</ToastProvider>
	</main>
</div>
