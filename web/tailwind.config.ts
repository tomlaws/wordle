import type { Config } from 'tailwindcss'

export default {
    darkMode: 'class',
    content: [
        "./src/**/*.{js,ts,svelte,html}",
    ],
    theme: {
        extend: {
            typography: {
                DEFAULT: {
                    css: {
                        color: '#1a202c', // light mode
                    },
                },
                dark: {
                    css: {
                        color: '#e5e7eb', // dark mode
                    },
                },
            },
        },
    },
    plugins: [],
} satisfies Config;