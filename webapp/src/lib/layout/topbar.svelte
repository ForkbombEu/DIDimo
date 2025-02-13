<script lang="ts">
	import { Button } from '@/components/ui/button';
	import { featureFlags } from '@/features';
	import { m } from '@/i18n';
	import BaseTopbar from '@/components/layout/topbar.svelte';
	import { currentUser } from '@/pocketbase';
	import UserNav from './userNav.svelte';
</script>

<BaseTopbar class="bg-card border-none">
	{#snippet left()}
		<!-- <AppLogo /> -->
		<Button variant="link" href="/">{m.Getting_started()}</Button>
		<Button variant="link" href="/">{m.Tests()}</Button>
		<Button variant="link" href="/providers">{m.Services()}</Button>
		<Button variant="link" href="/">{m.Apps()}</Button>
		<Button variant="link" href="/credentials">{m.Credentials()}</Button>
	{/snippet}

	{#snippet right()}
		<div class="space-x-2">
			<Button variant="default" href="/tests/new">{m.Start_a_new_check()}</Button>
			{#if $featureFlags.AUTH}
				{#if $currentUser}
					<UserNav />
				{:else}
					<Button variant="link" href="/login">{m.Login()}</Button>
				{/if}
			{/if}
		</div>
	{/snippet}
</BaseTopbar>
