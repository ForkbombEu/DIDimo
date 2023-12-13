<script lang="ts">
	import { Register } from 'flowbite-svelte-blocks';
	import { appTitle } from '$lib/strings';
</script>

<div class="grid h-screen place-items-center">
	<div class="w-2/5">
		<Register href="/">
			<svelte:fragment slot="top">
				<img class="h-10 mr-2" src="/logo.svg" alt={`${appTitle} Logo`} />
			</svelte:fragment>
			<div class="p-6 space-y-4 md:space-y-6 sm:p-8">
				<slot />
			</div>
		</Register>
	</div>
</div>
