<script>
	import { Logo, Topbar } from '$lib/layout';
	import { currentUser } from '$lib/pocketbase';
	import DIDButton from '$lib/components/DIDButton.svelte';
	import UserAvatar from '$lib/components/userAvatar.svelte';
</script>

<Topbar>
	<Logo slot="left" />
	<div class="flex items-center" slot="center">
		<div>
			<span class="whitespace-nowrap">
				Hello, <span class="font-semibold">{$currentUser?.email}</span>
			</span>
		</div>
		<div class="shrink-0">
			<DIDButton />
		</div>
	</div>
	<UserAvatar slot="right" />
</Topbar>
<div class="mx-auto max-w-md p-6">
	<slot />
</div>
