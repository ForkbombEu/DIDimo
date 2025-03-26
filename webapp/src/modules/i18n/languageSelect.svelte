<script lang="ts">
	import Icon from '@/components/ui-custom/icon.svelte';
	import * as Popover from '@/components/ui/popover';
	import { Button } from '@/components/ui/button';
	import BaseLanguageSelect from './baseLanguageSelect.svelte';
	import type { LanguageSelectTriggerSnippetProps } from './baseLanguageSelect.svelte';
	import type { Snippet } from 'svelte';
	import type { GenericRecord } from '@/utils/types';
	import { setLocale } from './paraglide/runtime';

	type Props = {
		contentClass?: string;
		trigger?: Snippet<
			[LanguageSelectTriggerSnippetProps & { triggerAttributes: GenericRecord }]
		>;
		flagsOnly?: boolean;
	};

	const { contentClass = '', trigger: triggerSnippet, flagsOnly }: Props = $props();
</script>

<Popover.Root>
	<BaseLanguageSelect>
		{#snippet trigger(data)}
			<Popover.Trigger>
				{#snippet child({ props: triggerAttributes })}
					{#if triggerSnippet}
						{@render triggerSnippet({ ...data, triggerAttributes })}
					{:else}
						{@const { icon: LanguageIcon, text } = data}
						<Button variant="outline" {...triggerAttributes}>
							<Icon src={LanguageIcon} />
							{text}
						</Button>
					{/if}
				{/snippet}
			</Popover.Trigger>
		{/snippet}

		{#snippet languages({ languages })}
			<Popover.Content class="space-y-0.5 p-1 {contentClass} w-[--bits-popover-anchor-width]">
				{#each languages as { name, flag, isCurrent, tag }}
					<Button
						onclick={() => setLocale(tag)}
						variant={isCurrent ? 'secondary' : 'ghost'}
						class="flex items-center justify-start gap-2"
						size="sm"
					>
						<span class="text-2xl">
							{flag}
						</span>
						{#if !flagsOnly}
							<span>
								{name}
							</span>
						{/if}
					</Button>
				{/each}
			</Popover.Content>
		{/snippet}
	</BaseLanguageSelect>
</Popover.Root>
