<script>
	import { appName } from '@/brand';
	import T from '@/components/ui-custom/t.svelte';
	import Button from '@/components/ui-custom/button.svelte';
	import { featureFlags } from '@/features';
	import { m } from '@/i18n';
</script>

<div class="bg-primary text-muted">
	<div class="mx-auto max-w-screen-xl">
		<div class="px-8 py-12">
			<div class="flex flex-col justify-between gap-6 sm:flex-row">
				<div>
					<T tag="h4">{appName}</T>
					<p class="border-b border-b-muted">{m.Test_and_find_decentralized_IDs()}</p>
				</div>
				<div class="flex gap-2 blur-sm">
					<Button
						variant="default"
						class="grow basis-1 bg-white text-primary hover:bg-white/90 sm:grow-0"
						href={$featureFlags.DEMO ? undefined : '/tests/new'}
						disabled={$featureFlags.DEMO}
					>
						{m.Start_a_new_test()}
					</Button>
					<Button
						variant="ghost"
						class="grow basis-1 border sm:grow-0"
						href={$featureFlags.DEMO ? undefined : '/login'}
						disabled={$featureFlags.DEMO}
					>
						{m.Login()}
					</Button>
				</div>
			</div>
		</div>
	</div>
	<div class="flex items-center justify-center bg-black/30 p-2 px-8">
		<small>
			<a href="https://forkbomb.solutions/" class="underline hover:no-underline">
				{m.Proudly_developed_by_ForkBomb_BV()}
			</a>
		</small>
	</div>
</div>
