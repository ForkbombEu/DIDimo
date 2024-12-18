<script lang="ts">
	import {
		getHMAC,
		matchPublicAndPrivateKeys,
		regenerateKeypair,
		saveKeyringToLocalStorage,
		type Keyring
	} from '@/keypairoom/keypair';
	import { currentUser, pb } from '@/pocketbase';

	import z from 'zod';
	import { Form, createForm } from '@/forms';
	import { Field, TextareaField } from '@/forms/fields';
	import Card from '@/components/ui-custom/card.svelte';
	import { featureFlags } from '@/features';
	import { zod } from 'sveltekit-superforms/adapters';
	import { getUserPublicKeys, RegenerateKeyringSession } from '@/keypairoom/utils';
	import { m } from '@/i18n';
	import RegenerateBanner from '../_partials/RegenerateBanner.svelte';
	import { log } from '@/utils/other';
	import T from '@/components/ui-custom/t.svelte';
	import Separator from '@/components/ui/separator/separator.svelte';
	import { PageCard } from '@/components/layout';
	import A from '@/components/ui-custom/a.svelte';

	//

	let success = $state(false);

	const schema = z.object({
		email: z.string().email(),
		seed: z.string()
	});

	const form = createForm({
		adapter: zod(schema),
		onSubmit: async ({ form }) => {
			const hmac = await getHMAC(form.data.email);

			let keyring: Keyring;
			try {
				const keypair = await regenerateKeypair(form.data.seed, hmac);
				keyring = keypair.keyring;
			} catch (e) {
				throw new Error(m.Invalid_seed());
			}

			if ($featureFlags.AUTH && $currentUser) {
				const publicKeys = await getUserPublicKeys();
				if (!publicKeys) {
					throw new Error(
						m.User_public_keys_are_missing_Please_generate_them_using_the_security_questions_()
					);
				} else {
					try {
						await matchPublicAndPrivateKeys(publicKeys, keyring);
					} catch (e) {
						throw new Error(m.Invalid_seed());
					}
				}

				if ($featureFlags.DID) {
					try {
						await pb.send('/api/did', {});
					} catch (e) {
						log(e);
					}
				}
			}

			saveKeyringToLocalStorage(keyring);
			success = true;

			if (RegenerateKeyringSession.isActive()) RegenerateKeyringSession.end();
		}
	});

	const { form: formData } = form;
	if ($currentUser) $formData.email = $currentUser.email;

	const textAreaPlaceholder =
		'skin buyer sunset person run push elevator under debris soft surge man';
</script>

{#if !success}
	<T tag="h4">{m.Regenerate_keys()}</T>

	{#if RegenerateKeyringSession.isActive()}
		<RegenerateBanner />
	{/if}

	{#if $currentUser}
		<T>{m.Please_type_here_your_seed_to_restore_your_keyring_()}</T>
	{:else}
		<T>{m.Please_type_here_your_email_and_your_seed_to_restore_your_keyring_()}</T>
	{/if}

	<Form {form} hideRequiredIndicator>
		{#if !$currentUser}
			<div class="space-y-1">
				<Field {form} name="email" options={{ label: m.User_email() }} />
				<T tag="small" class="text-gray-400">
					{m.Your_email_wont_be_stored_anywhere_it_will_be_used_only_to_generate_the_keys_()}
				</T>
			</div>
		{/if}

		<TextareaField {form} name="seed" options={{ placeholder: textAreaPlaceholder }} />

		{#snippet submitButtonContent()}
			{m.Regenerate_keys()}
		{/snippet}
	</Form>

	<Separator />

	<A href="/keypairoom" class="block text-sm">{m.Forgot_the_seed_Regenerate_it()}</A>
{:else}
	<div class="flex flex-col space-y-4 p-6">
		<T tag="h4">{m.Keys_regenerated()}</T>
		<T>
			{m.Your_keys_have_been_regenerated_You_can_now_go_back_to()}
			<A href="/my">{m.your_profile()}</A>
		</T>
	</div>
{/if}
