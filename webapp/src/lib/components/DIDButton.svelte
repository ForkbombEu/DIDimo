<script>
	import { featureFlags } from '$lib/features';
	import { currentUser } from '$lib/pocketbase';
	import { Button } from 'flowbite-svelte';

	$: eddsa = $currentUser?.eddsa_public_key;
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
