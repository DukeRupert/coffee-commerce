<script lang="ts">
	import '../app.css';
	import Sidebar from '$lib/components/sidebar.svelte';

	let { children } = $props();
</script>

<Sidebar>
	{@render children()}
</Sidebar>
