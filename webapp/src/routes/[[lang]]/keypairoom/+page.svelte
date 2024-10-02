<!--
SPDX-FileCopyrightText: 2024 The Forkbomb Company

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import {
		userChallengesSchema,
		userChallenges,
		type UserChallenges
	} from '$lib/keypairoom/userQuestions.js';
	import {
		generateKeypair,
		getHMAC,
		matchPublicAndPrivateKeys,
		saveKeyringToLocalStorage
	} from '$lib/keypairoom/keypair';
	import {
		getPublicKeysFromKeypair,
		saveUserPublicKeys,
		getUserPublicKeys
	} from '$lib/keypairoom/utils.js';
	import { currentUser, pb } from '$lib/pocketbase';
	import { z } from 'zod';
	import { featureFlags } from '$lib/features';
	import { page } from '$app/stores';
	import { welcomeSearchParamKey } from '$lib/utils/constants.js';
	import { appTitle } from '$lib/strings.js';

	// Components
	import { Form, createForm, Input, FormError, SubmitButton } from '$lib/forms';
	import { A, Alert, Button, Heading, Hr, P } from 'flowbite-svelte';
	import CopyButton from '$lib/components/copyButton.svelte';
	import Card from '$lib/components/card.svelte';
	import { InformationCircle } from 'svelte-heros-v2';

	//

	let seed = '';

	const schema = z.object({
		email: z.string().email(),
		questions: userChallengesSchema
	});

	const superform = createForm(schema, async ({ form }) => {
		const { email, questions } = form.data;
		const challenges = userChallengesSchema.parse(questions);
		const keypair = await createKeypairFromFormData(email, challenges);

		const privateKeys = keypair.keyring;
		const publicKeys = getPublicKeysFromKeypair(keypair);

		if ($featureFlags.AUTH && $currentUser) {
			const storedPublicKeys = await getUserPublicKeys();

			if (!storedPublicKeys) {
				await saveUserPublicKeys(publicKeys);
			} else {
				try {
					await matchPublicAndPrivateKeys(storedPublicKeys, privateKeys);
				} catch (e) {
					throw new Error('Wrong answers');
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

		saveKeyringToLocalStorage(privateKeys);
		seed = keypair.seed;
	});

	async function createKeypairFromFormData(email: string, challenges: UserChallenges) {
		const HMAC = await getHMAC(email);
		return await generateKeypair(email, HMAC, challenges);
	}

	const { form } = superform;

	if ($currentUser) $form.email = $currentUser.email;

	//

	$: isWelcome = $page.url.searchParams.has(welcomeSearchParamKey);
	$: searchParams = $page.url.searchParams.toString();
</script>

{#if isWelcome}
	<div class="mb-6 -rotate-1">
		<Alert color="yellow" border>
			<div class="space-y-3 overflow-hidden text-ellipsis">
				<Heading color="yellow" tag="h2" class="text-ellipsis">Welcome to {appTitle} 🎉</Heading>
				<P color="yellow" weight="bold">Thanks for joining us!</P>
				<P color="yellow">
					One last thing before to using the app:<br /> we need you to answer these questions, as they
					will be used to secure your data.
				</P>
			</div>
		</Alert>
	</div>
{/if}

<Card class="space-y-6 p-6">
	{#if !seed}
		<Heading tag="h4">Generate your keys</Heading>

		<Alert color="blue">
			<span class="sr-only">Info</span>
			<span class="text mb-2 flex items-center font-bold">
				<div class="mr-1">
					<InformationCircle size="20" />
				</div>
				Important information
			</span>
			<ul class="list-disc space-y-1 pl-4 pt-1">
				<li>
					By answering these questions, you will generate keys that will be used to encrypt your
					data
				</li>
				<li>
					Please remember the answers, as they will be the only way to restore the encryption keys
				</li>
				<li>Please answer at least 3 of the following questions</li>
			</ul>
		</Alert>

		<Hr />

		<Form {superform} className="space-y-6">
			{#if !$currentUser}
				<div class="space-y-1">
					<Input {superform} field="email" options={{ label: 'User email' }} />

					<P size="sm" color="text-gray-400">
						Your email won't be stored anywhere, it will be used only to generate the keys.
					</P>
				</div>

				<Hr />
			{/if}

			{#each userChallenges as question}
				<Input {superform} field={`questions.${question.id}`} options={{ label: question.text }} />
			{/each}

			<FormError />

			<div class="flex justify-end">
				<SubmitButton>Generate keys</SubmitButton>
			</div>
		</Form>

		<Hr />

		<A class="text-sm" href="/keypairoom/regenerate">I have the seed passphrase</A>
	{:else}
		<Heading tag="h4">Keypair creation successful!</Heading>
		<P size="sm" color="text-gray-400 dark:text-gray-600">
			Please store this in a safe place to recover your account in the future, this passphrase will
			be shown only one time!
		</P>
		<Alert color="blue">
			<span class="font-mono">
				{seed}
				<div class="flex flex-col items-end pt-4">
					<CopyButton textToCopy={seed}>Copy seed</CopyButton>
				</div>
			</span>
		</Alert>
		<Button href={`/my?${searchParams}`}>Go to Dashboard</Button>
	{/if}
</Card>
