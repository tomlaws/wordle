<script lang="ts">
    import type { HTMLButtonAttributes } from 'svelte/elements';
    let { children, ...props }: HTMLButtonAttributes = $props();
</script>

<button
    class="game-btn"
    {...props}
>
    {@render children?.()}
</button>

<style>
.game-btn {
    padding: 0.7em 1.8em;
    background: var(--color-btn);
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
    letter-spacing: 0.2em;
}

/* Use a slightly lighter and darker shade of the base color for hover and active */
.game-btn:hover:not(:disabled) {
    background: var(--color-btn-hover);
    transform: translateY(-1px) scale(1.02);
}

.game-btn:active:not(:disabled) {
    background: var(--color-btn-active);
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
    border: 2px solid var(--color-btn);
    color: var(--color-btn);
}
.game-btn.outline:hover:not(:disabled) {
    background: rgba(59, 130, 246, 0.1);
}
.game-btn.outline:active:not(:disabled) {
    background: rgba(59, 130, 246, 0.2);
    border-color: var(--color-btn-active);
}
</style>