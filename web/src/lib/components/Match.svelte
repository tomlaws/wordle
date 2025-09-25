<script lang="ts">
	import { TOAST_KEY, type ToastAPI } from '$lib/context/toast-context';
	import {
		FeedbackPayload,
		GuessPayload,
		InvalidWordPayload,
		RoundStartPayload
	} from '$lib/types/payload';
	import type { Payload } from '$lib/utils/message';
	import type { WebSocketConnection } from '$lib/utils/websocket';
	import { getContext } from 'svelte';

	let {
		websocket,
		playerInfo,
		rounds = 12
	}: {
		websocket: WebSocketConnection<Payload>;
		playerInfo: { id: string; nickname: string };
		rounds?: number;
	} = $props();
	// State for guesses and current input
	let myTurn = $state<boolean | null>(null);
	let loading = $state(false);
	let guesses = Array(rounds).fill(Array(5).fill(''));
	let currentRound = $state(0);
	let currentGuess = $state(Array(5).fill(''));
	const keyboardRows = [
		['Q', 'W', 'E', 'R', 'T', 'Y', 'U', 'I', 'O', 'P'],
		['A', 'S', 'D', 'F', 'G', 'H', 'J', 'K', 'L'],
		['Enter', 'Z', 'X', 'C', 'V', 'B', 'N', 'M', 'Backspace']
	];
	const toast = getContext<ToastAPI>(TOAST_KEY);
	function submitGuess() {
		loading = true;
		websocket.send(new GuessPayload({ guess: currentGuess.join('') }));
	}
	function handleKey(key) {
		if (!myTurn) return;
		if (key === 'Enter') {
			if (currentGuess.every((l) => l)) {
				submitGuess();
			}
		} else if (key === 'Backspace') {
			let idx = currentGuess.findLastIndex((l) => l);
			if (idx !== -1) currentGuess[idx] = '';
		} else if (/^[A-Z]$/.test(key)) {
			let idx = currentGuess.findIndex((l) => !l);
			if (idx !== -1) currentGuess[idx] = key;
		}
	}
	websocket.messages$.subscribe((msg) => {
		if (msg instanceof RoundStartPayload) {
			console.log('Round started', msg);
			myTurn = msg.player.id === playerInfo.id;
			if (myTurn) {
				currentRound = msg.round;
				currentGuess = Array(5).fill('');
			}
		}
		if (msg instanceof InvalidWordPayload) {
			console.log('Invalid word', msg);
			//loading = false;
			toast.error('Invalid word, try again.');
		}
		if (msg instanceof FeedbackPayload) {
			loading = false;
		}
	});
</script>

<div class="board">
	<h2>{myTurn == null ? 'Loading' : myTurn ? 'Your Turn!' : 'Waiting for Opponent...'}</h2>
	{#each guesses as guess, i}
		<div class="row">
			{#each guess as letter, j}
				<div class="box">
					{i === currentRound ? currentGuess[j] : letter}
				</div>
			{/each}
		</div>
	{/each}
</div>

<div class="keyboard">
	{#each keyboardRows as row}
		<div class="key-row">
			{#each row as key}
				<button class="key" disabled={!myTurn || loading} onclick={() => handleKey(key)}
					>{key}</button
				>
			{/each}
		</div>
	{/each}
</div>

<style>
	.board {
		display: grid;
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
		border: 2px solid #ccc;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 1.5em;
		background: #fff;
		text-transform: uppercase;
	}
	.keyboard {
		margin-top: 24px;
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
