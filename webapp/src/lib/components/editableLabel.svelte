<script lang="ts">
	import { Button } from "flowbite-svelte";
	import { ArrowLeft, Pencil } from "svelte-heros-v2";
    export let valueClass: string = "";
    export let value: string;
    export let handleChange: (value: string) => void;

    let isInputNameOpen = false;
    const handleSave = () => {
        handleChange(value);
        isInputNameOpen = false;
    }
</script>

<div class="flex flex-row items-center justify-between w-full">
	{#if !isInputNameOpen}
		<h5 class={`${valueClass}`}>
			{value}
		</h5>
		<Button on:click={() => (isInputNameOpen = true)} outline size="sm" class="w-14">
			<Pencil />
		</Button>
	{/if}

	{#if isInputNameOpen}
		<input
			type="text"
			bind:value={value}
			class="border border-gray-300 dark:border-gray-700 rounded-md w-full px-3 py-2 mr-2 text-sm text-gray-900 dark:text-white"
		/>
		<div class="flex flex-row">
			<Button
				on:click={() => (isInputNameOpen = false)}
				outline
				size="sm"
				class="px-0 text-bold w-8 h-8"
			>
				<div class="flex flex-col">
					<ArrowLeft class="w-8" />
					<!-- <span class="text-sm font-bold">Back</span> -->
				</div>
			</Button>
			<Button on:click={handleSave} outline size="sm" class="w-8 h-8 text-red-500 px-0">
				<div class="flex flex-col">
					<Pencil class="w-8" />
					<!-- <span class="text-sm  font-bold">Save</span> -->
				</div>
			</Button>
		</div>
	{/if}
</div>
