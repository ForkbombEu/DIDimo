<script lang="ts" generics="Response extends object">
	import { m } from '@/i18n';
	import type { CollectionName } from '@/pocketbase/collections-models';
	import type { CollectionResponses } from '@/pocketbase/types';
	import {
		RecordDelete,
		RecordEdit,
		RecordSelect,
		RecordShare,
		type RecordAction
	} from '../record-actions';
	import type { KeyOf } from '@/utils/types';
	import * as Table from '@/components/ui/table';
	import FieldTh from './fieldTh.svelte';
	import IconButton from '@/components/ui-custom/iconButton.svelte';
	import type { Snippet } from 'svelte';

	//

	interface Props {
		records: Response[];
		fields?: KeyOf<Response>[];
		hide?: Array<RecordAction>;
		header?: Snippet<[{ Th: typeof Table.Head }]>;
		row?: Snippet<[{ Td: typeof Table.Cell; record: Response }]>;
		actions?: Snippet<[{ record: Response }]>;
	}

	const {
		records,
		fields = ['id'] as KeyOf<Response>[],
		hide = [],
		header,
		row,
		actions
	}: Props = $props();
</script>

<Table.Root>
	<Table.Header>
		<Table.Row class="hover:bg-inherit">
			{#if !hide.includes('select')}
				<Table.Head />
			{/if}

			{#each fields as field}
				<FieldTh {field} />
			{/each}

			{@render header?.({ Th: Table.Head })}

			<Table.Head>
				{m.Actions()}
			</Table.Head>
		</Table.Row>
	</Table.Header>

	<Table.Body>
		{#each records as untypedRecord}
			{@const record = untypedRecord as CollectionResponses[CollectionName]}
			<Table.Row class="hover:bg-inherit">
				{#if !hide.includes('select')}
					<Table.Cell class="py-2">
						<RecordSelect {record} />
					</Table.Cell>
				{/if}

				{#each fields as field}
					<Table.Cell>
						{untypedRecord[field]}
					</Table.Cell>
				{/each}

				{@render row?.({ record: untypedRecord, Td: Table.Cell })}

				<Table.Cell class="py-2">
					{@render actions?.({ record: untypedRecord })}

					{#if !hide.includes('edit')}
						<RecordEdit {record}>
							{#snippet button({ triggerAttributes, icon })}
								<IconButton {icon} variant="ghost" {...triggerAttributes} />
							{/snippet}
						</RecordEdit>
					{/if}
					{#if !hide.includes('share')}
						<RecordShare {record}>
							{#snippet button({ triggerAttributes, icon })}
								<IconButton {icon} variant="ghost" {...triggerAttributes} />
							{/snippet}
						</RecordShare>
					{/if}
					{#if !hide.includes('delete')}
						<RecordDelete {record}>
							{#snippet button({ triggerAttributes, icon })}
								<IconButton {icon} variant="ghost" {...triggerAttributes} />
							{/snippet}
						</RecordDelete>
					{/if}
				</Table.Cell>
			</Table.Row>
		{/each}
	</Table.Body>
</Table.Root>
