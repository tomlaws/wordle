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
	let { websocket, matchInfo } = $derived(gameContext);

	let pressedKey: string | null = $state(null);

	function submitGuess() {
		matchInfo!.loading = true;
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
		if (key === 'Enter' || key === 'Backspace' || /^[A-Z]$/.test(key)) {
			event.preventDefault();
			pressedKey = key;
			handleKey(key);
		}
	}

	function onKeyUp(event: KeyboardEvent) {
		const key = event.key.length === 1 ? event.key.toUpperCase() : event.key;
		if (key === 'Enter' || key === 'Backspace' || /^[A-Z]$/.test(key)) {
			pressedKey = null;
		}
	}

	onMount(() => {
		console.log('Keyboard component mounted, adding keydown listener');
		window.addEventListener('keydown', onKeyDown);
		window.addEventListener('keyup', onKeyUp);
	});

	onDestroy(() => {
		console.log('Keyboard component unmounted, removing keydown listener');
		window.removeEventListener('keydown', onKeyDown);
		window.removeEventListener('keyup', onKeyUp);
	});
</script>

<div
	class="fixed bottom-0 left-0 right-0 flex h-[170px] flex-col items-center justify-center border-t border-t-gray-300 bg-white"
>
	{#if !matchInfo?.myTurn && matchInfo}
		<div class="keyboard-overlay">Waiting for Opponent...</div>
	{/if}
	<div class="keyboard" class:keyboard--disabled={!matchInfo?.myTurn || matchInfo?.loading}>
		{#each keyboardRows as row}
			<div class="key-row">
				{#each row as key}
					<button
						class="key {pressedKey === key ? 'pressed' : ''}"
						disabled={!matchInfo?.myTurn || matchInfo?.loading}
						onmousedown={() => (pressedKey = key)}
						onmouseup={() => (pressedKey = null)}
						onmouseleave={() => (pressedKey = null)}
						onclick={() => handleKey(key)}>{key}</button
					>
				{/each}
			</div>
		{/each}
	</div>
</div>

<style>
	.keyboard {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
	}
	.keyboard--disabled {
		opacity: 0.5;
		pointer-events: none;
		user-select: none;
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
		transition:
			background 0.1s,
			transform 0.05s;
	}
	.key:hover:enabled {
		background: #ddd;
	}
	.key.pressed:enabled,
	.key:active:enabled {
		background: #ccc;
		transform: scale(0.96);
	}
	.keyboard-overlay {
		position: absolute;
		inset: 0;
		width: 100vw;
		background: rgba(255, 255, 255, 0.75);
		text-align: center;
		padding: 16px 0;
		font-size: 1.2em;
		font-weight: bold;
		z-index: 1000;
		pointer-events: none;
		user-select: none;
		display: flex;
		align-items: center;
		justify-content: center;
	}
</style>
