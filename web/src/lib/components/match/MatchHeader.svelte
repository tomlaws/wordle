<script lang="ts">
	import { type GameContext, GAME_KEY } from '$lib/context/game-context';
	import { getContext } from 'svelte';
	import Countdown from '../Countdown.svelte';
	import { twJoin } from 'tailwind-merge';

	const gameContext = getContext<GameContext>(GAME_KEY);
	const { playerInfo, matchInfo } = $derived(gameContext);
	const player1Turn = $derived(
		matchInfo?.myTurn == (matchInfo?.player1.nickname === playerInfo.nickname)
	);
	const player2Turn = $derived(
		matchInfo?.myTurn == (matchInfo?.player2.nickname === playerInfo.nickname)
	);
</script>

{#if !matchInfo?.gameOver}
	<div
		class="sticky left-0 right-0 top-0 z-10 flex min-h-20 w-full flex-col items-center justify-center border-b border-gray-300 bg-white text-center dark:border-gray-700 dark:bg-[#191e25]"
	>
		<div class="flex w-full max-w-[400px] flex-col items-center justify-center px-4">
			<h2 class="text-lg font-semibold uppercase text-gray-700 dark:text-gray-300">
				{#if matchInfo?.gameOver}
					Game Over
				{:else}
					Round {matchInfo?.currentRound}
				{/if}
			</h2>
			<div class="mt-2 flex w-full items-center justify-between">
				<div
					class={twJoin(
						'player w-1/3 overflow-hidden text-ellipsis whitespace-nowrap rounded-lg px-2 text-center text-gray-800 dark:text-white',
						player1Turn ? 'bg-(--color-primary) text-white' : 'bg-gray-200 dark:bg-gray-700'
					)}
				>
					<span class="text-sm font-bold">{matchInfo?.player1.nickname}</span>
				</div>
				{#if matchInfo?.deadline}
					<div>
						<Countdown deadline={matchInfo?.deadline} />
					</div>
				{/if}
				<div
					class={twJoin(
						'player w-1/3 overflow-hidden text-ellipsis whitespace-nowrap rounded-lg px-2 text-center text-gray-800 dark:text-white',
						player2Turn ? 'bg-(--color-primary) text-white' : 'bg-gray-200 dark:bg-gray-700'
					)}
				>
					<span class="text-sm font-bold">{matchInfo?.player2.nickname}</span>
				</div>
			</div>
		</div>
	</div>
{/if}
