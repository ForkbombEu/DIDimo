<script lang="ts">
	import {
		registerUser,
		isWebauthnSupported,
		isPlatformAuthenticatorAvailable
	} from '$lib/webauthn';
	import { currentUser } from '$lib/pocketbase';

	import { Collections, type WebauthnCredentialsResponse } from '$lib/pocketbase/types';
	import { CollectionManager, DeleteRecord, EditRecord } from '$lib/collectionManager';

	import { InformationCircle, Plus } from 'svelte-heros-v2';
	import { Alert, Button, Card, Heading, P, Spinner } from 'flowbite-svelte';
	import { createTypeProp } from '$lib/utils/typeProp';

	const platformAuthenticatorAvailable = isPlatformAuthenticatorAvailable();
	const recordType = createTypeProp<WebauthnCredentialsResponse<{ ID: string }>>();

	const userEmailAddress = $currentUser?.email!;
</script>

<Heading tag="h6">Your devices</Heading>
<P color="gray" size="sm">Manage the devices you use to login.</P>

<CollectionManager
	{recordType}
	collection={Collections.WebauthnCredentials}
	let:records
	editFormSettings={{ exclude: ['user', 'credential'] }}
>
	<div class="space-y-2 py-4">
		{#each records as record}
			{@const label = Boolean(record.description) ? record.description : record.credential?.ID}
			<Card class="max-w-none">
				<div class="flex justify-between gap-4 items-center">
					<div class="w-0 grow overflow-hidden">
						<P>{label}</P>
					</div>
					<div class="flex gap-2">
						<EditRecord {record} />
						<DeleteRecord {record} />
					</div>
				</div>
			</Card>
		{/each}
	</div>
</CollectionManager>

{#await platformAuthenticatorAvailable}
	<div class="flex flex-col items-center">
		<Spinner />
		<P>Checking your device</P>
	</div>
{:then response}
	{#if !response || !isWebauthnSupported()}
		<Alert color="yellow" class="mt-2">
			<svelte:fragment slot="icon">
				<InformationCircle />
			</svelte:fragment>
			<span>
				Your device does not have integrated Webauthn support. An external authenticator is
				required.
			</span>
		</Alert>
	{/if}
{/await}

<div class="flex justify-end mt-2">
	<Button
		color="alternative"
		on:click={() => {
			registerUser(userEmailAddress, navigator.userAgent);
		}}
	>
		<Plus size="20" class="mr-1" /> Add a device
	</Button>
</div>
