<script lang="ts">
	import { getUserDidUrl } from '@/did';
	import { featureFlags } from '@/features';
	import { Button } from '@/components/ui/button';
	import { m } from '@/i18n';

	let url: string | undefined = $state(undefined);

	getUserDidUrl().then((res) => {
		url = res;
	});
</script>

{#if $featureFlags.DID && url}
	<Button href={url} target="_blank" size="sm" variant="outline">
		{m.my_DID()}
	</Button>
{/if}
