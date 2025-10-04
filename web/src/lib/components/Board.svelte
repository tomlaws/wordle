<script lang="ts">
	import { type GameContext, GAME_KEY } from '$lib/context/game-context';
	import { getContext, onMount } from 'svelte';
	import IconDazeSquare from './icons/IconDazeSquare.svelte';
	import { Tooltip } from '@svelte-plugins/tooltips';

	const gameContext = getContext<GameContext>(GAME_KEY);
	const { matchInfo } = $derived(gameContext);

	$effect(() => {
		document
			.getElementById('row-' + matchInfo!.currentRound)
			?.scrollIntoView({ behavior: 'smooth', block: 'center' });
	});

	function getCharacter(letter) {
		if (!letter) return null;
		if (letter.letter === '-'.charCodeAt(0)) {
			return null;
		}
		return String.fromCharCode(letter.letter);
	}
</script>

<div class="board">
	{#each matchInfo!.guesses as guess, i}
		<div class="row" id={'row-' + (i + 1)}>
			{#each guess as letter, j}
				<Tooltip
					content={letter?.letter === '-'.charCodeAt(0) && i < matchInfo!.currentRound - 1
						? 'Timeout'
						: (
							letter?.matchType === 2
								? 'Hit'
								: (letter?.matchType === 1 ? 'Present' : (letter?.matchType === 0 ? 'Miss' : ''))
						)}
					show={false}
				>
					<div
						class="box"
						class:miss={letter?.matchType === 0}
						class:present={letter?.matchType === 1}
						class:hit={letter?.matchType === 2}
						class:current={i === matchInfo!.currentRound - 1 &&
							matchInfo?.currentGuess.join('').length == j}
					>
						{i === matchInfo!.currentRound - 1 ? matchInfo!.currentGuess[j] : getCharacter(letter)}
						{#if getCharacter(letter) == null && i < matchInfo!.currentRound - 1 && j == 2}
							<IconDazeSquare class="h-6 w-6" />
						{/if}
					</div>
				</Tooltip>
			{/each}
		</div>
	{/each}
	<div class="pointer-events-none h-[184px]"></div>
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
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 1.5em;
		background: var(--color-input-bg);
		text-transform: uppercase;
		border-radius: 4px;
	}
	.box.miss {
		background: #f38484;
		color: #fff;
	}
	.box.present {
		background: #64c4fd;
		color: #fff;
	}
	.box.hit {
		background: #7fc579;
		color: #fff;
	}
	.box.current {
		--tw-ring-color: var(--color-primary);
		--tw-ring-shadow: var(--tw-ring-inset,) 0 0 0 calc(2px + var(--tw-ring-offset-width))
			var(--tw-ring-color);
		box-shadow:
			var(--tw-inset-shadow), var(--tw-inset-ring-shadow), var(--tw-ring-offset-shadow),
			var(--tw-ring-shadow), var(--tw-shadow);
	}
</style>
