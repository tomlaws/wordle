<script lang="ts">
	import { type GameContext, GAME_KEY } from '$lib/context/game-context';
	import { getContext, onMount } from 'svelte';

	const gameContext = getContext<GameContext>(GAME_KEY);
	const { matchInfo } = $derived(gameContext);

	$effect(() => {
		document
			.getElementById('row-' + matchInfo!.currentRound)
			?.scrollIntoView({ behavior: 'smooth', block: 'center' });
	});
</script>

<div class="board">
	{#each matchInfo!.guesses as guess, i}
		<div class="row" id={'row-' + (i + 1)}>
			{#each guess as letter, j}
				<div
					class="box"
					class:miss={letter?.matchType === 0}
					class:present={letter?.matchType === 1}
					class:hit={letter?.matchType === 2}
				>
					{i === matchInfo!.currentRound - 1
						? matchInfo!.currentGuess[j]
						: letter?.letter
							? String.fromCharCode(letter.letter)
							: ''}
				</div>
			{/each}
		</div>
	{/each}
	<div class="pointer-events-none h-[180px]"></div>
</div>

<style>
	.board {
		display: grid;
		width: 100%;
		gap: 8px;
		justify-content: center;
	}
	.row {
		display: flex;
		gap: 8px;
	}
	.box {
		width: 40px;
		height: 40px;
		border: 1px solid #ccc;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 1.5em;
		background: #fff;
		text-transform: uppercase;
		border-radius: 4px;
	}
	.box.miss {
		background: #9aa7ad;
		border-color: #9aa7ad;
		color: #fff;
	}
	.box.present {
		background: #c9c93e;
		border-color: #c9c93e;
		color: #fff;
	}
	.box.hit {
		background: #6aaa64;
		border-color: #6aaa64;
		color: #fff;
	}
</style>
