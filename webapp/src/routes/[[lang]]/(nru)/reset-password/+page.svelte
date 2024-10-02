<!--
SPDX-FileCopyrightText: 2024 The Forkbomb Company

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import { enhance } from '$app/forms';
	import { Label, Input, Heading, P, Button } from 'flowbite-svelte';

	export let form;
</script>

{#if !form}
	<form method="post" use:enhance class="space-y-8">
		<div class="space-y-1">
			<Heading tag="h4">Reset password</Heading>
			<P>Please enter here a new password.</P>
		</div>
		<Label class="space-y-2">
			<span>New password</span>
			<Input type="password" name="password" id="password" placeholder="•••••••••••" required />
		</Label>
		<Label class="space-y-2">
			<span>New password</span>
			<Input
				type="password"
				name="passwordConfirmation"
				id="passwordConfirmation"
				placeholder="•••••••••••"
				required
			/>
		</Label>
		<Button type="submit" class="w-full">Reset password</Button>
	</form>
{:else if form.success}
	<div class="space-y-4">
		<Heading tag="h4">Password reset successfully!</Heading>
		<P>Please click the link in the email to reset your password.</P>
	</div>
{/if}
