<script lang="ts">
    const { disabled, type, outline, children, ...props }: {
        disabled?: boolean;
        type?: "button" | "submit" | "reset" | null | undefined;
        outline?: boolean;
        children: () => any;
        [key: string]: any;
    } = $props();
</script>

<button
    class="game-btn"
    class:outline={outline}
    {disabled}
    type={type}
    {...props}
>
    {@render children()}
</button>

<style>
.game-btn {
    padding: 0.7em 1.8em;
    background: #3b82f6; /* blue */
    color: #fff;
    border: none;
    border-radius: 0.5em;
    font-size: 1rem;
    font-family: inherit;
    font-weight: 500;
    letter-spacing: 0.03em;
    cursor: pointer;
    transition: background 0.18s, transform 0.09s, box-shadow 0.18s;
    outline: none;
    display: inline-flex;
    align-items: center;
    gap: 0.4em;
    text-transform: uppercase;
}

/* Use a slightly lighter and darker shade of the base color for hover and active */
.game-btn:hover:not(:disabled) {
    background: #60a5fa; /* lighter blue */
    transform: translateY(-1px) scale(1.02);
}

.game-btn:active:not(:disabled) {
    background: #1d4ed8; /* deeper blue */
    transform: scale(0.97);
}

.game-btn:disabled {
    background: #cbd5e1;
    color: #64748b;
    cursor: not-allowed;
    opacity: 0.7;
}

.game-btn.outline {
    background: transparent;
    border: 2px solid #3b82f6;
    color: #3b82f6;
}
.game-btn.outline:hover:not(:disabled) {
    background: rgba(59, 130, 246, 0.1);
}
.game-btn.outline:active:not(:disabled) {
    background: rgba(59, 130, 246, 0.2);
    border-color: #2563eb;
}
</style>