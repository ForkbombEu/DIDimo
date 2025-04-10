<!--
SPDX-FileCopyrightText: 2025 Forkbomb BV

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import { appName } from '@/brand';
	import T from '@/components/ui-custom/t.svelte';
	import Button from '@/components/ui-custom/button.svelte';
	import { featureFlags } from '@/features';
	import { m } from '@/i18n';
	import { appVersion } from '@/utils/appVersion';

	const footer_links = {
		Platform: [
			{ label: m.Overview(), url: '/platform' },
			{ label: m.api_documentation(), url: '/docs/api' },
			{ label: m.platform_architecture(), url: '/platform/architecture' },
			{ label: m.deployment_options(), url: '/platform/deployment' },
			{ label: m.compliance_standards(), url: '/platform/compliance' }
		],
		Solutions: [
			{ label: m.didital_identity_for_governments(), url: '/solutions/government' },
			{ label: m.financial_services(), url: '/solutions/finance' },
			{ label: m.workforce_verification(), url: '/solutions/workforce' },
			{ label: m.education_credentials(), url: '/solutions/education' },
			{ label: m.transportation_smart_mobility(), url: '/solutions/transportation' }
		],
		Capabilities: [
			{ label: m.credential_issuance(), url: '/capabilities/issuance' },
			{ label: m.credential_verification(), url: '/capabilities/verification' },
			{ label: m.identity_wallet_integration(), url: '/capabilities/wallet-integration' },
			{ label: m.conformance_compliance_tools(), url: '/capabilities/conformance' },
			{ label: m.ecosystem_governance(), url: '/capabilities/ecosystem-governance' }
		],
		Developers: [
			{ label: m.get_started_with_api(), url: '/developers/get-started' },
			{ label: m.sdks_libraries(), url: '/developers/sdks' },
			{ label: m.sandbox_environment(), url: '/developers/sandbox' },
			{ label: m.test_suites_runs(), url: '/developers/test-suites' },
			{ label: m.developer_tools(), url: '/developers/tools' },
			{ label: m.Changelog(), url: '/developers/changelog' }
		],
		Resources: [
			{ label: m.Documentation(), url: '/resources/docs' },
			{ label: m.articles_tutorials(), url: '/resources/articles' },
			{ label: m.videos(), url: '/resources/videos' },
			{ label: m.case_studies(), url: '/resources/case-studies' },
			{ label: m.press_releases(), url: '/resources/press' },
			{ label: 'Blog', url: '/blog' }
		],
		Community: [
			{ label: m.join_our_slack(), url: '/community/slack' },
			{ label: 'GitHub', url: 'https://github.com/didimoid' },
			{ label: m.events_webinars(), url: '/community/events' },
			{ label: m.conformance_community_forum(), url: '/community/forum' },
			{ label: m.partner_network(), url: '/partners' }
		],
		Company: [
			{ label: m.about_credimi(), url: '/about' },
			{ label: m.our_vision_and_approach(), url: '/about/approach' },
			{ label: m.careers(), url: '/careers' },
			{ label: m.security(), url: '/security' },
			{ label: m.contact_us(), url: '/contact' },
			{ label: m.legal_privacy(), url: '/legal' }
		]
	};
</script>

<div class="bg-primary text-muted">
	<div class="mx-auto max-w-screen-xl">
		<div class="px-8 py-12">
			<div class="flex flex-col justify-between gap-6 sm:flex-row">
				<div>
					<T tag="h4">
						{appName}
						<span class="pl-1 text-sm text-white/40">{appVersion}</span>
					</T>
					<p class="border-b border-b-muted">{m.Test_and_find_decentralized_IDs()}</p>
				</div>
				<div class="flex gap-2 blur-sm">
					<Button
						variant="default"
						class="grow basis-1 bg-white text-primary hover:bg-white/90 sm:grow-0"
						href={$featureFlags.DEMO ? undefined : '/tests/new'}
						disabled={$featureFlags.DEMO}
					>
						{m.Start_a_new_test()}
					</Button>
					<Button
						variant="ghost"
						class="grow basis-1 border sm:grow-0"
						href={$featureFlags.DEMO ? undefined : '/login'}
						disabled={$featureFlags.DEMO}
					>
						{m.Login()}
					</Button>
				</div>
			</div>
			<div class="grid grid-cols-1 gap-12 pt-12 sm:grid-cols-2 md:grid-cols-4">
				{#each Object.entries(footer_links) as [category, links]}
					<div class="flex flex-row gap-2">
						<div class="flex flex-col gap-2">
							<T tag="h4">{category}</T>
							{#each links as link}
								<a href={link.url} class="text-muted-foreground hover:text-white"
									>{link.label}</a
								>
							{/each}
						</div>
					</div>
				{/each}
			</div>
		</div>
	</div>
	<div class="flex items-center justify-center bg-black/30 p-2 px-8">
		<small>
			<a href="https://forkbomb.solutions/" class="underline hover:no-underline">
				{m.Proudly_developed_by_ForkBomb_BV()}
			</a>
		</small>
	</div>
</div>
