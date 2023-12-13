<script lang="ts">
	import { Button } from 'flowbite-svelte';
	import { ClipboardDocument } from 'svelte-heros-v2';

	export let textToCopy: string;
	export let delay = 2000;

	let isCopied = false;

	function copyText() {
		navigator.clipboard.writeText(textToCopy);
		isCopied = true;
		setTimeout(() => {
			isCopied = false;
		}, delay);
	}
</script>

<Button on:click={copyText} color="alternative">
	{#if !isCopied}
		<ClipboardDocument size="20" />
		<span class="ml-2">
			<slot />
		</span>
	{:else}
		<span class="whitespace-nowrap">✅ Copied!</span>
	{/if}
</Button>
