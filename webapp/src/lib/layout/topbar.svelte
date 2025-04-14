<!--
SPDX-FileCopyrightText: 2025 Forkbomb BV

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import Button from '@/components/ui-custom/button.svelte';
	import { featureFlags } from '@/features';
	import { m } from '@/i18n';
	import BaseTopbar from '@/components/layout/topbar.svelte';
	import { currentUser } from '@/pocketbase';
	import UserNav from './userNav.svelte';
	import LanguageSelect from '@/i18n/languageSelect.svelte';
	import Icon from '@/components/ui-custom/icon.svelte';
	import { Sparkle } from 'lucide-svelte';
</script>

<BaseTopbar class="border-none bg-card">
	{#snippet left()}
		<!-- <AppLogo /> -->
		<Button variant="link" href={$featureFlags.DEMO ? '#waitlist' : '/'}>
			{m.Getting_started()}
		</Button>
		<!-- <Button variant="link" href={$featureFlags.DEMO ? '#waitlist' : '/'}>{m.Tests()}</Button> -->
		<Button variant="link" href={$featureFlags.DEMO ? '#waitlist' : '/services'}>
			{m.Services()}
		</Button>
		<Button variant="link" href={$featureFlags.DEMO ? '#waitlist' : '/apps'}>{m.Apps()}</Button>
		<Button variant="link" href={$featureFlags.DEMO ? '#waitlist' : '/credentials'}>
			{m.Credentials()}
		</Button>
		<Button variant="link" href={$featureFlags.DEMO ? '#waitlist' : '/organizations'}>
			{m.organizations()}
		</Button>
	{/snippet}

	{#snippet right()}
		<div class="flex items-center space-x-2">
			{#if !$featureFlags.DEMO && $featureFlags.AUTH}
				{#if $currentUser}
					<Button href="/my/tests/new">
						<Icon src={Sparkle} />
						{m.Start_a_new_test()}
					</Button>

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
