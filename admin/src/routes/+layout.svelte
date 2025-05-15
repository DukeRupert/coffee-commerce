<script lang="ts">
	import '../app.css';
	import Navbar from '$lib/navbar.svelte';

	let { children } = $props();
</script>

<Navbar />
<div class="relative isolate px-6 pt-14 lg:px-8">
	{@render children()}
</div>
