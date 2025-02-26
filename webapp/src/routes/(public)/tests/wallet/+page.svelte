<script lang="ts">
	import PageContent from '$lib/layout/pageContent.svelte';
	import T from '@/components/ui-custom/t.svelte';
	import Button from '@/components/ui/button/button.svelte';
	import TextareaField from '@/forms/fields/textareaField.svelte';
	import { Form, SubmitButton, createForm } from '@/forms';
	import { QrCode } from '@/qr';
	import { zod } from 'sveltekit-superforms/adapters';
	import { z } from 'zod';
	import Separator from '@/components/ui/separator/separator.svelte';
	import { pb } from '@/pocketbase/index.js';
	import Alert from '@/components/ui-custom/alert.svelte';

	let { data } = $props();

	let result = $state();

	const successForm = createForm({
		adapter: zod(z.object({})),
		onSubmit: async () => {
			const result = await pb.send('/wallet-test/confirm-success', {
				method: 'POST',
				body: {
					workflow_id: data.workflowId
				}
			});
			console.log('success', result);
		}
	});

	const failureForm = createForm({
		adapter: zod(z.object({ reason: z.string() })),
		onSubmit: async ({
			form: {
				data: { reason }
			}
		}) => {
			const result = await pb.send('/wallet-test/notify-failure', {
				method: 'POST',
				body: {
					workflow_id: data.workflowId,
					reason: reason
				}
			});
			console.log('failure', result);
		}
	});
</script>

<PageContent>
	<div class="space-y-4">
		<T tag="h1">Wallet test</T>

		<div class="step-container">
			{@render Step(1, 'Scan this QR with the wallet app to start the check')}

			<div class="bg-primary/10 ml-16 mt-4 flex items-center justify-center rounded-md p-2">
				<QrCode src={data.qrContent} class="size-40 rounded-sm" />
			</div>
		</div>

		<div class="step-container">
			{@render Step(2, 'Follow the procedure on the wallet app')}
		</div>

		<div class="step-container">
			{@render Step(3, 'Confirm the result')}

			{#if !result}
				<div class="ml-16 flex flex-col gap-8 sm:flex-row">
					<div class="grow basis-1">
						<Form form={successForm}>
							{#snippet submitButton()}
								<SubmitButton class="w-full bg-green-600 hover:bg-green-700">
									All good!
								</SubmitButton>
							{/snippet}
						</Form>
					</div>

					<div class="grow basis-1">
						<Form form={failureForm} hideRequiredIndicator>
							<TextareaField form={failureForm} name="reason" />
							{#snippet submitButton()}
								<SubmitButton class="w-full bg-red-600 hover:bg-red-700">
									There was an issue!
								</SubmitButton>
							{/snippet}
						</Form>
					</div>
				</div>
			{:else}
				<Alert variant="info">Your response was submitted!</Alert>
			{/if}
		</div>
	</div>
</PageContent>

{#snippet Step(n: number, text: string)}
	<div class="flex items-center gap-4">
		<div
			class="bg-primary text-primary-foreground flex size-12 shrink-0 items-center justify-center rounded-full text-lg font-semibold"
		>
			<p>{n}</p>
		</div>
		<T class="text-primary font-semibold">{text}</T>
	</div>
{/snippet}

<style lang="postcss">
	.step-container {
		@apply bg-secondary rounded-xl p-4;
	}
</style>
