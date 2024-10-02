<!--
SPDX-FileCopyrightText: 2024 The Forkbomb Company

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import { ParaglideJS } from '@inlang/paraglide-sveltekit';
	import { i18n } from '$lib/i18n';

	import '../app.postcss';
	import { version } from '$app/environment';
	import { appTitle } from '$lib/strings';

	console.info(`${appTitle} version: 🔖 ${version}`);
</script>

<ParaglideJS {i18n}>
	<div class="p-8">
		<slot />
	</div>
</ParaglideJS>
