<!--
SPDX-FileCopyrightText: 2025 Forkbomb BV

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import AuthLayout from '@/auth/authLayout.svelte';
	import LanguageSelect from '@/i18n/languageSelect.svelte';

	import { m } from '@/i18n';
	import { currentUser } from '@/pocketbase';
	import T from '@/components/ui-custom/t.svelte';
	import { Button } from '@/components/ui/button/index.js';

	export let data;
</script>

<AuthLayout>
	{#snippet topbarRight()}
		<LanguageSelect />
	{/snippet}

	{#if data.error}
		<div class="space-y-4">
			<T tag="h4">{m.Oh_no()}</T>
			<T>{m.An_error_occurred_while_verifying_your_email_()}</T>
		</div>
	{:else if data.verified}
		<div class="space-y-4">
			<T tag="h4">{m.Email_verified_succesfully()}</T>
			{#if !currentUser}
				<Button href="/login" class="w-full">{m.Go_to_login()}</Button>
			{:else}
				<Button href="/my" class="w-full">{m.Go_to_Dashboard()}</Button>
			{/if}
		</div>
	{/if}
</AuthLayout>
