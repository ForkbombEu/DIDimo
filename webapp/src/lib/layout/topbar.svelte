<script lang="ts">
	import { Button } from '@/components/ui/button';
	import { featureFlags } from '@/features';
	import { m } from '@/i18n';
	import BaseTopbar from '@/components/layout/topbar.svelte';
	import { currentUser } from '@/pocketbase';
	import UserNav from './userNav.svelte';
</script>

<BaseTopbar class="border-none bg-card">
	{#snippet left()}
		<!-- <AppLogo /> -->
		<Button variant="link" href={$featureFlags.DEMO ? '#waitlist' : '/'}
			>{m.Getting_started()}</Button
		>
		<Button variant="link" href={$featureFlags.DEMO ? '#waitlist' : '/'}>{m.Tests()}</Button>
		<Button variant="link" href={$featureFlags.DEMO ? '#waitlist' : '/'}>{m.Services()}</Button>
		<Button variant="link" href={$featureFlags.DEMO ? '#waitlist' : '/'}>{m.Apps()}</Button>
		<Button variant="link" href={$featureFlags.DEMO ? '#waitlist' : '/'}
			>{m.Credentials()}</Button
		>
	{/snippet}

	{#snippet right()}
		{#if !$featureFlags.DEMO}
			{#if $featureFlags.AUTH}
				<div class="space-x-2">
					<Button variant="default" href="/tests/new">{m.Start_a_new_check()}</Button>
					{#if $currentUser}
						<UserNav />
					{:else}
						<Button variant="link" href="/login">{m.Login()}</Button>
					{/if}
				</div>
			{/if}
		{/if}
	{/snippet}
</BaseTopbar>
