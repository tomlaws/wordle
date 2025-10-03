<script lang="ts">
	import { type GameContext, GAME_KEY } from '$lib/context/game-context';
	import { GuessPayload, TypingPayload } from '$lib/types/payload';
	import { onMount, onDestroy, getContext } from 'svelte';
	import { tooltip, Tooltip } from '@svelte-plugins/tooltips';
	import IconArrowLeft from './icons/IconArrowLeft.svelte';
	import IconSendSolid from './icons/IconSendSolid.svelte';

	const keyboardRows = [
		['Q', 'W', 'E', 'R', 'T', 'Y', 'U', 'I', 'O', 'P'],
		['A', 'S', 'D', 'F', 'G', 'H', 'J', 'K', 'L'],
		['Enter', 'Z', 'X', 'C', 'V', 'B', 'N', 'M', 'Backspace']
	];

	const gameContext = getContext<GameContext>(GAME_KEY);
	let { websocket, matchInfo } = $derived(gameContext);

	let pressedKey: string | null = $state(null);

	function sendTypingUpdate() {
		const typingPayload = new TypingPayload();
		typingPayload.word = matchInfo?.currentGuess.join('') || '';
		websocket.send(typingPayload);
	}

	function submitGuess() {
		matchInfo!.loading = true;
		const guessPayload = new GuessPayload();
		guessPayload.word = matchInfo!.currentGuess.join('');
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
			sendTypingUpdate();
		} else if (/^[A-Z]$/.test(key)) {
			let idx = matchInfo!.currentGuess.findIndex((l) => !l);
			if (idx !== -1) matchInfo!.currentGuess[idx] = key;
			sendTypingUpdate();
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
	class="fixed bottom-0 left-0 right-0 flex h-[170px] flex-col items-center justify-center bg-white dark:bg-[#191e25]"
>
	{#if !matchInfo?.myTurn && matchInfo}
		<div class="keyboard-overlay bg-white opacity-75 dark:bg-[#191e25]">
			Waiting for Opponent...
		</div>
	{/if}
	<div class="keyboard" class:keyboard--disabled={!matchInfo?.myTurn || matchInfo?.loading}>
		{#each keyboardRows as row}
			<div class="key-row">
				{#each row as key}
					<Tooltip content={key === 'Backspace' ? 'Delete' : (key === 'Enter' ? 'Submit' : '')} show={false}>
						<button
							class="key {pressedKey === key
								? 'pressed'
								: ''} bg-gray-200 text-gray-900 hover:bg-gray-300 active:bg-gray-400 dark:bg-gray-700 dark:text-gray-100 dark:hover:bg-gray-600 dark:active:bg-gray-500"
							disabled={!matchInfo?.myTurn || matchInfo?.loading}
							onmousedown={() => (pressedKey = key)}
							onmouseup={() => (pressedKey = null)}
							onmouseleave={() => (pressedKey = null)}
							onclick={() => handleKey(key)}
						>
							{#if key === 'Backspace'}
								<IconArrowLeft class="mx-auto p-1" />
							{:else if key === 'Enter'}
								<IconSendSolid class="mx-auto p-1" />
							{:else}
								{key}
							{/if}
						</button></Tooltip
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
		width: clamp(30px, 6vw, 32px);
		padding: clamp(6px, 1.5vw, 8px) 0;
		border: none;
		border-radius: 4px;
		cursor: pointer;
		text-transform: uppercase;
		transition:
			background 0.1s,
			transform 0.05s;
	}

	.key.pressed:enabled,
	.key:active:enabled {
		transform: scale(0.96);
	}
	.keyboard-overlay {
		position: absolute;
		inset: 0;
		width: 100vw;
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
