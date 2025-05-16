<script lang="ts">
	import type { PageData } from './$types';
	import { Plus } from '@lucide/svelte';


	let { data }: { data: PageData } = $props();
</script>

<div class="md:flex md:items-center md:justify-between">
  <div class="min-w-0 flex-1">
    <h2 class="text-2xl/7 font-bold text-gray-900 sm:truncate sm:text-3xl sm:tracking-tight">Products</h2>
  </div>
  <div class="mt-4 flex md:mt-0 md:ml-4">
    <a href="/products/create" class="inline-flex items-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-xs ring-1 ring-gray-300 ring-inset hover:bg-gray-50"><Plus size={16}/><span class="ml-3">Create</span></a>
  </div>
</div>
{#if data.products && data.products.length > 0}
	{#each data.products as p}
		<p>{p.name}</p>
	{/each}
{:else}
	<p>No products</p>
{/if}
