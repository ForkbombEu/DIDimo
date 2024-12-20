<script lang="ts">
	import { AppLogo } from '@/brand';
	import { Button } from '@/components/ui/button';
	import { featureFlags } from '@/features';
	import { m } from '@/i18n';
</script>

<div class="flex items-center justify-between border-b p-2">
	<AppLogo />

	{#if $featureFlags.AUTH}
		<div>
			<Button variant="default" href="/tests/new">{m.Start_a_new_test()}</Button>
			<Button variant="link" href="/login">{m.Login()}</Button>
		</div>
	{/if}
</div>
