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

        timeLeft = [
            minutes.toString().padStart(2, '0'),
            seconds.toString().padStart(2, '0')
        ].join(':');
    }

    onMount(() => {
        updateTimeLeft();
        const interval = setInterval(updateTimeLeft, 1000);
        return () => clearInterval(interval);
    });
</script>

<span style="color: {secondsLeft < 10 ? 'red' : 'inherit'}">{timeLeft}</span>
