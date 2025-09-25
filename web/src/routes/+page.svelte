<script lang="ts">
	import { webSocket } from 'rxjs/webSocket';
	import { map, filter, first } from 'rxjs/operators';
	import Match from '$lib/components/Match.svelte';
	import { Protocol, type Message, type Payload } from '$lib/utils/message';
	import { payloadRegistry } from './payload-registry';
	import { createWebSocket, type WebSocketConnection } from '$lib/utils/websocket';
	import { PlayerInfoPayload } from '$lib/types/payload';
	let nickname = '';
	let entered = false;
	let websocket: WebSocketConnection<Payload>;
	let playerInfo: { id: string; nickname: string } | null = null;

	function enterGame() {
		if (nickname.trim()) {
			entered = true;
		}
		const protocol = new Protocol(payloadRegistry);
		if (!websocket) {
			websocket = createWebSocket(
				'ws://127.0.0.1:8080/socket?nickname=' + encodeURIComponent(nickname),
				(payload: Payload) => protocol.createMessage(payload),
				(msg: Message) => protocol.parseMessage(msg)
			);
		}
		// wait for player info from server
		websocket.messages$
			.pipe(
				filter((msg) => msg instanceof PlayerInfoPayload),
				first(),
			)
			.subscribe((payload) => {
				playerInfo = payload as { id: string; nickname: string };
			});
	}
</script>

{#if !playerInfo}
	<div class="mt-0 flex min-h-screen flex-col items-center justify-center gap-4">
		<label for="nickname" class="text-lg font-semibold">Enter your nickname:</label>
		<input
			id="nickname"
			type="text"
			bind:value={nickname}
			placeholder="Nickname"
			class="rounded-md border border-gray-300 px-3 py-2 text-base focus:outline-none focus:ring-2 focus:ring-blue-500"
			on:keydown={(e) => e.key === 'Enter' && enterGame()}
		/>
		<button
			on:click={enterGame}
			disabled={!nickname.trim()}
			class="rounded-md bg-blue-600 px-4 py-2 font-medium text-white transition hover:bg-blue-700 disabled:cursor-not-allowed disabled:bg-gray-400"
		>
			Enter Game
		</button>
	</div>
{:else}
	<!-- Game UI goes here -->
	<!-- <Lobby {nickname} /> -->
	<Match {websocket} {playerInfo} />
{/if}
