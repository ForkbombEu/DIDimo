<script lang="ts">
	import { i18n } from '@/i18n';
	import { ParaglideJS } from '@inlang/paraglide-sveltekit';
	import '../app.css';
	import { appName } from '@/brand';
	import { Toaster } from '@/components/ui/sonner';

	let { children } = $props();
</script>

<svelte:head>
	<title>
		{appName}
	</title>
</svelte:head>

<ParaglideJS {i18n}>
	{@render children()}
</ParaglideJS>

<Toaster richColors closeButton />
