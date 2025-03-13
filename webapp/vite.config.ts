import { paraglide } from '@inlang/paraglide-sveltekit/vite';
import { defineConfig } from 'vitest/config';
import { sveltekit } from '@sveltejs/kit/vite';

export default defineConfig({
	plugins: [
		sveltekit(),
		paraglide({
			project: './project.inlang',
			outdir: './src/modules/i18n/paraglide'
		})
	],

	optimizeDeps: {
		exclude: [
			'svelte-codemirror-editor',
			'codemirror',
			'@codemirror/language-javascript',
			'@codemirror/lang-json',
			'thememirror'
		],
		include: ['date-fns', 'date-fns-tz']
	},

	test: {
		include: ['src/**/*.{test,spec}.{js,ts}']
	},

	server: {
		port: Number(process.env.PORT) || 5100
	},
	preview: {
		allowedHosts: true
	}
});
