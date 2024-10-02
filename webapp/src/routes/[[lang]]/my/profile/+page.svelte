<!--
SPDX-FileCopyrightText: 2024 The Forkbomb Company

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import UserData from './_partials/userData.svelte';
	import UserKeys from './_partials/userKeys.svelte';
	import UserWebauthn from './_partials/userWebauthn.svelte';

	import { Card } from 'flowbite-svelte';
	import { featureFlags } from '$lib/features';

	export let data;
	let { publicKeys } = data;
</script>

<div class="max-w-xl mx-auto space-y-8 p-8">
	<Card class="!max-w-none">
		<UserData />
	</Card>

	{#if $featureFlags.KEYPAIROOM}
		<Card class="!max-w-none">
			<UserKeys keys={publicKeys} />
		</Card>
	{/if}

	{#if $featureFlags.WEBAUTHN}
		<Card class="!max-w-none">
			<UserWebauthn />
		</Card>
	{/if}
</div>
