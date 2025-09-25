import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vitest/config';

export default defineConfig({
	resolve: process.env.VITEST
		? {
			conditions: ['browser']
		}
		: undefined,
	plugins: [tailwindcss(), sveltekit()],
	test: {
		globals: true,
		environment: 'jsdom',
		setupFiles: ['./src/setupTest.ts'],
	}
});
