<script lang="ts">
	import { featureFlags } from '$lib/features';
	import { getUserPublicKeys } from '$lib/keypairoom/utils';
	import { Button } from 'flowbite-svelte';

	let eddsa: string | undefined = undefined;

	getUserPublicKeys().then((res) => {
		eddsa = res?.eddsa_public_key;
	});
</script>

{#if $featureFlags.DID && eddsa}
	<Button
		href="https://explorer.did.dyne.org/details/did:dyne:sandbox.signroom:{eddsa}"
		target="_blank"
		size="xs"
		class="ml-3"
		color="light">My DID</Button
	>
{/if}
