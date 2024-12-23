<script lang="ts">
	import BaseLayout from '$lib/layout/baseLayout.svelte';
	import { Button } from '@/components/ui/button';
	import { m } from '@/i18n';
	import type { Snippet } from 'svelte';

	interface Props {
		children?: Snippet;
	}

	let { children }: Props = $props();
</script>

<BaseLayout>
	<div class="border-y">
		<div class="mx-auto max-w-screen-xl px-2 py-2">
			<Button variant="link" href="/my/profile">{m.My_profile()}</Button>
			<Button variant="link" href="/my/organizations">{m.organizations()}</Button>
		</div>
	</div>
	{@render children?.()}
</BaseLayout>
