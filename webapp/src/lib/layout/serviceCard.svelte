<script lang="ts">
	import T from '@/components/ui-custom/t.svelte';
	import { Card } from '@/components/ui/card';
	import { m } from '@/i18n';

	type Service = {
		type: 'Wallet' | 'Issuer' | 'Verifier';
		name: string;
		organization: string;
		description: string;
		status: number | 'checking';
		lastCheck: string;
	};

	type Props = {
		service?: Service;
		class?: string;
	};

	const { service = defaultService(), class: className = '' }: Props = $props();

	function defaultService(): Service {
		return {
			type: 'Wallet',
			name: 'Didimo Wallet',
			organization: 'Didimo',
			description:
				'Didimo Wallet is a digital wallet that allows you to store your identity credentials securely.',
			status: Number.parseFloat(Math.random().toFixed(2)),
			lastCheck: '2021-10-01'
		};
	}

	const serviceStatus = $derived(
		service.status === 'checking'
			? m.checking_results_()
			: m.percentage_compliant({ percentage: `${service.status * 100}%` })
	);

	const serviceColor = $derived.by(() => {
		if (service.status === 'checking') {
			return 'text-muted-foreground';
		} else if (service.status >= 0.7) {
			return 'text-green-700';
		} else if (service.status >= 0.3) {
			return 'text-yellow-700';
		} else {
			return 'text-red-700';
		}
	});
</script>

<Card class="rounded-xl p-6 {className}">
	<div class="space-y-4">
		<div class="space-y-1">
			<T tag="small" class="text-muted-foreground block">{service.type}</T>
			<T tag="h4" class="block">{service.name}</T>
			<T tag="small" class="text-muted-foreground block">{service.organization}</T>
		</div>
		<T tag="small" class="block font-normal leading-snug">{service.description}</T>
	</div>

	<div class="mt-10 space-y-1">
		<T tag="small" class="block {serviceColor}">{serviceStatus}</T>
		<T tag="small" class="block font-normal">{m.last_check()} {service.lastCheck}</T>
	</div>
</Card>
