<!--
SPDX-FileCopyrightText: 2024 The Forkbomb Company

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import {
		getHMAC,
		matchPublicAndPrivateKeys,
		regenerateKeypair,
		saveKeyringToLocalStorage,
		type Keyring
	} from '$lib/keypairoom/keypair';
	import { currentUser, pb } from '$lib/pocketbase';

	import { z } from 'zod';
	import { Form, createForm, FormError, SubmitButton, Textarea, Input } from '$lib/forms';
	import { A, Alert, Heading, Hr, P } from 'flowbite-svelte';
	import Card from '$lib/components/card.svelte';
	import { page } from '$app/stores';
	import { missingKeyringParamKey } from '$lib/utils/constants.js';
	import { ExclamationTriangle } from 'svelte-heros-v2';
	import { featureFlags } from '$lib/features';
	import { getUserPublicKeys } from '$lib/keypairoom/utils';

	//

	let success = false;

	const schema = z.object({
		email: z.string().email(),
		seed: z.string()
	});

	const superform = createForm(schema, async ({ form }) => {
		const hmac = await getHMAC(form.data.email);
		let keyring: Keyring;
		try {
			const keypair = await regenerateKeypair(form.data.seed, hmac);
			keyring = keypair.keyring;
		} catch (e) {
			throw new Error('Invalid seed');
		}

		if ($featureFlags.AUTH && $currentUser) {
			const publicKeys = await getUserPublicKeys();
			if (!publicKeys) {
				throw new Error(
					'User public keys are missing. Please generate them using the security questions.'
				);
			} else {
				try {
					await matchPublicAndPrivateKeys(publicKeys, keyring);
				} catch (e) {
					throw new Error('Invalid seed');
				}
			}

			if ($featureFlags.DID) {
				try {
					await pb.send('/api/did', {});
				} catch (e) {
					console.log(e);
				}
			}
		}

		saveKeyringToLocalStorage(keyring);
		success = true;
	});

	const { form } = superform;

	if ($currentUser) $form.email = $currentUser.email;

	const textAreaPlaceholder =
		'skin buyer sunset person run push elevator under debris soft surge man';

	$: isKeyringMissing = $page.url.searchParams.has(missingKeyringParamKey);
</script>

<Card class="p-6 space-y-6">
	{#if !success}
		<Heading tag="h4">Regenerate keys</Heading>

		{#if isKeyringMissing}
			<Alert color="yellow" border>
				<svelte:fragment slot="icon"><ExclamationTriangle /></svelte:fragment>
				<div class="space-y-1">
					<p>You have been redirected here because your private keys are missing.</p>
					<p>Before using the app again, you need to restore them.</p>
				</div>
			</Alert>
		{/if}

		<Hr />

		{#if $currentUser}
			<P>Please type here your seed to restore your keyring.</P>
		{:else}
			<P>Please type here your email and your seed to restore your keyring.</P>
		{/if}

		<Form {superform}>
			{#if !$currentUser}
				<div class="space-y-1">
					<Input {superform} field="email" options={{ label: 'User email' }} />
					<P size="sm" color="text-gray-400">
						Your email won't be stored anywhere, it will be used only to generate the keys.
					</P>
				</div>
			{/if}

			<Textarea {superform} field="seed" options={{ placeholder: textAreaPlaceholder }} />

			<FormError />

			<div class="flex justify-end">
				<SubmitButton>Regenerate keys</SubmitButton>
			</div>
		</Form>

		<Hr />

		<A href="/keypairoom" class="text-sm">Forgot the seed? Regenerate it</A>
	{:else}
		<div class="space-y-4 p-6 flex flex-col">
			<Heading tag="h4">Keys regenerated!</Heading>
			<P>
				Your keys have been regenerated. You can now go back to
				<A href="/my">your profile</A>.
			</P>
		</div>
	{/if}
</Card>
