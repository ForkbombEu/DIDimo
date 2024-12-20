<script>
	import { appName } from '@/brand';
	import T from '@/components/ui-custom/t.svelte';
	import { Button } from '@/components/ui/button';
	import { m } from '@/i18n';
</script>

<div class="bg-muted-foreground text-muted">
	<div class="mx-auto max-w-screen-xl">
		<div class="px-8 py-12">
			<div class="flex justify-between">
				<div>
					<T tag="h4">{appName}</T>
					<p class="border-b-muted border-b">{m.Test_and_find_decentralized_IDs()}</p>
				</div>
				<div>
					<Button
						variant="default"
						class="text-muted-foreground bg-white hover:bg-white/90"
						href="/tests/new"
					>
						{m.Start_a_new_test()}
					</Button>
					<Button variant="ghost" class="border" href="/login">{m.Login()}</Button>
				</div>
			</div>
		</div>
	</div>
	<div class="flex items-center justify-center bg-black/30 p-2 px-8">
		<small>Proudly developed by ForkBomb BV</small>
	</div>
</div>
