<script lang="ts">
	export let meta = {
		page: 1,
		per_page: 20,
		total: 0,
		total_pages: 1,
		has_next: false,
		has_prev: false
	};

	export let baseUrl = '';

	// Function to generate page URL
	function getPageUrl(page: number): string {
		const url = new URL(baseUrl, window.location.origin);
		url.searchParams.set('page', page.toString());
		url.searchParams.set('per_page', meta.per_page.toString());
		return url.toString();
	}

	// Calculate start and end item numbers
	const startItem = Math.min((meta.page - 1) * meta.per_page + 1, meta.total);
	const endItem = Math.min(meta.page * meta.per_page, meta.total);

	// Determine which page numbers to show
	// We'll show up to 7 pages: current, 3 before, and 3 after (when available)
	let pageNumbers: (number | null)[] = [];

	if (meta.total_pages <= 7) {
		// If 7 or fewer pages, show all
		pageNumbers = Array.from({ length: meta.total_pages }, (_, i) => i + 1);
	} else {
		// More than 7 pages, we need ellipses
		const leftBound = Math.max(1, meta.page - 2);
		const rightBound = Math.min(meta.total_pages, meta.page + 2);

		if (leftBound > 1) {
			pageNumbers.push(1);
			if (leftBound > 2) pageNumbers.push(null); // Left ellipsis
		}

		// Add the pages around current page
		for (let i = leftBound; i <= rightBound; i++) {
			pageNumbers.push(i);
		}

		if (rightBound < meta.total_pages) {
			if (rightBound < meta.total_pages - 1) pageNumbers.push(null); // Right ellipsis
			pageNumbers.push(meta.total_pages);
		}
	}
</script>

<div class="flex items-center justify-between border-t border-gray-200 bg-white px-4 py-3 sm:px-6">
	<!-- Mobile pagination controls -->
	<div class="flex flex-1 justify-between sm:hidden">
		<a
			href={meta.has_prev ? getPageUrl(meta.page - 1) : '#'}
			class={`relative inline-flex items-center rounded-md border border-gray-300 bg-white px-4 py-2 text-sm font-medium ${meta.has_prev ? 'text-gray-700 hover:bg-gray-50' : 'cursor-not-allowed text-gray-300'}`}
			aria-disabled={!meta.has_prev}
		>
			Previous
		</a>
		<a
			href={meta.has_next ? getPageUrl(meta.page + 1) : '#'}
			class={`relative ml-3 inline-flex items-center rounded-md border border-gray-300 bg-white px-4 py-2 text-sm font-medium ${meta.has_next ? 'text-gray-700 hover:bg-gray-50' : 'cursor-not-allowed text-gray-300'}`}
			aria-disabled={!meta.has_next}
		>
			Next
		</a>
	</div>

	<!-- Desktop pagination controls -->
	<div class="hidden sm:flex sm:flex-1 sm:items-center sm:justify-between">
		<div>
			<p class="text-sm text-gray-700">
				Showing
				<span class="font-medium">{startItem}</span>
				to
				<span class="font-medium">{endItem}</span>
				of
				<span class="font-medium">{meta.total}</span>
				results
			</p>
		</div>

		{#if meta.total_pages > 1}
			<div>
				<nav class="isolate inline-flex -space-x-px rounded-md shadow-xs" aria-label="Pagination">
					<!-- Previous page button -->
					<a
						href={meta.has_prev ? getPageUrl(meta.page - 1) : '#'}
						class={`relative inline-flex items-center rounded-l-md px-2 py-2 ${meta.has_prev ? 'text-gray-500 ring-1 ring-gray-300 ring-inset hover:bg-gray-50' : 'cursor-not-allowed text-gray-300 ring-1 ring-gray-200 ring-inset'} focus:z-20 focus:outline-offset-0`}
						aria-disabled={!meta.has_prev}
					>
						<span class="sr-only">Previous</span>
						<svg
							class="size-5"
							viewBox="0 0 20 20"
							fill="currentColor"
							aria-hidden="true"
							data-slot="icon"
						>
							<path
								fill-rule="evenodd"
								d="M11.78 5.22a.75.75 0 0 1 0 1.06L8.06 10l3.72 3.72a.75.75 0 1 1-1.06 1.06l-4.25-4.25a.75.75 0 0 1 0-1.06l4.25-4.25a.75.75 0 0 1 1.06 0Z"
								clip-rule="evenodd"
							/>
						</svg>
					</a>

					<!-- Page numbers -->
					{#each pageNumbers as pageNum, i (i)}
						{#if pageNum === null}
							<!-- Ellipsis -->
							<span
								class="relative inline-flex items-center px-4 py-2 text-sm font-semibold text-gray-700 ring-1 ring-gray-300 ring-inset focus:outline-offset-0"
							>
								...
							</span>
						{:else}
							<!-- Page number -->
							<a
								href={pageNum === meta.page ? '#' : getPageUrl(pageNum)}
								aria-current={pageNum === meta.page ? 'page' : undefined}
								class={`relative inline-flex items-center px-4 py-2 text-sm font-semibold ${
									pageNum === meta.page
										? 'z-10 bg-indigo-600 text-white focus:z-20 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600'
										: 'text-gray-900 ring-1 ring-gray-300 ring-inset hover:bg-gray-50 focus:z-20 focus:outline-offset-0'
								}`}
							>
								{pageNum}
							</a>
						{/if}
					{/each}

					<!-- Next page button -->
					<a
						href={meta.has_next ? getPageUrl(meta.page + 1) : '#'}
						class={`relative inline-flex items-center rounded-r-md px-2 py-2 ${meta.has_next ? 'text-gray-500 ring-1 ring-gray-300 ring-inset hover:bg-gray-50' : 'cursor-not-allowed text-gray-300 ring-1 ring-gray-200 ring-inset'} focus:z-20 focus:outline-offset-0`}
						aria-disabled={!meta.has_next}
					>
						<span class="sr-only">Next</span>
						<svg
							class="size-5"
							viewBox="0 0 20 20"
							fill="currentColor"
							aria-hidden="true"
							data-slot="icon"
						>
							<path
								fill-rule="evenodd"
								d="M8.22 5.22a.75.75 0 0 1 1.06 0l4.25 4.25a.75.75 0 0 1 0 1.06l-4.25 4.25a.75.75 0 0 1-1.06-1.06L11.94 10 8.22 6.28a.75.75 0 0 1 0-1.06Z"
								clip-rule="evenodd"
							/>
						</svg>
					</a>
				</nav>
			</div>
		{/if}
	</div>
</div>
