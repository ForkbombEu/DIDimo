<!--
SPDX-FileCopyrightText: 2024 The Forkbomb Company

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import FieldError from './fieldError.svelte';
	import FieldHelpText from './fieldHelpText.svelte';
	import FieldLabel from './fieldLabel.svelte';

	export let field: string;
	export let label: string | null | undefined = undefined;
	export let helpText: string | undefined = undefined;
</script>

<div class="space-y-2">
	<FieldLabel {field} text={label}>
		<slot slot="right" name="labelRight" />
	</FieldLabel>
	<slot />
	<FieldHelpText text={helpText} />
	<FieldError {field} />
</div>
