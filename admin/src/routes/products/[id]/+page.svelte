<!-- admin/src/routes/products/[id]/+page.svelte -->
<script lang="ts">
	import { goto } from '$app/navigation';
	import { ArrowLeft, Edit, Trash2, Coffee, Package, DollarSign, Calendar } from '@lucide/svelte';
	import PageHeader from '$lib/components/pageHeader.svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	// Enhanced product interface to include variants
	interface ProductVariant {
		id: string;
		active: boolean;
		options: Record<string, any>;
		price: {
			amount: number;
			currency: string;
			type: string;
		};
		price_id: string;
		product_id: string;
		stock_level: number;
		stripe_price_id: string;
		created_at: string;
		updated_at: string;
	}

	interface ProductWithVariants {
		product: any;
		variants: ProductVariant[];
	}

	// Handle the case where we got error or no product
	if (data.error || !data.product) {
		console.error('Product load error:', data.error);
	}

	const productData = data.product as ProductWithVariants;
	const product = productData?.product;
	const variants = productData?.variants || [];

	// Helper functions
	function formatCurrency(amount: number, currency: string = 'USD'): string {
		return new Intl.NumberFormat('en-US', {
			style: 'currency',
			currency: currency
		}).format(amount / 100);
	}

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'long',
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	function getStockLevelClass(stockLevel: number): string {
		if (stockLevel <= 0) {
			return 'bg-red-100 text-red-800';
		} else if (stockLevel < 10) {
			return 'bg-yellow-100 text-yellow-800';
		} else {
			return 'bg-green-100 text-green-800';
		}
	}

	function formatOptions(options: Record<string, any>): string {
		if (!options || Object.keys(options).length === 0) return 'Default';
		
		return Object.entries(options)
			.map(([key, value]) => `${key}: ${value}`)
			.join(', ');
	}

	// Navigation functions
	function handleBack() {
		goto('/products');
	}

	function handleEdit() {
		goto(`/products/${data.productId}/edit`);
	}

	function handleDelete() {
		// TODO: Implement delete confirmation and API call
		console.log('Delete product:', data.productId);
	}
</script>

<svelte:head>
	<title>{product?.name || 'Product'} - Coffee Admin</title>
</svelte:head>

{#if data.error}
	<div class="rounded-md bg-red-50 p-4">
		<div class="flex">
			<div class="flex-shrink-0">
				<svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
					<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.28 7.22a.75.75 0 00-1.06 1.06L8.94 10l-1.72 1.72a.75.75 0 101.06 1.06L10 11.06l1.72 1.72a.75.75 0 101.06-1.06L11.06 10l1.72-1.72a.75.75 0 00-1.06-1.06L10 8.94 8.28 7.22z" clip-rule="evenodd" />
				</svg>
			</div>
			<div class="ml-3">
				<h3 class="text-sm font-medium text-red-800">Error loading product</h3>
				<div class="mt-2 text-sm text-red-700">
					<p>{data.error}</p>
				</div>
			</div>
		</div>
	</div>
{:else if !product}
	<div class="text-center">
		<Coffee size={48} class="mx-auto text-gray-400" />
		<h3 class="mt-2 text-sm font-semibold text-gray-900">Product not found</h3>
		<p class="mt-1 text-sm text-gray-500">The product you're looking for doesn't exist.</p>
		<div class="mt-6">
			<button
				onclick={handleBack}
				class="inline-flex items-center rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500"
			>
				<ArrowLeft class="mr-1.5 -ml-0.5 h-5 w-5" />
				Back to Products
			</button>
		</div>
	</div>
{:else}
	<PageHeader
		title={product.name}
		buttons={[
			{
				type: 'button',
				label: 'Back',
				action: handleBack,
				icon: ArrowLeft
			},
			{
				type: 'button',
				label: 'Edit',
				action: handleEdit,
				primary: true,
				icon: Edit
			}
		]}
	/>

	<div class="mt-8 overflow-hidden bg-white shadow sm:rounded-lg">
		<!-- Product Header -->
		<div class="px-4 py-5 sm:px-6">
			<div class="flex items-center space-x-5">
				<div class="flex-shrink-0">
					{#if product.image_url}
						<img
							class="h-20 w-20 rounded-lg object-cover"
							src={product.image_url}
							alt={product.name}
						/>
					{:else}
						<div class="flex h-20 w-20 items-center justify-center rounded-lg bg-gray-200">
							<Coffee size={32} class="text-gray-500" />
						</div>
					{/if}
				</div>
				<div class="flex-1 min-w-0">
					<h1 class="text-2xl font-bold leading-7 text-gray-900 sm:truncate sm:text-3xl">
						{product.name}
					</h1>
					<div class="mt-1 flex flex-col sm:mt-0 sm:flex-row sm:flex-wrap sm:space-x-6">
						<div class="mt-2 flex items-center text-sm text-gray-500">
							<Package class="mr-1.5 h-4 w-4 flex-shrink-0 text-gray-400" />
							{product.origin || 'Unknown origin'}
						</div>
						<div class="mt-2 flex items-center text-sm text-gray-500">
							<Coffee class="mr-1.5 h-4 w-4 flex-shrink-0 text-gray-400" />
							{product.roast_level || 'Unknown roast'}
						</div>
					</div>
				</div>
				<div class="flex flex-col space-y-2">
					{#if product.active}
						<span class="inline-flex rounded-full bg-green-100 px-2 text-xs font-semibold leading-5 text-green-800">
							Active
						</span>
					{:else}
						<span class="inline-flex rounded-full bg-gray-100 px-2 text-xs font-semibold leading-5 text-gray-800">
							Inactive
						</span>
					{/if}
					{#if product.allow_subscription}
						<span class="inline-flex rounded-full bg-blue-100 px-2 text-xs font-semibold leading-5 text-blue-800">
							Subscription
						</span>
					{/if}
				</div>
			</div>
		</div>

		<!-- Product Details -->
		<div class="border-t border-gray-200 px-4 py-5 sm:px-6">
			<dl class="grid grid-cols-1 gap-x-4 gap-y-6 sm:grid-cols-2">
				<!-- Description -->
				<div class="sm:col-span-2">
					<dt class="text-sm font-medium text-gray-500">Description</dt>
					<dd class="mt-1 text-sm text-gray-900">{product.description}</dd>
				</div>

				<!-- Flavor Notes -->
				{#if product.flavor_notes}
					<div class="sm:col-span-2">
						<dt class="text-sm font-medium text-gray-500">Flavor Notes</dt>
						<dd class="mt-1 text-sm text-gray-900">{product.flavor_notes}</dd>
					</div>
				{/if}

				<!-- Stock Level -->
				<div>
					<dt class="text-sm font-medium text-gray-500">Stock Level</dt>
					<dd class="mt-1">
						<span class="inline-flex rounded-full px-2 text-xs font-semibold leading-5 {getStockLevelClass(product.stock_level)}">
							{product.stock_level} units
						</span>
					</dd>
				</div>

				<!-- Weight -->
				{#if product.weight}
					<div>
						<dt class="text-sm font-medium text-gray-500">Weight</dt>
						<dd class="mt-1 text-sm text-gray-900">{product.weight}g</dd>
					</div>
				{/if}

				<!-- Created At -->
				<div>
					<dt class="text-sm font-medium text-gray-500">Created</dt>
					<dd class="mt-1 text-sm text-gray-900">{formatDate(product.created_at)}</dd>
				</div>

				<!-- Updated At -->
				<div>
					<dt class="text-sm font-medium text-gray-500">Last Updated</dt>
					<dd class="mt-1 text-sm text-gray-900">{formatDate(product.updated_at)}</dd>
				</div>
			</dl>
		</div>

		<!-- Product Variants -->
		{#if variants.length > 0}
			<div class="border-t border-gray-200">
				<div class="px-4 py-5 sm:px-6">
					<h3 class="text-lg font-medium leading-6 text-gray-900">Product Variants</h3>
					<p class="mt-1 max-w-2xl text-sm text-gray-500">
						Different pricing and option combinations for this product.
					</p>
				</div>
				<div class="border-t border-gray-200">
					<div class="overflow-hidden shadow ring-1 ring-black ring-opacity-5">
						<table class="min-w-full divide-y divide-gray-300">
							<thead class="bg-gray-50">
								<tr>
									<th scope="col" class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wide text-gray-500">
										Options
									</th>
									<th scope="col" class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wide text-gray-500">
										Price
									</th>
									<th scope="col" class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wide text-gray-500">
										Stock
									</th>
									<th scope="col" class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wide text-gray-500">
										Status
									</th>
									<th scope="col" class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wide text-gray-500">
										Stripe Price ID
									</th>
								</tr>
							</thead>
							<tbody class="divide-y divide-gray-200 bg-white">
								{#each variants as variant (variant.id)}
									<tr>
										<td class="whitespace-nowrap px-6 py-4 text-sm text-gray-900">
											{formatOptions(variant.options)}
										</td>
										<td class="whitespace-nowrap px-6 py-4 text-sm text-gray-900">
											<div class="flex items-center">
												<DollarSign size={16} class="mr-1 text-gray-400" />
												{formatCurrency(variant.price.amount, variant.price.currency)}
												<span class="ml-2 text-xs text-gray-500">({variant.price.type})</span>
											</div>
										</td>
										<td class="whitespace-nowrap px-6 py-4 text-sm">
											<span class="inline-flex rounded-full px-2 text-xs font-semibold leading-5 {getStockLevelClass(variant.stock_level)}">
												{variant.stock_level} units
											</span>
										</td>
										<td class="whitespace-nowrap px-6 py-4 text-sm">
											{#if variant.active}
												<span class="inline-flex rounded-full bg-green-100 px-2 text-xs font-semibold leading-5 text-green-800">
													Active
												</span>
											{:else}
												<span class="inline-flex rounded-full bg-gray-100 px-2 text-xs font-semibold leading-5 text-gray-800">
													Inactive
												</span>
											{/if}
										</td>
										<td class="whitespace-nowrap px-6 py-4 text-sm text-gray-500 font-mono">
											{variant.stripe_price_id}
										</td>
									</tr>
								{/each}
							</tbody>
						</table>
					</div>
				</div>
			</div>
		{/if}

		<!-- Action Buttons -->
		<div class="border-t border-gray-200 px-4 py-4 sm:px-6">
			<div class="flex justify-end space-x-3">
				<button
					type="button"
					onclick={handleDelete}
					class="inline-flex items-center rounded-md border border-red-300 bg-white px-4 py-2 text-sm font-medium text-red-700 shadow-sm hover:bg-red-50 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2"
				>
					<Trash2 class="mr-2 h-4 w-4" />
					Delete Product
				</button>
				<button
					type="button"
					onclick={handleEdit}
					class="inline-flex items-center rounded-md border border-transparent bg-indigo-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
				>
					<Edit class="mr-2 h-4 w-4" />
					Edit Product
				</button>
			</div>
		</div>
	</div>
{/if}