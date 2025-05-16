<script lang="ts">
	// Define the type for button objects
	type ButtonProps = {
		type: "button";
		label: string;
		action: () => void;
		primary?: boolean; // Optional prop to determine if button should use primary styling
		icon?: any; // Optional icon component
	};

	type LinkProps = {
		type: 'link';
		label: string;
		href: string;
		primary?: boolean; // Optional prop to determine if button should use primary styling
		icon?: any; // Optional icon component
	};

    interface PageProps {
        title: string;
        buttons: (ButtonProps | LinkProps)[]
    }

	// Component props
	let { title = 'Default Title', buttons = [] }: PageProps = $props();
</script>

<div class="md:flex md:items-center md:justify-between">
	<div class="min-w-0 flex-1">
		<h2
			class="text-2xl/7 font-bold capitalize text-gray-900 sm:truncate sm:text-3xl sm:tracking-tight"
		>
			{title}
		</h2>
	</div>
	<div class="mt-4 flex md:ml-4 md:mt-0">
		{#each buttons as button, i}
			{#if button.type == 'button'}
				<button
					type="button"
					class="{i > 0 ? 'ml-3' : ''} inline-flex items-center rounded-md {button.primary
						? 'bg-indigo-600 text-white hover:bg-indigo-700 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600'
						: 'shadow-xs bg-white text-gray-900 ring-1 ring-inset ring-gray-300 hover:bg-gray-50'} px-3 py-2 text-sm font-semibold"
					onclick={button.action}
				>
					{#if button.icon}
						<span class="mr-2">
							<button.icon size={16}></button.icon>
						</span>
					{/if}
					{button.label}
				</button>
			{:else if button.type == 'link'}
				<a
					href={button.href}
					class="{i > 0 ? 'ml-3' : ''} inline-flex items-center rounded-md {button.primary
						? 'bg-indigo-600 text-white hover:bg-indigo-700 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600'
						: 'shadow-xs bg-white text-gray-900 ring-1 ring-inset ring-gray-300 hover:bg-gray-50'} px-3 py-2 text-sm font-semibold"
				>
					{#if button.icon}
						<span class="mr-2">
							<button.icon size={16}></button.icon>
						</span>
					{/if}
					{button.label}
				</a>
			{/if}
		{/each}
	</div>
</div>
