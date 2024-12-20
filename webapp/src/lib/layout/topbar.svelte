<script lang="ts">
	import { AppLogo } from '@/brand';
	import { Button } from '@/components/ui/button';
	import { featureFlags } from '@/features';
	import { m } from '@/i18n';
	import BaseTopbar from '@/components/layout/topbar.svelte';
</script>

<BaseTopbar class="bg-background border-none">
	{#snippet left()}
		<!-- <AppLogo /> -->
		<Button variant="link" href="/">{m.Getting_started()}</Button>
		<Button variant="link" href="/">{m.Tests()}</Button>
		<Button variant="link" href="/">{m.Apps()}</Button>
		<Button variant="link" href="/">{m.Credentials()}</Button>
	{/snippet}

	{#snippet right()}
		{#if $featureFlags.AUTH}
			<div>
				<Button variant="default" href="/tests/new">{m.Start_a_new_test()}</Button>
				<Button variant="link" href="/login">{m.Login()}</Button>
			</div>
		{/if}
	{/snippet}
</BaseTopbar>
