<script lang="ts">
	export let text: string | undefined = undefined;
</script>

{#if text}
	<p class="text-xs text-gray-400">{text}</p>
{/if}
