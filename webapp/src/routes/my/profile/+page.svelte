<script lang="ts">
	import UserData from './_partials/userData.svelte';
	import UserKeys from './_partials/userKeys.svelte';
	import UserWebauthn from './_partials/userWebauthn.svelte';

	import { featureFlags } from '@/features';
	import Card from '@/components/ui-custom/card.svelte';

	let { data } = $props();
	let { publicKeys } = data;
</script>

<div class="mx-auto max-w-xl space-y-8 p-8">
	<Card>
		<UserData />
	</Card>

	{#if $featureFlags.KEYPAIROOM}
		<Card>
			<UserKeys keys={publicKeys} />
		</Card>
	{/if}

	{#if $featureFlags.WEBAUTHN}
		<Card>
			<UserWebauthn />
		</Card>
	{/if}
</div>
