<script lang="ts">
	import CopyButton from "./copyButton.svelte";

    export let rows:{title:string, value:string}[]
</script>


<div class="flex flex-col">
    {#each rows as row}
        <div class="text-md text-black font-semibold max-w-sm">{row.title}</div>
        <div class="flex flex-row items-center justify-between mb-4 text-md font-bold border-t">
                    <div class="text-ellipsis overflow-hidden w-full">{row.value}</div>
            <CopyButton textToCopy={row.value} />
        </div>
    {/each}
</div>