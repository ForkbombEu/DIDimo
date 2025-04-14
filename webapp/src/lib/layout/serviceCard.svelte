<!--
SPDX-FileCopyrightText: 2025 Forkbomb BV

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import T from '@/components/ui-custom/t.svelte';
	import { m, localizeHref } from '@/i18n';
	import { type CredentialIssuersResponse } from '@/pocketbase/types';
	import { String } from 'effect';

	//

	type Props = {
		service: CredentialIssuersResponse;
		class?: string;
	};

	const { service, class: className = '' }: Props = $props();

	const title = $derived(String.isNonEmpty(service.name) ? service.name : service.url);
</script>

<a
	href={localizeHref(`/services/${service.id}`)}
	class="border-primary bg-card text-card-foreground ring-primary block rounded-lg border p-6 shadow-sm transition-all hover:-translate-y-2 hover:ring-2 {className}"
>
	<div class="space-y-4">
		<div class="space-y-1">
			<T tag="small" class="text-primary block">Credential issuer</T>
			<T tag="h4" class="block break-words">{title}</T>
		</div>
		{#if String.isNonEmpty(service.description)}
			<T tag="small" class="block font-normal leading-snug">{service.description}</T>
		{/if}
	</div>
</a>
