<script lang="ts">
	import { Coffee, Edit } from '@lucide/svelte';
	import type { Product } from '$lib/coffeServiceClient/main';
	let { product }: { product: Product } = $props();

	// Generate badge color class based on stock level
	function getStockLevelClass(stockLevel: number): string {
		if (stockLevel <= 0) {
			return 'bg-red-100 text-red-800';
		} else if (stockLevel < 10) {
			return 'bg-yellow-100 text-yellow-800';
		} else {
			return 'bg-green-100 text-green-800';
		}
	}

	// Format options for display
	function formatOptions(options: Record<string, string[]> | undefined): string {
		if (!options) return 'None';

		const parts = [];
		if (options.weights && options.weights.length > 0) {
			parts.push(`${options.weights.length} weights`);
		}
		if (options.grinds && options.grinds.length > 0) {
			parts.push(`${options.grinds.length} grinds`);
		}

		return parts.length > 0 ? parts.join(', ') : 'None';
	}

	// Format currency
	function formatCurrency(amount: number): string {
		return new Intl.NumberFormat('en-US', {
			style: 'currency',
			currency: 'USD'
		}).format(amount / 100);
	}
</script>

<tr>
	<td class="py-4 pr-3 pl-4 text-sm sm:pl-0">
		<div class="flex items-center">
			<div class="h-10 w-10 flex-shrink-0">
				{#if product.image_url}
					<img
						class="h-10 w-10 rounded-full object-cover"
						src={product.image_url}
						alt={product.name}
					/>
				{:else}
					<div class="flex h-10 w-10 items-center justify-center rounded-full bg-gray-200">
						<Coffee size={20} class="text-gray-500" />
					</div>
				{/if}
			</div>
			<div class="ml-4">
				<div class="font-medium text-gray-900">{product.name}</div>
				<div class="max-w-xs truncate text-gray-500">
					{product.description.length > 60
						? product.description.substring(0, 60) + '...'
						: product.description}
				</div>
			</div>
		</div>
	</td>
	<td class="px-3 py-4 text-sm text-gray-500">{product.origin || 'N/A'}</td>
	<td class="px-3 py-4 text-sm text-gray-500">{product.roast_level || 'N/A'}</td>
	<td class="px-3 py-4 text-sm text-gray-500">
		{formatOptions(product.options)}
	</td>
	<td class="px-3 py-4 text-sm whitespace-nowrap">
		<span
			class="inline-flex rounded-full px-2 text-xs leading-5 font-semibold {getStockLevelClass(
				product.stock_level
			)}"
		>
			{product.stock_level} units
		</span>
	</td>
	<td class="px-3 py-4 text-sm whitespace-nowrap">
		{#if product.active}
			<span
				class="inline-flex rounded-full bg-green-100 px-2 text-xs leading-5 font-semibold text-green-800"
			>
				Active
			</span>
		{:else}
			<span
				class="inline-flex rounded-full bg-gray-100 px-2 text-xs leading-5 font-semibold text-gray-800"
			>
				Inactive
			</span>
		{/if}
	</td>
	<td class="px-3 py-4 text-sm whitespace-nowrap">
		{#if product.allow_subscription}
			<span
				class="inline-flex rounded-full bg-blue-100 px-2 text-xs leading-5 font-semibold text-blue-800"
			>
				Subscription
			</span>
		{:else}
			<span
				class="inline-flex rounded-full bg-gray-200 px-2 text-xs leading-5 font-semibold text-gray-700"
			>
				One-time
			</span>
		{/if}
	</td>
	<td class="relative py-4 pr-4 pl-3 text-right text-sm font-medium whitespace-nowrap sm:pr-0">
		<a
			href={`/products/${product.id}/edit`}
			class="flex items-center justify-end text-indigo-600 hover:text-indigo-900"
		>
			<Edit size={16} class="mr-1" />
			<span>Edit</span>
			<span class="sr-only">, {product.name}</span>
		</a>
	</td>
</tr>
