<script lang="ts">
	import { featureFlags } from '@/features';
	import { WelcomeSession } from '@/auth/welcome';
	import A from '@/components/ui-custom/a.svelte';

	if (WelcomeSession.isActive()) WelcomeSession.end();
</script>

<div class="flex flex-col p-4">
	<A href="/my/profile">My profile</A>
	{#if $featureFlags.ORGANIZATIONS}
		<A href="/my/organizations">My organizations</A>
	{/if}
</div>
