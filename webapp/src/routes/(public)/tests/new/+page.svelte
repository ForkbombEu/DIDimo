<script lang="ts">
	import { CollectionForm } from '@/collections-components';
	import T from '@/components/ui-custom/t.svelte';
	import { m } from '@/i18n';
</script>

<div class="mx-auto max-w-screen-xl px-8 py-12">
	<T tag="h1" class="mb-12">{m.Upload_a_new_compliance_test()}</T>
	<CollectionForm collection="tests">
		{#snippet submitButtonContent()}
			{m.Upload_test()}
		{/snippet}
	</CollectionForm>
</div>
