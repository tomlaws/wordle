<script lang="ts">
	import { GAME_KEY, type GameContext } from '$lib/context/game-context';
	import { PlayAgainPayload, type GameOverPayload } from '$lib/types/payload';
	export { GameOverPayload } from '$lib/types/payload';
	import { getContext } from 'svelte';
	import Button from './common/Button.svelte';

	const gameContext = getContext<GameContext>(GAME_KEY);
	const { websocket, playerInfo, matchInfo } = $derived(gameContext);

	function playAgain(confirm: boolean) {
		if (confirm) {
			const playAgainPayload = new PlayAgainPayload();
			playAgainPayload.confirm = true;
			websocket.send(playAgainPayload);
		} else {
			location.reload();
		}
	}
</script>

<h2 class="text-2xl font-bold">
{#if matchInfo!.gameOver!.winner}
	{#if matchInfo!.gameOver!.winner.id === playerInfo.id}
		<p>🎉 Congratulations, you won!</p>
	{:else}
		<p>🫢 Oops, you lost.</p>
	{/if}
{:else}
	<p>The game ended in a draw.</p>
{/if}
</h2>
<p class="text-lg mt-2">The word was {matchInfo!.gameOver!.answer.toUpperCase()}.</p>

<div class="mt-8 flex gap-4">
	<Button onclick={() => playAgain(true)}>Play Again</Button>
	<Button outline={true} onclick={() => playAgain(false)}>Quit</Button>
</div>
