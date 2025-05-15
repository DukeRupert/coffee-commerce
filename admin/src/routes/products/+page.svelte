<script lang="ts">
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
	let { products, meta, filters } = $derived(data);
    $inspect(products, meta, filters)
</script>

<h1>Product Page</h1>

{#if products && products.length > 0}
	{#each products as p}
		<p>{p.name}</p>
	{/each}
{:else}
	<p>No products</p>
{/if}
