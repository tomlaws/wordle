<script lang="ts">
	import { type GameContext, GAME_KEY } from '$lib/context/game-context';
	import { getContext } from 'svelte';
	import Countdown from '../Countdown.svelte';

	export let currentRound: number;
	export let myTurn: boolean | null;
	export let deadline: Date | null;

	const { playerInfo, matchInfo } = getContext<GameContext>(GAME_KEY);
    const myNickname = playerInfo?.nickname;
</script>

<div
	class="sticky top-0 left-0 right-0 bg-white flex min-h-20 w-full flex-col items-center justify-center border-b border-b-[1px] border-b-gray-300 text-center"
>
	<div class="flex w-[400px] flex-col items-center justify-center">
		<h2>
			Round {currentRound} -
			{myTurn == null ? 'Loading' : myTurn ? 'Your Turn!' : 'Waiting for Opponent...'}
		</h2>
		<div class="mt-2 flex w-full items-center justify-between">
            <div
                class="text-center w-1/3 rounded-lg px-2 text-ellipsis overflow-hidden whitespace-nowrap"
                class:bg-yellow-200={myTurn == (matchInfo?.player1.nickname === myNickname)}
                class:bg-gray-200={myTurn != (matchInfo?.player1.nickname === myNickname)}
            >
                <span class="font-semibold">{matchInfo?.player1.nickname}</span>
            </div>
            {#if deadline}
                <div>
                    <span class="ml-4 text-sm text-gray-500">
                        <Countdown {deadline} />
                    </span>
                </div>
            {/if}
            <div
                class="text-center w-1/3 rounded-lg px-2 text-ellipsis overflow-hidden whitespace-nowrap"
                class:bg-yellow-200={myTurn == (matchInfo?.player2.nickname === myNickname)}
                class:bg-gray-200={myTurn != (matchInfo?.player2.nickname === myNickname)}
            >
                <span class="font-semibold">{matchInfo?.player2.nickname}</span>
            </div>
		</div>
	</div>
</div>