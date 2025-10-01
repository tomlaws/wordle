<script lang="ts">
    import { onMount, onDestroy } from 'svelte';

    export let deadline: Date;

    let timeLeft = '';
    let interval: ReturnType<typeof setInterval>;
    let secondsLeft = 0;

    function updateTimeLeft() {
        const now = new Date();
        const diff = deadline.getTime() - now.getTime();

        if (diff <= 0) {
            timeLeft = '00:00';
            secondsLeft = 0;
            clearInterval(interval);
            return;
        }

        const minutes = Math.floor((diff / (1000 * 60)) % 60);
        const seconds = Math.floor((diff / 1000) % 60);
        secondsLeft = Math.floor(diff / 1000);

        timeLeft = [
            minutes.toString().padStart(2, '0'),
            seconds.toString().padStart(2, '0')
        ].join(':');
    }

    onMount(() => {
        updateTimeLeft();
        interval = setInterval(updateTimeLeft, 1000);
    });

    onDestroy(() => {
        clearInterval(interval);
    });
</script>

<span style="color: {secondsLeft < 10 ? 'red' : 'inherit'}">{timeLeft}</span>
