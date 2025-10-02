<script lang="ts">
	import Match from '$lib/components/Match.svelte';
	import { Protocol, type Message, type Payload } from '$lib/utils/message';
	import { payloadRegistry } from './payload-registry';
	import { createWebSocket, type WebSocketConnection } from '$lib/utils/websocket';
	import { FeedbackPayload, GameStartPayload, MatchingPayload, PlayerInfoPayload } from '$lib/types/payload';
	import { onMount, setContext } from 'svelte';
	import { GAME_KEY, type GameContext } from '$lib/context/game-context';
	import Lobby from '$lib/components/Lobby.svelte';
	import { GameState } from '$lib/types/state';
	let nickname = $state('');
	let gameState = $state<GameState>(GameState.UNAUTHENTICATED);
	let gameContext = $state<Partial<GameContext>>({});
	setContext<Partial<GameContext>>(GAME_KEY, gameContext);
	onMount(() => {
		console.log('Page component mounted');
		return () => {
			console.log('Page component unmounted');
		};
	});
	function enterGame() {
		console.log('Entering game with nickname:', nickname);
		const protocol = new Protocol(payloadRegistry);
		if (!gameContext.websocket) {
			gameContext.websocket = createWebSocket(
				'ws://127.0.0.1:8080/socket?nickname=' + encodeURIComponent(nickname),
				(payload: Payload) => protocol.createMessage(payload),
				(msg: Message) => protocol.parseMessage(msg)
			);
		}
		gameContext.websocket.messages$.subscribe((msg) => {
			console.log('Received message', msg);
			if (msg instanceof PlayerInfoPayload) {
				gameState = GameState.AUTHENTICATED;
				gameContext.playerInfo = { id: msg.id, nickname: msg.nickname };
			}
			if (msg instanceof MatchingPayload) {
				gameState = GameState.MATCHING;
			}
			if (msg instanceof GameStartPayload) {
				gameState = GameState.IN_GAME;
				gameContext.matchInfo = {
					player1: msg.player1,
					player2: msg.player2,
					guesses: Array.from({ length: 12 }, () => Array(5).fill(null)) ,
					currentRound: -1,
					currentGuess: Array(5).fill(''),
					myTurn: false,
				};
			}
		});
	}
</script>

{#if gameState == GameState.UNAUTHENTICATED}
	<div class="mt-0 flex min-h-screen flex-col items-center justify-center gap-4">
		<label for="nickname" class="text-lg font-semibold">Enter your nickname:</label>
		<input
			id="nickname"
			type="text"
			bind:value={nickname}
			placeholder="Nickname"
			class="rounded-md border border-gray-300 px-3 py-2 text-base focus:outline-none focus:ring-2 focus:ring-blue-500"
			onkeydown={(e) => e.key === 'Enter' && enterGame()}
		/>
		<button
			onclick={enterGame}
			disabled={!nickname.trim()}
			class="rounded-md bg-blue-600 px-4 py-2 font-medium text-white transition hover:bg-blue-700 disabled:cursor-not-allowed disabled:bg-gray-400"
		>
			Enter Game
		</button>
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
