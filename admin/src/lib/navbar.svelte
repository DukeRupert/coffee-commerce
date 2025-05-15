<script lang="ts">
	import { Navbar, NavBrand, NavLi, NavUl, NavHamburger } from 'flowbite-svelte';

	interface Link {
		href: string;
		label: string;
	}
	const links: Link[] = [
		{ href: '/', label: 'Home' },
		{ href: '/products', label: 'Products' }
	];
</script>

<Navbar>
	<NavBrand href="/">
		<span class="self-center whitespace-nowrap text-xl font-semibold dark:text-white">Flowbite</span
		>
	</NavBrand>
	<NavHamburger />
	<NavUl>
		{#each links as link}
			<NavLi href={link.href}>{link.label}</NavLi>
		{/each}
	</NavUl>
</Navbar>
