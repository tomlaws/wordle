<script lang="ts">
	import { GAME_KEY, type GameContext } from '$lib/context/game-context';
	import { PlayAgainPayload, type GameOverPayload } from '$lib/types/payload';
	export { GameOverPayload } from '$lib/types/payload';
	import { getContext } from 'svelte';

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

<h2>Game Over</h2>
<p>The word was {matchInfo!.gameOver!.answer}.</p>
{#if matchInfo!.gameOver!.winner}
	{#if matchInfo!.gameOver!.winner.id === playerInfo.id}
		<p>Congratulations, you won!</p>
	{:else}
		<p>{matchInfo!.gameOver!.winner.nickname} won the game.</p>
	{/if}
{:else}
	<p>The game ended in a draw.</p>
{/if}
<button onclick={() => playAgain(true)}>Play Again</button>
<button onclick={() => playAgain(false)}>Quit</button>
