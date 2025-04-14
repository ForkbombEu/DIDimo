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
	import { AppLogo } from '@/brand';
</script>

<BaseTopbar class="border-none bg-card">
	{#snippet left()}
		<Button variant="link" href={$featureFlags.DEMO ? '#waitlist' : '/'}>
			<AppLogo />
		</Button>
		<div class="hidden flex-row sm:flex">
			<Button variant="link" href={$featureFlags.DEMO ? '#waitlist' : '/credentials'}>
				{m.Credentials()}
			</Button>
			<Button variant="link" href={$featureFlags.DEMO ? '#waitlist' : '/apps'}
				>{m.Apps()}</Button
			>
			<Button variant="link" href={$featureFlags.DEMO ? '#waitlist' : '/services'}>
				{m.Issuers()}
			</Button>
			<Button variant="link" href={$featureFlags.DEMO ? '#waitlist' : '/organizations'}>
				{m.organizations()}
			</Button>
		</div>
	{/snippet}

	{#snippet right()}
		<div class="flex items-center space-x-2">
			<Button variant="link" href="/news">{m.News()}</Button>
			<Button variant="link" href="/help">{m.Help()}</Button>
			{#if !$featureFlags.DEMO && $featureFlags.AUTH}
				{#if !$currentUser}
					<Button variant="secondary" href="/login">{m.Login()}</Button>
				{:else}
					<Button variant="link" href="/my/tests/new">
						<Icon src={Sparkle} />
						{m.Start_a_new_test()}
					</Button>
					<UserNav />
				{/if}
			{/if}
		</div>
	{/snippet}
</BaseTopbar>
