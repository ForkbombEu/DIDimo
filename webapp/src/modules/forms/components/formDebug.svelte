<script lang="ts" module>
	import type { GenericRecord } from '@/utils/types';
</script>

<script lang="ts" generics="T extends GenericRecord">
	import SuperDebug, { type SuperForm } from 'sveltekit-superforms';

	interface Props {
		form: SuperForm<T>;
	}

	let { form }: Props = $props();
	const { form: formData } = form;
</script>

<SuperDebug data={$formData}></SuperDebug>
