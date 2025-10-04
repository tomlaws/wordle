<script lang="ts">
	import { GAME_KEY, type GameContext } from '$lib/context/game-context';
	import { getContext } from "svelte";
	import Loader from './Loader.svelte';

	const { playerInfo } = getContext<GameContext>(GAME_KEY);
</script>

<div class="flex flex-col items-center justify-center min-h-screen text-[1.2rem] text-[#333]">
    <h1 class="mb-4 text-xl font-bold text-center text-gray-700">Just a moment, {playerInfo.nickname}!</h1>
    <p class="text-center text-gray-600 dark:text-gray-300 text-sm font-bold mb-6">We're finding opponents...</p>
    <Loader />
</div>

<style>
</style>
