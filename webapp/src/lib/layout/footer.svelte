<script>
	import { appName } from '@/brand';
	import T from '@/components/ui-custom/t.svelte';
	import { Button } from '@/components/ui/button';
	import { m } from '@/i18n';
</script>

<div class="bg-muted-foreground text-muted">
	<div class="mx-auto max-w-screen-xl">
		<div class="px-8 py-12">
			<div class="flex flex-col justify-between gap-6 sm:flex-row">
				<div>
					<T tag="h4">{appName}</T>
					<p class="border-b-muted border-b">{m.Test_and_find_decentralized_IDs()}</p>
				</div>
				<div class="flex gap-2">
					<Button
						variant="default"
						class="text-muted-foreground grow basis-1 bg-white hover:bg-white/90 sm:grow-0"
						href="/tests/new"
					>
						{m.Start_a_new_test()}
					</Button>
					<Button variant="ghost" class="grow basis-1 border sm:grow-0" href="/login"
						>{m.Login()}</Button
					>
				</div>
			</div>
		</div>
	</div>
	<div class="flex items-center justify-center bg-black/30 p-2 px-8">
		<small>{m.Proudly_developed_by_ForkBomb_BV()}</small>
	</div>
</div>
