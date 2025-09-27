<script lang="ts">
	import { GAME_KEY, type GameContext } from '$lib/context/game-context';
	import { TOAST_KEY, type ToastAPI } from '$lib/context/toast-context';
	import {
		FeedbackPayload,
		GameOverPayload,
		GuessPayload,
		GuessTimeoutPayload,
		InvalidWordPayload,
		PlayAgainPayload,
		RoundStartPayload
	} from '$lib/types/payload';
	import { getContext, onMount } from 'svelte';

	// State for guesses and current input
	let {
		rounds = 12
	}: {
		rounds?: number;
	} = $props();
	let { websocket, playerInfo } = getContext<GameContext>(GAME_KEY);
	let myTurn = $state<boolean | null>(null);
	let loading = $state(false);
	let gameOver = $state<GameOverPayload | null>(null);
	let guesses = $state<Array<FeedbackPayload['feedback']>>(
		Array.from({ length: rounds }, () => Array(5).fill(null))
	);
	let currentRound = $state(-1);
	let currentGuess = $state(Array(5).fill(''));

	const keyboardRows = [
		['Q', 'W', 'E', 'R', 'T', 'Y', 'U', 'I', 'O', 'P'],
		['A', 'S', 'D', 'F', 'G', 'H', 'J', 'K', 'L'],
		['Enter', 'Z', 'X', 'C', 'V', 'B', 'N', 'M', 'Backspace']
	];
	const toast = getContext<ToastAPI>(TOAST_KEY);
	function submitGuess() {
		loading = true;
		console.log(currentGuess.join(''));
		const guessPayload = new GuessPayload();
		guessPayload.word = currentGuess.join('');
		console.log(guessPayload);
		websocket.send(guessPayload);
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
	function playAgain(confirm: boolean) {
		if (confirm) {
			const playAgainPayload = new PlayAgainPayload();
			playAgainPayload.confirm = true;
			websocket.send(playAgainPayload);
		} else {
			location.reload();
		}
	}
	onMount(() => {
		websocket.messages$.subscribe((msg) => {
			if (msg instanceof RoundStartPayload) {
				myTurn = msg.player.id === playerInfo.id;
				if (myTurn) {
					currentRound = msg.round;
					currentGuess = Array(5).fill('');
				}
			}
			if (msg instanceof InvalidWordPayload) {
				console.log('Invalid word received', msg);
				if (msg.round === currentRound) {
					if (msg.player.id === playerInfo.id) {
						loading = false;
						currentGuess = Array(5).fill('');
						toast.info(`${msg.word} is not a valid word.`);
					} else {
						toast.info(`${msg.player.nickname} guessed an invalid word ${msg.word}.`);
					}
				}
			}
			if (msg instanceof GuessTimeoutPayload) {
				loading = false;
				if (msg.player.id === playerInfo.id) {
					toast.error('You ran out of time!');
				} else {
					toast.info(`${msg.player.nickname} ran out of time.`);
				}
				guesses[msg.round - 1] = Array.from({ length: 5 }, (_, i) => ({
					position: i,
					letter: '-'.charCodeAt(0),
					matchType: 0
				}));
			}
			if (msg instanceof FeedbackPayload) {
				loading = false;
				msg.feedback.forEach((item) => {
					guesses[msg.round - 1][item.position] = item;
				});
			}
			if (msg instanceof GameOverPayload) {
				gameOver = msg;
				toast.info(`Game over! The word was ${msg.answer}.`);
				if (msg.winner) {
					if (msg.winner.id === playerInfo.id) {
						toast.success('You won!');
					} else {
						toast.error(`${msg.winner.nickname} won the game.`);
					}
				} else {
					toast.info('The game ended in a draw.');
				}
			}
		});
	});
</script>

{#if gameOver}
	<h2>Game Over</h2>
	<p>The word was {gameOver.answer}.</p>
	{#if gameOver.winner}
		{#if gameOver.winner.id === playerInfo.id}
			<p>Congratulations, you won!</p>
		{:else}
			<p>{gameOver.winner.nickname} won the game.</p>
		{/if}
	{:else}
		<p>The game ended in a draw.</p>
	{/if}
	<button onclick={() => playAgain(true)}>Play Again</button>
	<button onclick={() => playAgain(false)}>Quit</button>
{:else}
	<h2>{myTurn == null ? 'Loading' : myTurn ? 'Your Turn!' : 'Waiting for Opponent...'}</h2>
	<div class="board">
		{#each guesses as guess, i}
			<div class="row">
				{#each guess as letter, j}
					<div
						class="box"
						class:miss={letter?.matchType === 0}
						class:present={letter?.matchType === 1}
						class:hit={letter?.matchType === 2}
					>
						{i === currentRound - 1
							? currentGuess[j]
							: letter?.letter
								? String.fromCharCode(letter.letter)
								: ''}
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
{/if}

<style>
	.board {
		display: grid;
		gap: 8px;
		justify-content: center;
		margin-top: 32px;
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
	.box.miss {
		background: #787c7e;
		border-color: #787c7e;
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
