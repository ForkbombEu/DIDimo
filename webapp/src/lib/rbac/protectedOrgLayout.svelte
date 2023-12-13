<script lang="ts">
	import { verifyAuthorizations } from '.';
	import { onNavigate } from '$app/navigation';
	import { browser } from '$app/environment';

	export let orgId: string;

	let verificationPromise: Promise<void>;
	if (browser) verificationPromise = verifyAuthorizations(orgId);

	onNavigate(async () => {
		verificationPromise = verifyAuthorizations(orgId);
	});
</script>

{#await verificationPromise then _}
	<slot />
{:catch _}
	unauthorized
{/await}
