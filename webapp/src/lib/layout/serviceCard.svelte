<script lang="ts">
	import T from '@/components/ui-custom/t.svelte';
	import { m } from '@/i18n';
	import { type ServicesRecord, type ServicesResponse } from '@/pocketbase/types';

	//

	type Props = {
		service: ServicesResponse;
		class?: string;
	};

	const { service, class: className = '' }: Props = $props();

	//

	type ServiceType = 'Wallet' | 'Issuer' | 'Verifier' | 'Unknown';

	const serviceType = $derived(getServiceType(service));

	function getServiceType(service: ServicesRecord): ServiceType {
		if (Boolean(service.credential_issuers?.length)) return 'Issuer';
		else return 'Unknown';
	}

	//

	type ServiceStatus = 'checking' | number;

	const serviceStatus = $state<ServiceStatus>('checking');

	const serviceStatusText = $derived(
		serviceStatus === 'checking'
			? m.checking_results_()
			: m.percentage_compliant({ percentage: `${serviceStatus * 100}%` })
	);

	const serviceColor = $derived.by(() => {
		if (serviceStatus === 'checking') {
			return 'text-primary';
		} else if (serviceStatus >= 0.7) {
			return 'text-green-700';
		} else if (serviceStatus >= 0.3) {
			return 'text-yellow-700';
		} else {
			return 'text-red-700';
		}
	});
</script>

<a
	href="/services/{service.id}"
	class="bg-card text-card-foreground border-primary ring-primary rounded-lg border p-6 shadow-sm transition-all hover:-translate-y-2 hover:ring-2 {className}"
>
	<div class="space-y-4">
		<div class="space-y-1">
			<T tag="small" class="text-primary block">{serviceType}</T>
			<T tag="h4" class="block">{service.name}</T>
			<!-- <T tag="small" class="text-muted-foreground block">{service.organization}</T> -->
		</div>
		<T tag="small" class="block font-normal leading-snug">{service.description}</T>
	</div>

	<div class="mt-10 space-y-1">
		<T tag="small" class="block font-semibold {serviceColor}">{serviceStatusText}</T>
		<T tag="small" class="block font-light">
			{m.last_check()}
			{service.updated.split(' ').at(0)}
		</T>
	</div>
</a>
