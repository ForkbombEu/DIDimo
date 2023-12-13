<script lang="ts">
	import { Heading, P, Button, Hr } from 'flowbite-svelte';
	import { currentUser } from '$lib/pocketbase';
	import UserDataForm from './userDataForm.svelte';
	import { Pencil, XMark } from 'svelte-heros-v2';
	import UserAvatar from '$lib/components/userAvatar.svelte';

	let edit = false;
	const toggleEdit = () => {
		edit = !edit;
	};
</script>

<div class="space-y-6">
	<div class="flex flex-row gap-6 items-center">
		<UserAvatar size="lg" />
		<div class="flex flex-col">
			<Heading tag="h4">{$currentUser?.name}</Heading>
			<P>
				{$currentUser?.email}
				<span class="text-gray-400 text-sm ml-1">
					({$currentUser?.emailVisibility ? 'public' : 'not public'})
				</span>
			</P>
		</div>
	</div>

	<div class="flex items-center gap-4 justify-end">
		{#if edit}
			<Hr />
		{/if}
		<Button color="alternative" on:click={toggleEdit}>
			{#if !edit}
				<Pencil size="20" />
				<span class="ml-2">Edit profile</span>
			{:else}
				<XMark size="20" />
				<span class="ml-2">Cancel</span>
			{/if}
		</Button>
	</div>

	{#if edit}
		<UserDataForm on:success={toggleEdit} />
	{/if}
</div>
