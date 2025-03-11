<script lang="ts">
	import { Button } from '@/components/ui/button';
	import { featureFlags } from '@/features';
	import { m } from '@/i18n';
	import BaseTopbar from '@/components/layout/topbar.svelte';
	import { currentUser } from '@/pocketbase';
	import UserNav from './userNav.svelte';
	import LanguageSelect from '@/i18n/languageSelect.svelte';
</script>

<BaseTopbar class="bg-card border-none">
	{#snippet left()}
		<!-- <AppLogo /> -->
		<Button variant="link" href={$featureFlags.DEMO ? '#waitlist' : '/'}>
			{m.Getting_started()}
		</Button>
		<Button variant="link" href={$featureFlags.DEMO ? '#waitlist' : '/'}>{m.Tests()}</Button>
		<Button variant="link" href={$featureFlags.DEMO ? '#waitlist' : '/providers'}>
			{m.Services()}
		</Button>
		<Button variant="link" href={$featureFlags.DEMO ? '#waitlist' : '/'}>{m.Apps()}</Button>
		<Button variant="link" href={$featureFlags.DEMO ? '#waitlist' : '/credentials'}>
			{m.Credentials()}
		</Button>
	{/snippet}

	{#snippet right()}
		<div class="flex items-center space-x-2">
			{#if !$featureFlags.DEMO}
				<Button variant="default" href="/tests/new">{m.Start_a_new_check()}</Button>
			{/if}

			{#if !$featureFlags.DEMO && $featureFlags.AUTH}
				{#if $currentUser}
					<UserNav />
				{:else}
					<Button variant="link" href="/login">{m.Login()}</Button>
				{/if}
			{/if}

			{#if !$currentUser}
				<LanguageSelect flagsOnly>
					{#snippet trigger({ triggerAttributes, language })}
						<Button variant="outline" class="w-14 text-2xl" {...triggerAttributes}>
							{language.flag}
						</Button>
					{/snippet}
				</LanguageSelect>
			{/if}
		</div>
	{/snippet}
</BaseTopbar>
