<script lang="ts">
	import { onMount } from 'svelte';

	export let deadline: Date;

	let timeLeft = '';
	let secondsLeft = 0;

	function updateTimeLeft() {
		const now = new Date();
		const diff = deadline.getTime() - now.getTime();

		if (diff <= 0) {
			timeLeft = '00:00';
			secondsLeft = 0;
			return;
		}

		const minutes = Math.floor((diff / (1000 * 60)) % 60);
		const seconds = Math.floor((diff / 1000) % 60);
		secondsLeft = Math.floor(diff / 1000);

		timeLeft = [minutes.toString().padStart(2, '0'), seconds.toString().padStart(2, '0')].join(':');
	}

	onMount(() => {
		updateTimeLeft();
		const interval = setInterval(updateTimeLeft, 1000);
		return () => clearInterval(interval);
	});
</script>

<span class="countdown tracking-wide {secondsLeft <= 15 ? 'urgent' : ''}">
    {timeLeft}
</span>

<style>
    .countdown {
        font-family: monospace;
        font-size: 1rem;
        font-weight: 600;
        color: var(--color-gray-600);
    }
    .countdown.urgent {
        color: var(--color-red-400);
        animation: pulse 1s infinite;
    }
    @keyframes pulse {
        0%, 100% {
            opacity: 1;
        }
        50% {
            opacity: 0.5;
        }
    }
</style>