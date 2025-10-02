<script lang="ts">
	import { type GameContext, GAME_KEY } from '$lib/context/game-context';
	import { GuessPayload } from '$lib/types/payload';
	import { onMount, onDestroy, getContext } from 'svelte';

	const keyboardRows = [
		['Q', 'W', 'E', 'R', 'T', 'Y', 'U', 'I', 'O', 'P'],
		['A', 'S', 'D', 'F', 'G', 'H', 'J', 'K', 'L'],
		['Enter', 'Z', 'X', 'C', 'V', 'B', 'N', 'M', 'Backspace']
	];

	const gameContext = getContext<GameContext>(GAME_KEY);
	let { websocket, matchInfo, loading } = $derived(gameContext);

	function submitGuess() {
		loading = true;
		console.log(matchInfo!.currentGuess.join(''));
		const guessPayload = new GuessPayload();
		guessPayload.word = matchInfo!.currentGuess.join('');
		console.log(guessPayload);
		websocket.send(guessPayload);
	}

	function handleKey(key) {
		if (!matchInfo?.myTurn) return;
		if (key === 'Enter') {
			if (matchInfo!.currentGuess.every((l) => l)) {
				submitGuess();
			}
		} else if (key === 'Backspace') {
			let idx = matchInfo!.currentGuess.findLastIndex((l) => l);
			if (idx !== -1) matchInfo!.currentGuess[idx] = '';
		} else if (/^[A-Z]$/.test(key)) {
			let idx = matchInfo!.currentGuess.findIndex((l) => !l);
			if (idx !== -1) matchInfo!.currentGuess[idx] = key;
		}
	}

	function onKeyDown(event: KeyboardEvent) {
		const key = event.key.length === 1 ? event.key.toUpperCase() : event.key;
		if (
			key === 'Enter' ||
			key === 'Backspace' ||
			/^[A-Z]$/.test(key)
		) {
			event.preventDefault();
			handleKey(key);
		}
	}

	onMount(() => {
		console.log('Keyboard component mounted, adding keydown listener');
		window.addEventListener('keydown', onKeyDown);
	});

	onDestroy(() => {
		console.log('Keyboard component unmounted, removing keydown listener');
		window.removeEventListener('keydown', onKeyDown);
	});
</script>

<div class="keyboard">
	{#each keyboardRows as row}
		<div class="key-row">
			{#each row as key}
				<button
					class="key"
					disabled={!matchInfo?.myTurn || loading}
					onclick={() => handleKey(key)}>{key}</button
				>
			{/each}
		</div>
	{/each}
</div>

<style>
	.keyboard {
		position: fixed;
		left: 0;
		right: 0;
		bottom: 0;
		height: 170px;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		background: #fff;
		border-top: 1px solid #ccc;
	}
	.key-row {
		display: flex;
		justify-content: center;
		gap: 6px;
		margin-bottom: 6px;
	}
	.key {
		min-width: 32px;
		padding: 8px 12px;
		background: #eee;
		border: none;
		border-radius: 4px;
		font-size: 1em;
		cursor: pointer;
		text-transform: uppercase;
	}
</style>
