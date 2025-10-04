<script lang="ts">
	import { GAME_KEY, type GameContext } from '$lib/context/game-context';
	import { TOAST_KEY, type ToastAPI } from '$lib/context/toast-context';
	import {
		FeedbackPayload,
		GameOverPayload,
		GuessTimeoutPayload,
		InvalidWordPayload,
		RoundStartPayload,

		TypingPayload

	} from '$lib/types/payload';
	import { getContext, onMount } from 'svelte';
	import MatchHeader from './match/MatchHeader.svelte';
	import Keyboard from './Keyboard.svelte';
	import Board from './Board.svelte';
	import GameOver from './GameOver.svelte';

	// State for guesses and current input
	const gameContext = getContext<GameContext>(GAME_KEY);
	const { websocket, playerInfo } = gameContext;
	const toast = getContext<ToastAPI>(TOAST_KEY);
	const { matchInfo } = $derived(gameContext);

	onMount(() => {
		console.log('Match component mounted, subscribing to websocket messages');
		const sub = websocket.messages$.subscribe((msg) => {
			if (msg instanceof RoundStartPayload) {
				matchInfo!.myTurn = msg.player.id === playerInfo.id;
				matchInfo!.currentRound = msg.round;
				matchInfo!.deadline = msg.getDeadline();
				matchInfo!.currentGuess = Array(5).fill('');
				matchInfo!.loading = false;
				// if (matchInfo!.myTurn) {
				// 	toast.info(`Round ${msg.round} started! It's your turn.`);
				// }
			}
			if (msg instanceof InvalidWordPayload) {
				console.log('Invalid word received', msg);
				if (msg.round === matchInfo!.currentRound) {
					if (msg.player.id === playerInfo.id) {
						matchInfo!.loading = false;
						matchInfo!.currentGuess = Array(5).fill('');
						toast.info(`${msg.word} is not a valid word.`);
					} else {
						toast.info(`${msg.player.nickname} guessed an invalid word ${msg.word}.`);
					}
				}
			}
			if (msg instanceof GuessTimeoutPayload) {
				matchInfo!.loading = false;
				if (msg.player.id === playerInfo.id) {
					toast.error('You ran out of time!');
				} else {
					// toast.info(`${msg.player.nickname} ran out of time.`);
				}
				matchInfo!.guesses[msg.round - 1] = Array.from({ length: 5 }, (_, i) => ({
					position: i,
					letter: '-'.charCodeAt(0),
					matchType: 0
				}));
				matchInfo!.guesses = matchInfo!.guesses;
			}
			if (msg instanceof FeedbackPayload) {
				matchInfo!.loading = false;
				msg.feedback.forEach((item) => {
					matchInfo!.guesses[msg.round - 1][item.position] = item;
				});
				matchInfo!.guesses = matchInfo!.guesses;
			}
			if (msg instanceof GameOverPayload) {
				matchInfo!.gameOver = msg;
				matchInfo!.deadline = undefined;
			}
			if (msg instanceof TypingPayload) {
				matchInfo!.currentGuess = msg.word.split('');
			}
		});

		return () => {
			sub.unsubscribe();
			console.log('Match component unmounted, unsubscribed from websocket messages');
		};
	});
</script>

<div class="flex min-h-screen flex-col items-center justify-center">
	<MatchHeader />
	<div class="flex w-full flex-1 flex-col items-center justify-center py-4">
		{#if matchInfo?.gameOver}
			<GameOver />
		{:else}
			<Board />
		{/if}
	</div>
	{#if !matchInfo?.gameOver}
		<Keyboard />
	{/if}
</div>
