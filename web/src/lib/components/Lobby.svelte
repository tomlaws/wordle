<script>
	export let nickname = '';
	export let message = 'Finding opponents...';
</script>
<div class="flex flex-col items-center justify-center min-h-screen text-[1.2rem] text-[#333]">
    <p class="mb-4">Welcome to Wordle, {nickname}!</p>
    <div class="spinner mb-4"></div>
    <p>{message}</p>
</div>

<style>
	.spinner {
		width: 48px;
		height: 48px;
		border: 6px solid #eee;
		border-top: 6px solid #3498db;
		border-radius: 50%;
		animation: spin 1s linear infinite;
		margin-bottom: 16px;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}
</style>
