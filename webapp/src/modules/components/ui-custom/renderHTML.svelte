<!--
SPDX-FileCopyrightText: 2025 Forkbomb BV

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import DOMPurify from 'dompurify';

	type Props = {
		content: string;
		className?: string;
	};

	const { content, className }: Props = $props();
</script>

<div class={className}>
	{@html DOMPurify.sanitize(content)}
</div>
