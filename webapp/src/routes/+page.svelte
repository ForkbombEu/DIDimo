<script lang="ts">
	import { currentUser } from '$lib/pocketbase';
	import {
		Footer,
		FooterBrand,
		FooterCopyright,
		FooterIcon,
		FooterLink,
		FooterLinkGroup,
		Input,
		Label,
		NavBrand,
		NavHamburger,
		NavLi,
		NavUl,
		Navbar,
		P,
		Textarea
	} from 'flowbite-svelte';
	import { Button } from 'flowbite-svelte';
	import {
		Contact,
		ContentWithImage,
		Facebook,
		FeatureDefault,
		FeatureItem,
		Github,
		HeroHeader,
		PricingBodyHead,
		PricingCard,
		PricingHead,
		PricingItem,
		PricingItemWrapper,
		Section,
		Twitter
	} from 'flowbite-svelte-blocks';
	import {
		PresentationChartLine,
		AcademicCap,
		Cog,
		Briefcase,
		CurrencyDollar,
		Bolt
	} from 'svelte-heros-v2';
	import { appTitle } from '$lib/strings';
	import { featureFlags } from '$lib/features';

</script>

<div class="flex flex-col gap-20">
	<header>
		<Navbar let:hidden let:toggle fluid={false}>
			<NavBrand href="/">
				<img src="/logo.svg" class="mr-3 h-6 sm:h-9" alt={`${appTitle} Logo`} />
			</NavBrand>
			<div class="flex items-center lg:order-2">
				{#if $featureFlags.AUTH}
					<Button id="btn-login" href="/login" color="dark">Log in</Button>
					<Button href="/register" class="ml-2" color="red">Sign up</Button>
				{/if}

				<NavHamburger
					on:click={toggle}
					btnClass="inline-flex items-center p-2 ml-1 text-sm text-gray-500 rounded-lg lg:hidden hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-gray-200 dark:text-gray-400 dark:hover:bg-gray-700 dark:focus:ring-gray-600"
				/>
			</div>
			<NavUl
				{hidden}
				divClass="justify-between items-center w-full lg:flex lg:w-auto lg:order-1"
				ulClass="flex flex-col mt-4 font-medium lg:flex-row lg:space-x-8 lg:mt-0"
			>
				<NavLi href="/" active={true}>Home</NavLi>
				<NavLi href="/">Company</NavLi>
				<NavLi href="/">Marketplace</NavLi>
				<NavLi href="/">Features</NavLi>
				<NavLi href="/">Team</NavLi>
				<NavLi href="/">Contact</NavLi>
			</NavUl>
		</Navbar>
	</header>

	<Section name="contentwithimg">
		<ContentWithImage>
			<svelte:fragment slot="h2">We did reinvent the wheel</svelte:fragment>
			<p class="mb-4">
				We are strategists, designers and developers. Innovators and problem solvers. Small enough
				to be simple and quick, but big enough to deliver the scope you want at the pace you need.
				Small enough to be simple and quick, but big enough to deliver the scope you want at the
				pace you need.
			</p>
			<p>
				We are strategists, designers and developers. Innovators and problem solvers. Small enough
				to be simple and quick.
			</p>

			<svelte:fragment slot="image">
				<img
					class="w-full rounded-lg"
					src="https://flowbite.s3.amazonaws.com/blocks/marketing-ui/content/office-long-2.png"
					alt="office content 1"
				/>
				<img
					class="mt-4 w-full lg:mt-10 rounded-lg"
					src="https://flowbite.s3.amazonaws.com/blocks/marketing-ui/content/office-long-1.png"
					alt="office content 2"
				/>
			</svelte:fragment>
		</ContentWithImage>
	</Section>

	<Section name="pricing">
		<PricingHead>
			<svelte:fragment slot="h2">Designed for business teams like yours</svelte:fragment>
			<svelte:fragment slot="paragraph"
				>Here at {appTitle} we focus on markets where technology, innovation, and capital can unlock
				long-term value and drive economic growth.</svelte:fragment
			>
		</PricingHead>
		<div class="space-y-8 lg:grid lg:grid-cols-3 sm:gap-6 xl:gap-10 lg:space-y-0">
			<PricingCard>
				<PricingBodyHead>
					<svelte:fragment slot="h3">Starter</svelte:fragment>
					<svelte:fragment slot="paragraph">
						Best option for personal use & for your next project.
					</svelte:fragment>
					<svelte:fragment slot="price">
						<span class="mr-2 text-5xl font-extrabold">$29</span>
						<span class="text-gray-500 dark:text-gray-400">/month</span>
					</svelte:fragment>
				</PricingBodyHead>
				<PricingItemWrapper>
					<PricingItem>
						<span>Individual configuration</span>
					</PricingItem>
					<PricingItem>
						<span>No setup, or hidden fees</span>
					</PricingItem>
					<PricingItem>
						<span>Team size: <span class="font-semibold">1 developer</span></span>
					</PricingItem>
					<PricingItem>
						<span>Premium support: <span class="font-semibold">6 months</span></span>
					</PricingItem>
					<PricingItem>
						<span>Free updates: <span class="font-semibold">6 months</span></span>
					</PricingItem>

					<svelte:fragment slot="btn">
						<Button color="red">Get started</Button>
					</svelte:fragment>
				</PricingItemWrapper>
			</PricingCard>

			<PricingCard>
				<PricingBodyHead>
					<svelte:fragment slot="h3">Company</svelte:fragment>
					<svelte:fragment slot="paragraph">
						Relevant for multiple users, extended & premium support.
					</svelte:fragment>
					<svelte:fragment slot="price">
						<span class="mr-2 text-5xl font-extrabold">$99</span>
						<span class="text-gray-500 dark:text-gray-400">/month</span>
					</svelte:fragment>
				</PricingBodyHead>
				<PricingItemWrapper>
					<PricingItem>
						<span>Individual configuration</span>
					</PricingItem>
					<PricingItem>
						<span>No setup, or hidden fees</span>
					</PricingItem>
					<PricingItem>
						<span>Team size: <span class="font-semibold">10 developer</span></span>
					</PricingItem>
					<PricingItem>
						<span>Premium support: <span class="font-semibold">24 months</span></span>
					</PricingItem>
					<PricingItem>
						<span>Free updates: <span class="font-semibold">24 months</span></span>
					</PricingItem>

					<svelte:fragment slot="btn">
						<Button color="red">Get started</Button>
					</svelte:fragment>
				</PricingItemWrapper>
			</PricingCard>

			<PricingCard>
				<PricingBodyHead>
					<svelte:fragment slot="h3">Enterprise</svelte:fragment>
					<svelte:fragment slot="paragraph">
						Best for large scale uses and extended redistribution rights.
					</svelte:fragment>
					<svelte:fragment slot="price">
						<span class="mr-2 text-5xl font-extrabold">$499</span>
						<span class="text-gray-500 dark:text-gray-400">/month</span>
					</svelte:fragment>
				</PricingBodyHead>
				<PricingItemWrapper>
					<PricingItem>
						<span>Individual configuration</span>
					</PricingItem>
					<PricingItem>
						<span>No setup, or hidden fees</span>
					</PricingItem>
					<PricingItem>
						<span>Team size: <span class="font-semibold">100+ developer</span></span>
					</PricingItem>
					<PricingItem>
						<span>Premium support: <span class="font-semibold">36 months</span></span>
					</PricingItem>
					<PricingItem>
						<span>Free updates: <span class="font-semibold">36 months</span></span>
					</PricingItem>

					<svelte:fragment slot="btn">
						<Button color="red">Get started</Button>
					</svelte:fragment>
				</PricingItemWrapper>
			</PricingCard>
		</div>
	</Section>

	<Section name="feature">
		<HeroHeader
			class="max-w-screen-md mb-8 lg:mb-16"
			h2Class="mb-4 text-4xl tracking-tight font-extrabold text-gray-900 dark:text-white"
			pClass="text-gray-500 sm:text-xl dark:text-gray-400"
		>
			<svelte:fragment slot="h2">Designed for business teams like yours</svelte:fragment>
			<svelte:fragment slot="paragraph"
				>Here at {appTitle} we focus on markets where technology, innovation, and capital can unlock
				long-term value and drive economic growth.</svelte:fragment
			>
		</HeroHeader>
		<FeatureDefault>
			<FeatureItem>
				<svelte:fragment slot="icon"
					><PresentationChartLine
						class="w-5 h-5 text-primary-600 lg:w-6 lg:h-6 dark:text-primary-300"
					/></svelte:fragment
				>
				<svelte:fragment slot="h3">Marketing</svelte:fragment>
				<svelte:fragment slot="paragraph"
					>Plan it, create it, launch it. Collaborate seamlessly with all the organization and hit
					your marketing goals every month with our marketing plan.</svelte:fragment
				>
			</FeatureItem>
			<FeatureItem>
				<svelte:fragment slot="icon"
					><AcademicCap
						class="w-5 h-5 text-primary-600 lg:w-6 lg:h-6 dark:text-primary-300"
					/></svelte:fragment
				>
				<svelte:fragment slot="h3">Legal</svelte:fragment>
				<svelte:fragment slot="paragraph"
					>Protect your organization, devices and stay compliant with our structured workflows and
					custom permissions made for you.</svelte:fragment
				>
			</FeatureItem>
			<FeatureItem>
				<svelte:fragment slot="icon"
					><Briefcase
						class="w-5 h-5 text-primary-600 lg:w-6 lg:h-6 dark:text-primary-300"
					/></svelte:fragment
				>
				<svelte:fragment slot="h3">Business Automation</svelte:fragment>
				<svelte:fragment slot="paragraph"
					>Auto-assign tasks, send Slack messages, and much more. Now power up with hundreds of new
					templates to help you get started.</svelte:fragment
				>
			</FeatureItem>
			<FeatureItem>
				<svelte:fragment slot="icon"
					><CurrencyDollar
						class="w-5 h-5 text-primary-600 lg:w-6 lg:h-6 dark:text-primary-300"
					/></svelte:fragment
				>
				<svelte:fragment slot="h3">Finance</svelte:fragment>
				<svelte:fragment slot="paragraph"
					>Audit-proof software built for critical financial operations like month-end close and
					quarterly budgeting.</svelte:fragment
				>
			</FeatureItem>
			<FeatureItem>
				<svelte:fragment slot="icon"
					><Bolt
						class="w-5 h-5 text-primary-600 lg:w-6 lg:h-6 dark:text-primary-300"
					/></svelte:fragment
				>
				<svelte:fragment slot="h3">Enterprise Design</svelte:fragment>
				<svelte:fragment slot="paragraph"
					>Craft beautiful, delightful experiences for both marketing and product with real
					cross-company collaboration.</svelte:fragment
				>
			</FeatureItem>
			<FeatureItem>
				<svelte:fragment slot="icon"
					><Cog
						class="w-5 h-5 text-primary-600 lg:w-6 lg:h-6 dark:text-primary-300"
					/></svelte:fragment
				>
				<svelte:fragment slot="h3">Operations</svelte:fragment>
				<svelte:fragment slot="paragraph"
					>Keep your company’s lights on with customizable, iterative, and structured workflows
					built for all efficient teams and individual.</svelte:fragment
				>
			</FeatureItem>
		</FeatureDefault>
	</Section>

	<Section name="contact">
		<Contact>
			<svelte:fragment slot="h2">Contact Us</svelte:fragment>
			<svelte:fragment slot="paragraph">
				Got a technical issue? Want to send feedback about a beta feature? Need details about our
				Business plan? Let us know.
			</svelte:fragment>
			<form class="space-y-8">
				<div>
					<Label for="email" class="block mb-2">Your email</Label>
					<Input id="email" name="email" placeholder={`name@${appTitle}.com`} required />
				</div>
				<div>
					<Label for="subject" class="block mb-2">Subject</Label>
					<Input
						id="subject"
						name="subject"
						placeholder="Let us know how we can help you"
						required
					/>
				</div>
				<div>
					<Textarea
						id="subject"
						name="subject"
						placeholder="Leave a comment..."
						label="Your message"
					/>
				</div>
				<Button>Send message</Button>
			</form>
		</Contact>
	</Section>

	<!-- svelte-ignore missing-declaration -->
	<Section>
		<Footer footerType="socialmedia">
			<div class="grid gap-12 lg:grid-cols-6 lg:gap-24">
				<div class="col-span-2">
					<FooterBrand
						href="/"
						src="/logo.svg"
						alt={`${appTitle} Logo`}
						name=""
						aClass="flex mb-6"
					/>
					<p class="text-gray-600 dark:text-gray-400">
						{appTitle} is an ecosystem built on top of Tailwind CSS including a component library, block
						sections, a Figma design system and other resources.
					</p>
				</div>
				<div>
					<h2 class="mb-6 text-sm font-semibold text-gray-400 uppercase dark:text-white">
						Resources
					</h2>
					<FooterLinkGroup>
						<FooterLink liClass="mb-4" href="/">{appTitle}</FooterLink>
						<FooterLink liClass="mb-4" href="https://tailwindcss.com/">Tailwind CSS</FooterLink>
						<FooterLink liClass="mb-4" href="https://flowbite-svelte.com/"
							>{appTitle}-Svelte</FooterLink
						>
					</FooterLinkGroup>
				</div>
				<div>
					<h2 class="mb-6 text-sm font-semibold text-gray-400 uppercase dark:text-white">
						HELP & SUPPORT
					</h2>
					<FooterLinkGroup>
						<FooterLink liClass="mb-4" href="https://discord.gg/4eeurUVvTy">Discord</FooterLink>
						<FooterLink
							liClass="mb-4"
							href="https://github.com/shinokada/flowbite-svelte-blocks/issues">GitHub</FooterLink
						>
					</FooterLinkGroup>
				</div>
				<div>
					<h2 class="mb-6 text-sm font-semibold text-gray-400 uppercase dark:text-white">
						Follow us
					</h2>
					<FooterLinkGroup>
						<FooterLink liClass="mb-4" href="https://github.com/shinokada/flowbite-svelte-blocks"
							>Gihub</FooterLink
						>
						<FooterLink liClass="mb-4" href="https://discord.gg/4eeurUVvTy">Discord</FooterLink>
					</FooterLinkGroup>
				</div>
				<div>
					<h2 class="mb-6 text-sm font-semibold text-gray-400 uppercase dark:text-white">Legal</h2>
					<FooterLinkGroup>
						<FooterLink
							liClass="mb-4"
							href="https://github.com/shinokada/flowbite-svelte-blocks/LICENSE">LICENSE</FooterLink
						>
					</FooterLinkGroup>
				</div>
			</div>
			<hr class="my-6 border-gray-200 sm:mx-auto dark:border-gray-700 lg:my-8" />
			<div class="sm:flex sm:items-center sm:justify-between">
				<FooterCopyright href="/" by={appTitle} />
				<div class="flex mt-4 space-x-6 sm:justify-center sm:mt-0">
					<FooterIcon href="/" class="text-gray-400 hover:text-gray-900" icon={Facebook} />
					<FooterIcon href="/" class="text-gray-400 hover:text-gray-900" icon={Twitter} />
					<FooterIcon href="/" class="text-gray-400 hover:text-gray-900" icon={Github} />
				</div>
			</div>
		</Footer>
	</Section>
</div>
