<script lang="ts">
	import {
		DarkMode,
		Dropdown,
		DropdownItem,
		NavBrand,
		NavLi,
		NavUl,
		Navbar,
		Search
	} from 'flowbite-svelte';
	import {
		ArchiveSolid,
		ArrowRightToBracketOutline,
		CogOutline,
		DollarOutline,
		InboxOutline,
		ProfileCardOutline,
		SalePercentOutline,
		ShoppingBagSolid,
		UsersGroupSolid,
		ChevronDownOutline
	} from 'flowbite-svelte-icons';

	import '../app.css';

	interface Props {
		fluid?: boolean;
		drawerHidden?: boolean;
		list?: boolean;
	}

	let { fluid = true, drawerHidden = $bindable(false), list = false }: Props = $props();

	interface Link {
		href: string;
		label: string;
	}
	const links: Link[] = [
		{ href: '/', label: 'Home' },
		{ href: '/products', label: 'Products' }
	];

	const menu = [
		{ name: 'Sales', href: '/', icon: ShoppingBagSolid },
		{ name: 'Users', href: '/', icon: UsersGroupSolid },
		{ name: 'Inbox', href: '/', icon: InboxOutline },
		{ name: 'Profile', href: '/', icon: ProfileCardOutline },
		{ name: 'Settings', href: '/settings', icon: CogOutline },
		{ name: 'Prouducts', href: '/', icon: ArchiveSolid },
		{ name: 'Pricing', href: '/pages/pricing', icon: DollarOutline },
		{ name: 'Billing', href: '/', icon: SalePercentOutline },
		{ name: 'Logout', href: '/', icon: ArrowRightToBracketOutline }
	];
</script>

<Navbar class="mx-10 sm:mx-0">
	<NavBrand href="/" class="mx-10">
		<img
			src="/images/flowbite-svelte-icon-logo.svg"
			class="me-2.5 h-6 sm:h-8"
			alt="Flowbite Logo"
		/>
		<span
			class="ml-px self-center text-xl font-semibold whitespace-nowrap sm:text-2xl dark:text-white"
		>
			Flowbite
		</span>
	</NavBrand>
	<div class="hidden lg:block lg:ps-3">
		{#if list}
			<NavUl class="ml-2" activeUrl="/" activeClass="text-primary-600 dark:text-primary-500">
				<NavLi href="/">Home</NavLi>
				<NavLi href="#top">Messages</NavLi>
				<NavLi href="#top">Profile</NavLi>
				<NavLi href="#top">Settings</NavLi>
				<NavLi class="cursor-pointer">
					Dropdown
					<ChevronDownOutline class="ms-0 inline" />
				</NavLi>
				<Dropdown class="z-20 w-44">
					<DropdownItem href="#top">Item 1</DropdownItem>
					<DropdownItem href="#top">Item 2</DropdownItem>
					<DropdownItem href="#top">Item 3</DropdownItem>
				</Dropdown>
			</NavUl>
		{:else}
			<form>
				<Search size="md" class="mt-1 w-96 border focus:outline-none" />
			</form>
		{/if}
	</div>
	<div class="ms-auto flex items-center text-gray-500 sm:order-2 dark:text-gray-300">
		<DarkMode />
	</div>
</Navbar>
