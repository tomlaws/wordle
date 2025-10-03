<script lang="ts">
	import { type GameContext, GAME_KEY } from '$lib/context/game-context';
	import { getContext } from 'svelte';
	import Countdown from '../Countdown.svelte';

	const gameContext = getContext<GameContext>(GAME_KEY);
	const { playerInfo, matchInfo } = $derived(gameContext);
</script>

<div
	class="sticky left-0 right-0 top-0 flex min-h-20 w-full flex-col items-center justify-center text-center bg-white dark:bg-[#191e25]"
>
	<div class="flex w-full max-w-[400px] flex-col items-center justify-center px-4">
		<h2 class="text-lg font-semibold text-gray-700">
			{#if matchInfo?.gameOver}
				Game Over
			{:else}
				Round {matchInfo?.currentRound} -
				{matchInfo?.myTurn == null
					? 'Loading'
					: matchInfo?.myTurn
						? 'Your Turn!'
						: 'Waiting for Opponent...'}
			{/if}
		</h2>
		<div class="mt-2 flex w-full items-center justify-between">
			<div
				class="w-1/3 overflow-hidden text-ellipsis whitespace-nowrap rounded-lg px-2 text-center text-gray-800"
				class:bg-blue-200={matchInfo?.myTurn ==
					(matchInfo?.player1.nickname === playerInfo.nickname)}
				class:bg-gray-200={matchInfo?.myTurn !=
					(matchInfo?.player1.nickname === playerInfo.nickname)}
			>
				<span class="text-sm font-bold">{matchInfo?.player1.nickname}</span>
			</div>
			{#if matchInfo?.deadline}
				<div>
					<Countdown deadline={matchInfo?.deadline} />
				</div>
			{/if}
			<div
				class="w-1/3 overflow-hidden text-ellipsis whitespace-nowrap rounded-lg px-2 text-center text-gray-800"
				class:bg-blue-200={matchInfo?.myTurn ==
					(matchInfo?.player2.nickname === playerInfo.nickname)}
				class:bg-gray-200={matchInfo?.myTurn !=
					(matchInfo?.player2.nickname === playerInfo.nickname)}
			>
				<span class="text-sm font-bold">{matchInfo?.player2.nickname}</span>
			</div>
		</div>
	</div>
</div>
