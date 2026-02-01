<script lang="ts">
	import Match from '$lib/components/Match.svelte';
	import { Protocol, type Message, type Payload } from '$lib/utils/message';
	import { payloadRegistry } from './payload-registry';
	import { createWebSocket } from '$lib/utils/websocket';
	import { GameStartPayload, MatchingPayload, PlayerInfoPayload } from '$lib/types/payload';
	import { getContext, setContext } from 'svelte';
	import { GAME_KEY, type GameContext } from '$lib/context/game-context';
	import Lobby from '$lib/components/Lobby.svelte';
	import { GameState } from '$lib/types/state';
	import Button from '$lib/components/common/Button.svelte';
	import Input from '$lib/components/common/Input.svelte';
	import { catchError, of } from 'rxjs';
	import { type ToastAPI, TOAST_KEY } from '$lib/context/toast-context';
	import { PUBLIC_WEBSOCKET_SERVER } from '$env/static/public';
	let nickname = $state('');
	let gameState = $state<GameState>(GameState.UNAUTHENTICATED);
	let gameContext = $state<Partial<GameContext>>({});
	setContext<Partial<GameContext>>(GAME_KEY, gameContext);
	const toast = getContext<ToastAPI>(TOAST_KEY);
	let isConnecting = $state(false);

	function enterGame() {
		if (isConnecting) return;
		const trimmed = nickname.trim();
		if (!trimmed || trimmed.length < 3 || trimmed.length > 16) {
			toast.error('Nickname must be between 3 and 16 characters long.');
			return;
		}
		isConnecting = true;
		const protocol = new Protocol(payloadRegistry);
		if (!gameContext.websocket) {
			gameContext.websocket = createWebSocket(
				PUBLIC_WEBSOCKET_SERVER + '/socket?nickname=' + encodeURIComponent(nickname),
				(payload: Payload) => protocol.createMessage(payload),
				(msg: Message) => protocol.parseMessage(msg)
			);
		}
		gameContext.websocket.messages$
			.pipe(
				catchError((error) => {
					console.error('WebSocket error:', error);
					isConnecting = false;
					return of([]);
				})
			)
			.subscribe((msg) => {
				if (msg instanceof PlayerInfoPayload) {
					gameState = GameState.AUTHENTICATED;
					gameContext.playerInfo = { id: msg.id, nickname: msg.nickname };
				}
				if (msg instanceof MatchingPayload) {
					gameState = GameState.MATCHING;
					// find header element and make it visible
					document.getElementById('header')?.classList.remove('hidden');
				}
				if (msg instanceof GameStartPayload) {
					gameState = GameState.IN_GAME;
					gameContext.matchInfo = {
						loading: true,
						player1: msg.player1,
						player2: msg.player2,
						guesses: Array.from({ length: msg.maxGuesses }, () => Array(5).fill(null)),
						currentRound: -1,
						currentGuess: Array(5).fill(''),
						myTurn: false
					};
					// find header element and make it invisible
					document.getElementById('header')?.classList.add('hidden');
				}
			});
	}
</script>

{#if gameState == GameState.UNAUTHENTICATED}
	<div class="mt-0 flex min-h-screen flex-col items-center justify-center gap-4">
		<h1 class="text-xl font-bold uppercase">Welcome!</h1>
		<label for="nickname" class="text-m text-sm font-semibold">
			What's your Wordle warrior name? 🌟
		</label>
		<Input
			class="my-2"
			placeholder="Nickname"
			bind:value={nickname}
			onkeydown={(e) => e.key === 'Enter' && enterGame()}
		/>
		<Button onclick={enterGame} disabled={!nickname.trim() || isConnecting}>
			{#if isConnecting}
				<div class="flex items-center gap-2">
					<svg
						class="animate-spin h-5 w-5"
						xmlns="http://www.w3.org/2000/svg"
						fill="none"
						viewBox="0 0 24 24"
					>
						<circle
							class="opacity-25"
							cx="12"
							cy="12"
							r="10"
							stroke="currentColor"
							stroke-width="4"
						></circle>
						<path
							class="opacity-75"
							fill="currentColor"
							d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
						></path>
					</svg>
					Connecting
				</div>
			{:else}
				Play
			{/if}
		</Button>
	</div>
{:else}
	<!-- Game UI goes here -->
	{#if gameState === GameState.MATCHING}
		<Lobby />
	{:else if gameState === GameState.IN_GAME}
		<Match />
	{:else}
		<p>Loading...</p>
	{/if}
{/if}
