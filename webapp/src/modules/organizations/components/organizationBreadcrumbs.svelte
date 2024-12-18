<script lang="ts">
	import { m } from '@/i18n';
	import { pb } from '@/pocketbase';
	import type { BreadcrumbsOptions } from '@/components/ui-custom/breadcrumbs.svelte';
	import Breadcrumbs from '@/components/ui-custom/breadcrumbs.svelte';

	//

	const breadcrumbsOptions: BreadcrumbsOptions = {
		renamers: {
			'[id]': getOrganizationNameById,
			organizations: () => m.organizations(),
			my: () => m.My()
		},
		exclude: ['[[lang]]']
	};

	async function getOrganizationNameById(id: string): Promise<string> {
		const organization = await pb.collection('organizations').getOne(id);
		return organization.name;
	}
</script>

<Breadcrumbs options={breadcrumbsOptions} />
