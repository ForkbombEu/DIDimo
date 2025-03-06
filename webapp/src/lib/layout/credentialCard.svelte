<script lang="ts">
	import Avatar from '@/components/ui-custom/avatar.svelte';
	import T from '@/components/ui-custom/t.svelte';
	import { m } from '@/i18n';
	import type { CredentialsResponse } from '@/pocketbase/types';

	type Props = {
		credential: CredentialsResponse;
		class?: string;
	};

	const { credential, class: className = '' }: Props = $props();

	const properties = {
		[m.Issuer()]: credential.issuer_name,
		// [m.Duration()]: 'credential_duration',
		// [m.Specification()]: 'credential_specification',
		// [m.Category()]: 'credential_category',
		[m.Format()]: credential.format,
		[m.Locale()]: credential.locale.toUpperCase()
		// [m.Locale()]: `${emojiFlag(credential.locale.trim())} ${credential.locale.toUpperCase()}`
	};
</script>

<a
	href="/credentials/{credential.id}"
	class="bg-card text-card-foreground border-primary ring-primary flex flex-col gap-6 rounded-xl border p-6 shadow-sm transition-transform hover:-translate-y-2 hover:ring-2 {className}"
>
	<div class="flex items-center gap-2">
		{#if credential.logo}
			<Avatar src={credential.logo} class="!rounded-sm" hideIfLoadingError />
		{/if}
		<T class="font-semibold">{credential.name}</T>
	</div>

	<div class="space-y-1">
		{#each Object.entries(properties) as [key, value]}
			<T class="text-sm text-slate-400">{key}: <span class="text-primary">{value}</span></T>
		{/each}
	</div>
</a>
