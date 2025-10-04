<script lang="ts">
	import Match from '$lib/components/Match.svelte';
	import { Protocol, type Message, type Payload } from '$lib/utils/message';
	import { payloadRegistry } from './payload-registry';
	import { createWebSocket } from '$lib/utils/websocket';
	import { GameStartPayload, MatchingPayload, PlayerInfoPayload } from '$lib/types/payload';
	import { onMount, setContext } from 'svelte';
	import { GAME_KEY, type GameContext } from '$lib/context/game-context';
	import Lobby from '$lib/components/Lobby.svelte';
	import { GameState } from '$lib/types/state';
	import Button from '$lib/components/common/Button.svelte';
	import Input from '$lib/components/common/Input.svelte';
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
				// find header element and make it visible
				document.getElementById('header')?.classList.remove('hidden');
			}
			if (msg instanceof GameStartPayload) {
				gameState = GameState.IN_GAME;
				gameContext.matchInfo = {
					loading: true,
					player1: msg.player1,
					player2: msg.player2,
					guesses: Array.from({ length: 12 }, () => Array(5).fill(null)),
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
		<Button onclick={enterGame} disabled={!nickname.trim()}>Play</Button>
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
