<script lang="ts">
	import {
		Form,
		createForm,
		createFormData,
		FormError,
		SubmitButton,
		Input,
		Checkbox,
		File as FileInput,
		zodFile
	} from '$lib/forms';

	import { currentUser, pb } from '$lib/pocketbase';
	import { createEventDispatcher } from 'svelte';
	import { z } from 'zod';

	const dispatch = createEventDispatcher<{ success: undefined }>();

	const schema = z.object({
		name: z.string().min(3).optional(),
		email: z.string().email(),
		emailVisibility: z.boolean().optional(),
		avatar: zodFile({ types: ['image/png', 'image/jpeg'], size: 1024 * 1024 * 2 }).optional()
	});

	const initialData: Partial<z.infer<typeof schema>> = {
		name: $currentUser!.name,
		email: $currentUser!.email,
		emailVisibility: $currentUser!.emailVisibility
	};

	const superform = createForm(
		schema,
		async ({ form }) => {
			const formData = createFormData(form.data);
			$currentUser = await pb.collection('users').update($currentUser!.id, formData);
			dispatch('success');
		},
		initialData
	);
</script>

<Form {superform}>
	<Input {superform} field="name" options={{ label: 'Username' }} />

	<div class="space-y-2">
		<Input {superform} field="email" options={{ type: 'email' }} />

		<Checkbox {superform} field="emailVisibility">
			<span>Show email to other users</span>
		</Checkbox>
	</div>

	<FileInput {superform} field="avatar" />

	<FormError />

	<div class="flex justify-end">
		<SubmitButton>Update profile</SubmitButton>
	</div>
</Form>
