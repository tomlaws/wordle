import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vitest/config';
import adapter from '@sveltejs/adapter-static';

export default defineConfig({
	envPrefix: ['VITE_', 'PUBLIC_'],
	resolve: process.env.VITEST
		? {
			conditions: ['browser']
		}
		: undefined,
	plugins: [tailwindcss(), sveltekit({
		alias: { '$lib': 'src/lib' },
		adapter: adapter({ pages: 'build', assets: 'build', strict: true })
	})],
	test: {
		globals: true,
		environment: 'jsdom',
		setupFiles: ['./src/setupTest.ts'],
	}
});
