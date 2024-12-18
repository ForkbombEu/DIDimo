<script lang="ts">
	import * as Form from '@/components/ui/form';
	import { getFormContext } from '@/forms';
	import { FORM_ERROR_PATH } from '../form';
	import type { ButtonProps } from '@/components/ui/button';

	//

	const { children, ...props }: ButtonProps = $props();

	const { form } = $derived(getFormContext());
	const { allErrors } = $derived(form);
	const formHasErrors = $derived($allErrors.filter((e) => e.path != FORM_ERROR_PATH).length > 0);
</script>

<Form.Button {...props} disabled={formHasErrors}>
	{@render children?.()}
</Form.Button>
