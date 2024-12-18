<script lang="ts">
	import Button, { type ButtonProps } from '@/components/ui/button/button.svelte';
	import { ClipboardPlus } from 'lucide-svelte';
	import Icon from '@/components/ui-custom/icon.svelte';
	import type { Snippet } from 'svelte';
	import { m } from '@/i18n';

	//

	type Props = ButtonProps & {
		delay?: number;
		textToCopy: string;
		children?: Snippet;
	};

	let { textToCopy, delay = 2000, children, ...rest }: Props = $props();

	let isCopied = $state(false);

	function copyText() {
		navigator.clipboard.writeText(textToCopy);
		isCopied = true;
		setTimeout(() => {
			isCopied = false;
		}, delay);
	}
</script>

<Button variant="outline" {...rest} onclick={copyText}>
	{#if !isCopied}
		{@render children?.()}
		<Icon src={ClipboardPlus} ml={Boolean(children)} />
	{:else}
		<span class="whitespace-nowrap">✅ {m.Copied()}</span>
	{/if}
</Button>
