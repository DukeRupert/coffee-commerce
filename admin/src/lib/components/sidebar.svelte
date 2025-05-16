<script lang="ts">
	import type { Component, Snippet } from 'svelte';
	import { Home, Coffee } from '@lucide/svelte';

	interface PageProps {
		children: Snippet;
	}
	let { children }: PageProps = $props();

	interface Link {
		href: string;
		label: string;
		icon: Component;
	}

	const links: Link[] = [
		{
			href: '/',
			label: 'home',
			icon: Home
		},
		{
			href: 'products/',
			label: 'products',
			icon: Coffee
		}
	];

	// Reactive state variables for UI controls
	let isMobileSidebarOpen = $state(false);
	let isUserMenuOpen = $state(false);

	// Function to toggle mobile sidebar
	function toggleMobileSidebar() {
		isMobileSidebarOpen = !isMobileSidebarOpen;
	}

	// Function to toggle user menu dropdown
	function toggleUserMenu() {
		isUserMenuOpen = !isUserMenuOpen;
	}

	// Close mobile sidebar
	function closeMobileSidebar() {
		isMobileSidebarOpen = false;
	}

	// Close menus when clicking outside
	function handleClickOutside(event: MouseEvent) {
		// User menu close logic
		const userMenuButton: HTMLElement | null = document.getElementById('user-menu-button');
		const userMenu: HTMLElement | null = document.getElementById('user-menu');
		const target = event.target as Node;

		if (
			isUserMenuOpen &&
			userMenu &&
			userMenuButton &&
			!userMenu.contains(target) &&
			!userMenuButton.contains(target)
		) {
			isUserMenuOpen = false;
		}
	}

	// Add global click handler when component is mounted
	$effect(() => {
		document.addEventListener('click', handleClickOutside);

		// Cleanup when component is destroyed
		return () => {
			document.removeEventListener('click', handleClickOutside);
		};
	});
</script>

{#snippet Link(link: Link)}
	<li>
		<a
			href={link.href}
			class="group flex gap-x-3 rounded-md p-2 text-sm/6 font-semibold capitalize text-gray-400 hover:bg-gray-800 hover:text-white"
		>
			<link.icon></link.icon>
			{link.label}
		</a>
	</li>
{/snippet}

{#snippet MobileLink(link: Link)}
	<li>
		<a
			href={link.href}
			class="group flex gap-x-3 rounded-md bg-gray-800 p-2 text-sm/6 font-semibold capitalize text-white"
		>
			<link.icon></link.icon>
			{link.label}
		</a>
	</li>
{/snippet}

<div>
	<!-- Mobile menu overlay, conditionally rendered with transitions -->
	<div
		class="relative z-50 lg:hidden"
		role="dialog"
		aria-modal="true"
		class:hidden={!isMobileSidebarOpen}
	>
		<!-- Backdrop with transition -->
		<div
			class="fixed inset-0 bg-gray-900/80 transition-opacity duration-300 ease-linear"
			class:opacity-100={isMobileSidebarOpen}
			class:opacity-0={!isMobileSidebarOpen}
			aria-hidden="true"
		></div>

		<div class="fixed inset-0 flex">
			<!-- Mobile sidebar with transition -->
			<div
				class="relative mr-16 flex w-full max-w-xs flex-1 transform transition-transform duration-300 ease-in-out"
				class:translate-x-0={isMobileSidebarOpen}
				class:translate-x-full={!isMobileSidebarOpen}
			>
				<!-- Close button with transition -->
				<div
					class="absolute left-full top-0 flex w-16 transform justify-center pt-5 transition-opacity duration-300 ease-in-out"
					class:opacity-100={isMobileSidebarOpen}
					class:opacity-0={!isMobileSidebarOpen}
				>
					<button type="button" class="-m-2.5 p-2.5" onclick={closeMobileSidebar}>
						<span class="sr-only">Close sidebar</span>
						<svg
							class="size-6 text-white"
							fill="none"
							viewBox="0 0 24 24"
							stroke-width="1.5"
							stroke="currentColor"
							aria-hidden="true"
							data-slot="icon"
						>
							<path stroke-linecap="round" stroke-linejoin="round" d="M6 18 18 6M6 6l12 12" />
						</svg>
					</button>
				</div>

				<!-- Sidebar component -->
				<div
					class="flex grow flex-col gap-y-5 overflow-y-auto bg-gray-900 px-6 pb-4 ring-1 ring-white/10"
				>
					<div class="flex h-16 shrink-0 items-center">
						<img
							class="h-8 w-auto"
							src="https://tailwindcss.com/plus-assets/img/logos/mark.svg?color=indigo&shade=500"
							alt="Your Company"
						/>
					</div>
					<nav class="flex flex-1 flex-col">
						<ul role="list" class="flex flex-1 flex-col gap-y-7">
							<li>
								<ul role="list" class="-mx-2 space-y-1">
									{#each links as link}
										{@render MobileLink(link)}
									{/each}
								</ul>
							</li>
							<li class="mt-auto">
								<a
									href="#"
									class="group -mx-2 flex gap-x-3 rounded-md p-2 text-sm/6 font-semibold text-gray-400 hover:bg-gray-800 hover:text-white"
								>
									<svg
										class="size-6 shrink-0"
										fill="none"
										viewBox="0 0 24 24"
										stroke-width="1.5"
										stroke="currentColor"
										aria-hidden="true"
										data-slot="icon"
									>
										<path
											stroke-linecap="round"
											stroke-linejoin="round"
											d="M9.594 3.94c.09-.542.56-.94 1.11-.94h2.593c.55 0 1.02.398 1.11.94l.213 1.281c.063.374.313.686.645.87.074.04.147.083.22.127.325.196.72.257 1.075.124l1.217-.456a1.125 1.125 0 0 1 1.37.49l1.296 2.247a1.125 1.125 0 0 1-.26 1.431l-1.003.827c-.293.241-.438.613-.43.992a7.723 7.723 0 0 1 0 .255c-.008.378.137.75.43.991l1.004.827c.424.35.534.955.26 1.43l-1.298 2.247a1.125 1.125 0 0 1-1.369.491l-1.217-.456c-.355-.133-.75-.072-1.076.124a6.47 6.47 0 0 1-.22.128c-.331.183-.581.495-.644.869l-.213 1.281c-.09.543-.56.94-1.11.94h-2.594c-.55 0-1.019-.398-1.11-.94l-.213-1.281c-.062-.374-.312-.686-.644-.87a6.52 6.52 0 0 1-.22-.127c-.325-.196-.72-.257-1.076-.124l-1.217.456a1.125 1.125 0 0 1-1.369-.49l-1.297-2.247a1.125 1.125 0 0 1 .26-1.431l1.004-.827c.292-.24.437-.613.43-.991a6.932 6.932 0 0 1 0-.255c.007-.38-.138-.751-.43-.992l-1.004-.827a1.125 1.125 0 0 1-.26-1.43l1.297-2.247a1.125 1.125 0 0 1 1.37-.491l1.216.456c.356.133.751.072 1.076-.124.072-.044.146-.086.22-.128.332-.183.582-.495.644-.869l.214-1.28Z"
										/>
										<path
											stroke-linecap="round"
											stroke-linejoin="round"
											d="M15 12a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z"
										/>
									</svg>
									Settings
								</a>
							</li>
						</ul>
					</nav>
				</div>
			</div>
		</div>
	</div>

	<!-- Static sidebar for desktop -->
	<div class="hidden lg:fixed lg:inset-y-0 lg:z-50 lg:flex lg:w-72 lg:flex-col">
		<!-- Sidebar component, swap this element with another sidebar if you like -->
		<div class="flex grow flex-col gap-y-5 overflow-y-auto bg-gray-900 px-6 pb-4">
			<div class="flex h-16 shrink-0 items-center">
				<img
					class="h-8 w-auto"
					src="https://tailwindcss.com/plus-assets/img/logos/mark.svg?color=indigo&shade=500"
					alt="Your Company"
				/>
			</div>
			<nav class="flex flex-1 flex-col">
				<ul role="list" class="flex flex-1 flex-col gap-y-7">
					<li>
						<ul role="list" class="-mx-2 space-y-1">
							{#each links as link}
								{@render Link(link)}
							{/each}
							<!-- Other menu items... -->
						</ul>
					</li>
					<li>
						<div class="text-xs/6 font-semibold text-gray-400">Your teams</div>
						<ul role="list" class="-mx-2 mt-2 space-y-1">
							<li>
								<a
									href="#"
									class="group flex gap-x-3 rounded-md p-2 text-sm/6 font-semibold text-gray-400 hover:bg-gray-800 hover:text-white"
								>
									<span
										class="flex size-6 shrink-0 items-center justify-center rounded-lg border border-gray-700 bg-gray-800 text-[0.625rem] font-medium text-gray-400 group-hover:text-white"
										>H</span
									>
									<span class="truncate">Heroicons</span>
								</a>
							</li>
							<!-- Other teams... -->
						</ul>
					</li>
					<li class="mt-auto">
						<a
							href="#"
							class="group -mx-2 flex gap-x-3 rounded-md p-2 text-sm/6 font-semibold text-gray-400 hover:bg-gray-800 hover:text-white"
						>
							<svg
								class="size-6 shrink-0"
								fill="none"
								viewBox="0 0 24 24"
								stroke-width="1.5"
								stroke="currentColor"
								aria-hidden="true"
								data-slot="icon"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									d="M9.594 3.94c.09-.542.56-.94 1.11-.94h2.593c.55 0 1.02.398 1.11.94l.213 1.281c.063.374.313.686.645.87.074.04.147.083.22.127.325.196.72.257 1.075.124l1.217-.456a1.125 1.125 0 0 1 1.37.49l1.296 2.247a1.125 1.125 0 0 1-.26 1.431l-1.003.827c-.293.241-.438.613-.43.992a7.723 7.723 0 0 1 0 .255c-.008.378.137.75.43.991l1.004.827c.424.35.534.955.26 1.43l-1.298 2.247a1.125 1.125 0 0 1-1.369.491l-1.217-.456c-.355-.133-.75-.072-1.076.124a6.47 6.47 0 0 1-.22.128c-.331.183-.581.495-.644.869l-.213 1.281c-.09.543-.56.94-1.11.94h-2.594c-.55 0-1.019-.398-1.11-.94l-.213-1.281c-.062-.374-.312-.686-.644-.87a6.52 6.52 0 0 1-.22-.127c-.325-.196-.72-.257-1.076-.124l-1.217.456a1.125 1.125 0 0 1-1.369-.49l-1.297-2.247a1.125 1.125 0 0 1 .26-1.431l1.004-.827c.292-.24.437-.613.43-.991a6.932 6.932 0 0 1 0-.255c.007-.38-.138-.751-.43-.992l-1.004-.827a1.125 1.125 0 0 1-.26-1.43l1.297-2.247a1.125 1.125 0 0 1 1.37-.491l1.216.456c.356.133.751.072 1.076-.124.072-.044.146-.086.22-.128.332-.183.582-.495.644-.869l.214-1.28Z"
								/>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									d="M15 12a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z"
								/>
							</svg>
							Settings
						</a>
					</li>
				</ul>
			</nav>
		</div>
	</div>

	<div class="lg:pl-72">
		<div
			class="shadow-xs sticky top-0 z-40 flex h-16 shrink-0 items-center gap-x-4 border-b border-gray-200 bg-white px-4 sm:gap-x-6 sm:px-6 lg:px-8"
		>
			<button
				type="button"
				class="-m-2.5 p-2.5 text-gray-700 lg:hidden"
				onclick={toggleMobileSidebar}
			>
				<span class="sr-only">Open sidebar</span>
				<svg
					class="size-6"
					fill="none"
					viewBox="0 0 24 24"
					stroke-width="1.5"
					stroke="currentColor"
					aria-hidden="true"
					data-slot="icon"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5"
					/>
				</svg>
			</button>

			<!-- Separator -->
			<div class="h-6 w-px bg-gray-900/10 lg:hidden" aria-hidden="true"></div>

			<div class="flex flex-1 gap-x-4 self-stretch lg:gap-x-6">
				<form class="grid flex-1 grid-cols-1" action="#" method="GET">
					<input
						type="search"
						name="search"
						aria-label="Search"
						class="outline-hidden col-start-1 row-start-1 block size-full bg-white pl-8 text-base text-gray-900 placeholder:text-gray-400 sm:text-sm/6"
						placeholder="Search"
					/>
					<svg
						class="pointer-events-none col-start-1 row-start-1 size-5 self-center text-gray-400"
						viewBox="0 0 20 20"
						fill="currentColor"
						aria-hidden="true"
						data-slot="icon"
					>
						<path
							fill-rule="evenodd"
							d="M9 3.5a5.5 5.5 0 1 0 0 11 5.5 5.5 0 0 0 0-11ZM2 9a7 7 0 1 1 12.452 4.391l3.328 3.329a.75.75 0 1 1-1.06 1.06l-3.329-3.328A7 7 0 0 1 2 9Z"
							clip-rule="evenodd"
						/>
					</svg>
				</form>
				<div class="flex items-center gap-x-4 lg:gap-x-6">
					<button type="button" class="-m-2.5 p-2.5 text-gray-400 hover:text-gray-500">
						<span class="sr-only">View notifications</span>
						<svg
							class="size-6"
							fill="none"
							viewBox="0 0 24 24"
							stroke-width="1.5"
							stroke="currentColor"
							aria-hidden="true"
							data-slot="icon"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								d="M14.857 17.082a23.848 23.848 0 0 0 5.454-1.31A8.967 8.967 0 0 1 18 9.75V9A6 6 0 0 0 6 9v.75a8.967 8.967 0 0 1-2.312 6.022c1.733.64 3.56 1.085 5.455 1.31m5.714 0a24.255 24.255 0 0 1-5.714 0m5.714 0a3 3 0 1 1-5.714 0"
							/>
						</svg>
					</button>

					<!-- Separator -->
					<div class="hidden lg:block lg:h-6 lg:w-px lg:bg-gray-900/10" aria-hidden="true"></div>

					<!-- Profile dropdown -->
					<div class="relative">
						<button
							type="button"
							class="-m-1.5 flex items-center p-1.5"
							id="user-menu-button"
							onclick={toggleUserMenu}
							aria-expanded={isUserMenuOpen}
							aria-haspopup="true"
						>
							<span class="sr-only">Open user menu</span>
							<img
								class="size-8 rounded-full bg-gray-50"
								src="https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80"
								alt=""
							/>
							<span class="hidden lg:flex lg:items-center">
								<span class="ml-4 text-sm/6 font-semibold text-gray-900" aria-hidden="true"
									>Tom Cook</span
								>
								<svg
									class="ml-2 size-5 text-gray-400"
									viewBox="0 0 20 20"
									fill="currentColor"
									aria-hidden="true"
									data-slot="icon"
								>
									<path
										fill-rule="evenodd"
										d="M5.22 8.22a.75.75 0 0 1 1.06 0L10 11.94l3.72-3.72a.75.75 0 1 1 1.06 1.06l-4.25 4.25a.75.75 0 0 1-1.06 0L5.22 9.28a.75.75 0 0 1 0-1.06Z"
										clip-rule="evenodd"
									/>
								</svg>
							</span>
						</button>

						<!-- 
              Dropdown menu with transition
            -->
						{#if isUserMenuOpen}
							<div
								id="user-menu"
								class="focus:outline-hidden absolute right-0 z-10 mt-2.5 w-32 origin-top-right rounded-md bg-white py-2 shadow-lg ring-1 ring-gray-900/5 transition duration-100 ease-out"
								class:transform-opacity-100-scale-100={isUserMenuOpen}
								class:transform-opacity-0-scale-95={!isUserMenuOpen}
								role="menu"
								aria-orientation="vertical"
								aria-labelledby="user-menu-button"
								tabindex="-1"
							>
								<a
									href="#"
									class="block px-3 py-1 text-sm/6 text-gray-900"
									role="menuitem"
									tabindex="-1"
									id="user-menu-item-0">Your profile</a
								>
								<a
									href="#"
									class="block px-3 py-1 text-sm/6 text-gray-900"
									role="menuitem"
									tabindex="-1"
									id="user-menu-item-1">Sign out</a
								>
							</div>
						{/if}
					</div>
				</div>
			</div>
		</div>

		<main class="py-10">
			<div class="px-4 sm:px-6 lg:px-8">
				{@render children()}
			</div>
		</main>
	</div>
</div>
