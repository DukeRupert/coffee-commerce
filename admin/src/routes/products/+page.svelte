<script lang="ts">
	import type { PageData } from './$types';
	import { Plus, Coffee } from '@lucide/svelte';
	import PageHeader from '$lib/components/pageHeader.svelte';
	import Pagination from '$lib/components/pagination.svelte';
	import ListItem from './ListItem.svelte';

	let { data }: { data: PageData } = $props();
	$inspect(data);
</script>

<PageHeader
	title="Products"
	buttons={[
		{
			type: 'link',
			label: 'Add Product',
			href: '/products/create',
			primary: true,
			icon: Plus
		}
	]}
/>

<p class="mt-2 text-sm text-gray-700">
	Manage your coffee products, including basic information, pricing, and inventory levels.
</p>

<div class="mt-8 flow-root">
	<div class="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
		<div class="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
			{#if data.products && data.products.length > 0}
				<table class="min-w-full divide-y divide-gray-300">
					<thead>
						<tr>
							<th
								scope="col"
								class="py-3.5 pr-3 pl-4 text-left text-sm font-semibold text-gray-900 sm:pl-0"
							>
								Product
							</th>
							<th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">
								Origin
							</th>
							<th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">
								Roast
							</th>
							<th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">
								Options
							</th>
							<th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">
								Stock
							</th>
							<th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">
								Status
							</th>
							<th scope="col" class="relative py-3.5 pr-4 pl-3 sm:pr-0">
								<span class="sr-only">Edit</span>
							</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-gray-200">
						{#each data.products as product, index (index)}
							<ListItem {product} />
						{/each}
					</tbody>
				</table>
			{:else}
				<div class="rounded-lg border border-gray-200 bg-white py-12 text-center">
					<Coffee size={48} class="mx-auto text-gray-400" />
					<h3 class="mt-2 text-sm font-semibold text-gray-900">No products</h3>
					<p class="mt-1 text-sm text-gray-500">Get started by creating a new product.</p>
					<div class="mt-6">
						<a
							href="/products/create"
							class="inline-flex items-center rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
						>
							<Plus class="mr-1.5 -ml-0.5 h-5 w-5" />
							New Product
						</a>
					</div>
				</div>
			{/if}
		</div>
	</div>
</div>
<Pagination meta={data.meta} baseUrl="/products" />
