<!-- admin/src/routes/products/[id]/edit/+page.svelte -->
<script lang="ts">
	import { enhance } from '$app/forms';
	import { goto } from '$app/navigation';
	import { X, Plus, Save, ArrowLeft } from '@lucide/svelte';
	import type { PageData, ActionData } from './$types';

	let { data, form }: { data: PageData; form: ActionData } = $props();

	// Handle the case where we got error or no product
	if (data.error || !data.product) {
		console.error('Product load error:', data.error);
	}

	const product = data.product;

	// Define our option structure for local state management
	interface ProductOption {
		key: string;
		values: string[];
	}

	// Initialize form state with existing product data
	let formData = $state({
		name: product?.name || '',
		description: product?.description || '',
		image_url: product?.image_url || '',
		origin: product?.origin || '',
		roast_level: product?.roast_level || 'Medium',
		flavor_notes: product?.flavor_notes || '',
		stock_level: product?.stock_level || 0,
		weight: product?.weight || 0,
		active: product?.active ?? true,
		allow_subscription: product?.allow_subscription ?? false
	});

	// Initialize options from existing product data
	let productOptions = $state<ProductOption[]>([]);

	// Initialize productOptions when product data is available
	$effect(() => {
		if (product?.options) {
			productOptions = Object.entries(product.options).map(([key, values]) => ({
				key,
				values: Array.isArray(values) ? [...values] : []
			}));
		} else {
			productOptions = [];
		}
	});

	// State for options management
	let currentOptionKey = $state('');
	let currentOptionValue = $state('');
	let selectedOptionIndex = $state<number | null>(null);

	// State for form management
	let error = $state('');
	let isSubmitting = $state(false);
	let validationErrors = $state<Record<string, string>>({});

	// Update validation errors when form action result changes
	$effect(() => {
		if (form) {
			if (form.success && form.redirect) {
				// Handle successful update with redirect
				goto(form.redirect);
			} else if (!form.success && form.validationErrors) {
				validationErrors = form.validationErrors as Record<string, string>;
				error = form.message || 'Validation failed';
			} else if (!form.success) {
				error = form.message || 'Update failed';
				validationErrors = {};
			} else if (form.success) {
				// Clear any previous errors on success
				error = '';
				validationErrors = {};
			}
		}
	});

	// Predefined options
	const roastLevels = ['Light', 'Medium', 'Medium-Dark', 'Dark'];

	// Handle adding a new option
	function addOption() {
		// Ensure productOptions is initialized as an array
		if (!Array.isArray(productOptions)) {
			productOptions = [];
		}

		if (productOptions.length >= 3) {
			error = 'Maximum of 3 options allowed';
			return;
		}

		if (!currentOptionKey.trim()) {
			error = 'Option key is required';
			return;
		}

		// Check if option key already exists
		if (productOptions.some((opt) => opt.key === currentOptionKey)) {
			error = `Option "${currentOptionKey}" already exists`;
			return;
		}

		productOptions.push({
			key: currentOptionKey,
			values: []
		});

		currentOptionKey = '';
		error = '';
	}

	// Handle selecting an option for adding values
	function selectOption(index: number) {
		selectedOptionIndex = index;
		currentOptionValue = '';
	}

	// Handle adding a value to the selected option
	function addValueToOption() {
		if (selectedOptionIndex === null) {
			error = 'Select an option first';
			return;
		}

		if (!currentOptionValue.trim()) {
			error = 'Value cannot be empty';
			return;
		}

		const option = productOptions[selectedOptionIndex];

		// Check if value already exists
		if (option.values.includes(currentOptionValue)) {
			error = `Value "${currentOptionValue}" already exists in this option`;
			return;
		}

		option.values.push(currentOptionValue);
		currentOptionValue = '';
		error = '';
	}

	// Handle removing an option
	function removeOption(index: number) {
		productOptions = productOptions.filter((_, i) => i !== index);
		if (selectedOptionIndex === index) {
			selectedOptionIndex = null;
		} else if (selectedOptionIndex !== null && selectedOptionIndex > index) {
			selectedOptionIndex--;
		}
	}

	// Handle removing a value from an option
	function removeValueFromOption(optionIndex: number, valueIndex: number) {
		productOptions[optionIndex].values = productOptions[optionIndex].values.filter(
			(_, i) => i !== valueIndex
		);
	}

	// Handle cancel button
	function handleCancel() {
		goto(`/products/${data.productId}`);
	}

	// Prepare options for submission
	function prepareOptionsForSubmission(): Record<string, string[]> {
		const options: Record<string, string[]> = {};

		// Ensure productOptions is an array before iterating
		if (Array.isArray(productOptions)) {
			productOptions.forEach((option) => {
				if (option.values && option.values.length > 0) {
					options[option.key] = [...option.values];
				}
			});
		}

		return options;
	}

	// Helper for validation errors
	function hasError(field: string): boolean {
		return Object.keys(validationErrors).some(
			(key) => key === field || key.startsWith(`${field}.`)
		);
	}

	function getError(field: string): string {
		const directError = validationErrors[field];
		if (directError) return directError;

		// Look for nested errors
		const nestedErrors = Object.entries(validationErrors)
			.filter(([key]) => key.startsWith(`${field}.`))
			.map(([key, value]) => value);

		return nestedErrors.join(' ');
	}
</script>

<svelte:head>
	<title>Edit {product?.name || 'Product'} - Coffee Admin</title>
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
		<h3 class="mt-2 text-sm font-semibold text-gray-900">Product not found</h3>
		<p class="mt-1 text-sm text-gray-500">The product you're trying to edit doesn't exist.</p>
	</div>
{:else}
	<div class="md:flex md:items-center md:justify-between">
		<div class="min-w-0 flex-1">
			<h2 class="text-2xl/7 font-bold text-gray-900 sm:truncate sm:text-3xl sm:tracking-tight">
				Edit Product: {product.name}
			</h2>
		</div>
		<div class="mt-4 flex md:mt-0 md:ml-4">
			<button
				type="button"
				onclick={handleCancel}
				class="inline-flex items-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-xs ring-1 ring-gray-300 ring-inset hover:bg-gray-50"
			>
				<ArrowLeft size={16} /><span class="ml-3">Back</span>
			</button>
		</div>
	</div>

	<div class="py-4">
		{#if error}
			<div class="rounded-md border border-red-300 bg-red-50 p-4">
				<p class="text-sm text-red-700">{error}</p>
			</div>
		{/if}

		<form 
			method="POST" 
			action="?/update"
			use:enhance={() => {
				isSubmitting = true;
				return async ({ update, result }) => {
					// Handle different result types
					if (result.type === 'redirect') {
						// Let SvelteKit handle the redirect
						await update();
					} else if (result.type === 'failure') {
						// Handle validation errors
						await update();
					} else {
						// Handle other cases
						await update();
					}
					isSubmitting = false;
				};
			}}
		>
			<!-- Hidden field for options -->
			<input 
				type="hidden" 
				name="options" 
				value={JSON.stringify(prepareOptionsForSubmission())} 
			/>

			<div class="space-y-12">
				<div class="border-b border-gray-900/10 pb-12">
					<h2 class="text-base/7 font-semibold text-gray-900">Coffee Product Information</h2>
					<p class="mt-1 text-sm/6 text-gray-600">
						Update the details for this coffee product.
					</p>

					<div class="mt-10 grid grid-cols-1 gap-x-6 gap-y-8 sm:grid-cols-6">
						<!-- Product Name -->
						<div class="sm:col-span-4">
							<label for="name" class="block text-sm/6 font-medium text-gray-900">Product Name</label>
							<div class="mt-2">
								<input
									type="text"
									id="name"
									name="name"
									required
									bind:value={formData.name}
									class="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
									placeholder="Ethiopian Yirgacheffe"
								/>
							</div>
							{#if hasError('name')}
								<p class="mt-1 text-sm text-red-600">{getError('name')}</p>
							{/if}
						</div>

						<!-- Product Description -->
						<div class="col-span-full">
							<label for="description" class="block text-sm/6 font-medium text-gray-900"
								>Description</label
							>
							<div class="mt-2">
								<textarea
									id="description"
									name="description"
									rows="3"
									required
									bind:value={formData.description}
									class="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
									placeholder="Describe the flavor profile, processing method, and other notable characteristics"
								></textarea>
							</div>
							{#if hasError('description')}
								<p class="mt-1 text-sm text-red-600">{getError('description')}</p>
							{/if}
						</div>

						<!-- Image URL -->
						<div class="sm:col-span-4">
							<label for="image_url" class="block text-sm/6 font-medium text-gray-900">Image URL</label>
							<div class="mt-2">
								<input
									type="url"
									id="image_url"
									name="image_url"
									bind:value={formData.image_url}
									class="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
									placeholder="https://example.com/images/coffee.jpg"
								/>
							</div>
							{#if hasError('image_url')}
								<p class="mt-1 text-sm text-red-600">{getError('image_url')}</p>
							{/if}
						</div>

						<!-- Coffee Origin -->
						<div class="sm:col-span-3">
							<label for="origin" class="block text-sm/6 font-medium text-gray-900">Origin</label>
							<div class="mt-2">
								<input
									type="text"
									id="origin"
									name="origin"
									bind:value={formData.origin}
									class="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
									placeholder="Ethiopia, Yirgacheffe"
								/>
							</div>
							{#if hasError('origin')}
								<p class="mt-1 text-sm text-red-600">{getError('origin')}</p>
							{/if}
						</div>

						<!-- Roast Level -->
						<div class="sm:col-span-3">
							<label for="roast_level" class="block text-sm/6 font-medium text-gray-900"
								>Roast Level</label
							>
							<div class="mt-2 grid grid-cols-1">
								<select
									id="roast_level"
									name="roast_level"
									bind:value={formData.roast_level}
									class="col-start-1 row-start-1 w-full appearance-none rounded-md bg-white py-1.5 pr-8 pl-3 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
								>
									{#each roastLevels as level, index (index)}
										<option value={level}>{level}</option>
									{/each}
								</select>
								<svg
									class="pointer-events-none col-start-1 row-start-1 mr-2 size-5 self-center justify-self-end text-gray-500 sm:size-4"
									viewBox="0 0 16 16"
									fill="currentColor"
									aria-hidden="true"
									data-slot="icon"
								>
									<path
										fill-rule="evenodd"
										d="M4.22 6.22a.75.75 0 0 1 1.06 0L8 8.94l2.72-2.72a.75.75 0 1 1 1.06 1.06l-3.25 3.25a.75.75 0 0 1-1.06 0L4.22 7.28a.75.75 0 0 1 0-1.06Z"
										clip-rule="evenodd"
									/>
								</svg>
							</div>
							{#if hasError('roast_level')}
								<p class="mt-1 text-sm text-red-600">{getError('roast_level')}</p>
							{/if}
						</div>

						<!-- Flavor Notes -->
						<div class="sm:col-span-4">
							<label for="flavor_notes" class="block text-sm/6 font-medium text-gray-900"
								>Flavor Notes</label
							>
							<div class="mt-2">
								<input
									type="text"
									id="flavor_notes"
									name="flavor_notes"
									bind:value={formData.flavor_notes}
									class="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
									placeholder="Blueberry, Chocolate, Citrus"
								/>
							</div>
							<p class="mt-1 text-sm/6 text-gray-600">Separate flavor notes with commas.</p>
							{#if hasError('flavor_notes')}
								<p class="mt-1 text-sm text-red-600">{getError('flavor_notes')}</p>
							{/if}
						</div>

						<!-- Stock Level -->
						<div class="sm:col-span-3">
							<label for="stock_level" class="block text-sm/6 font-medium text-gray-900"
								>Stock Level</label
							>
							<div class="mt-2">
								<input
									type="number"
									id="stock_level"
									name="stock_level"
									min="0"
									bind:value={formData.stock_level}
									class="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
								/>
							</div>
							{#if hasError('stock_level')}
								<p class="mt-1 text-sm text-red-600">{getError('stock_level')}</p>
							{/if}
						</div>

						<!-- Weight -->
						<div class="sm:col-span-3">
							<label for="weight" class="block text-sm/6 font-medium text-gray-900"
								>Weight (grams)</label
							>
							<div class="mt-2">
								<input
									type="number"
									id="weight"
									name="weight"
									min="1"
									bind:value={formData.weight}
									class="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
								/>
							</div>
							{#if hasError('weight')}
								<p class="mt-1 text-sm text-red-600">{getError('weight')}</p>
							{/if}
						</div>
					</div>
				</div>

				<!-- Product Options Section -->
				<div class="border-b border-gray-900/10 pb-12">
					<div class="flex items-center justify-between">
						<div>
							<h2 class="text-base/7 font-semibold text-gray-900">Product Options</h2>
							<p class="mt-1 text-sm/6 text-gray-600">
								Configure custom options for this coffee (max 3). Common examples might be "weight" or "grind".
							</p>
						</div>
						<div class="text-sm text-gray-500">
							{Array.isArray(productOptions) ? productOptions.length : 0}/3 options
						</div>
					</div>

					<div class="mt-6 space-y-6">
						<!-- Add new option -->
						<div class="sm:col-span-full">
							<label for="option_key" class="block text-sm/6 font-medium text-gray-900"
								>Add New Option</label
							>
							<div class="mt-2 flex gap-x-2">
								<input
									type="text"
									id="option_key"
									bind:value={currentOptionKey}
									placeholder="Option name (e.g. weight, grind, size)"
									disabled={(Array.isArray(productOptions) ? productOptions.length : 0) >= 3}
									class="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 disabled:bg-gray-100 disabled:text-gray-500 sm:text-sm/6"
								/>
								<button
									type="button"
									onclick={addOption}
									disabled={(Array.isArray(productOptions) ? productOptions.length : 0) >= 3 || !currentOptionKey.trim()}
									class="rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-xs hover:bg-indigo-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 disabled:bg-gray-300 disabled:hover:bg-gray-300"
								>
									Add Option
								</button>
							</div>
							{#if hasError('options')}
								<p class="mt-1 text-sm text-red-600">{getError('options')}</p>
							{/if}
						</div>

						<!-- Option list -->
						{#if Array.isArray(productOptions) && productOptions.length > 0}
							<div class="mt-4 space-y-4">
								<h3 class="text-sm font-medium text-gray-700">Your Options</h3>
								<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
									{#each productOptions as option, optionIndex (optionIndex)}
										<div
											class="rounded-md border p-4 {selectedOptionIndex === optionIndex
												? 'border-indigo-500 bg-indigo-50'
												: 'border-gray-300'}"
										>
											<div class="mb-3 flex items-center justify-between">
												<h4 class="font-medium text-gray-900 capitalize">{option.key}</h4>
												<button
													type="button"
													onclick={() => removeOption(optionIndex)}
													class="text-red-500 hover:text-red-700"
												>
													<X size={16} />
												</button>
											</div>

											<!-- Values list -->
											<div class="mb-3 flex flex-wrap gap-2">
												{#if option.values.length === 0}
													<p class="text-sm text-gray-500 italic">No values added yet</p>
												{:else}
													{#each option.values as value, valueIndex (valueIndex)}
														<div
															class="inline-flex items-center rounded-md bg-gray-100 px-2 py-1 text-sm font-medium text-gray-700"
														>
															{value}
															<button
																type="button"
																class="ml-1 inline-flex h-4 w-4 flex-shrink-0 items-center justify-center rounded-full text-gray-400 hover:bg-gray-200 hover:text-gray-500"
																onclick={() => removeValueFromOption(optionIndex, valueIndex)}
															>
																<span class="sr-only">Remove {value}</span>
																<svg
																	class="h-2 w-2"
																	stroke="currentColor"
																	fill="none"
																	viewBox="0 0 8 8"
																>
																	<path
																		stroke-linecap="round"
																		stroke-width="1.5"
																		d="M1 1l6 6m0-6L1 7"
																	/>
																</svg>
															</button>
														</div>
													{/each}
												{/if}
											</div>

											<!-- Add values input -->
											<div class="flex gap-2">
												<button
													type="button"
													onclick={() => selectOption(optionIndex)}
													class="text-sm text-indigo-600 underline underline-offset-2 hover:text-indigo-800"
												>
													{selectedOptionIndex === optionIndex ? 'Adding values...' : 'Add values'}
												</button>
											</div>
										</div>
									{/each}
								</div>
							</div>
						{/if}

						<!-- Add values to selected option -->
						{#if selectedOptionIndex !== null}
							<div class="mt-4 rounded-md border border-indigo-300 bg-indigo-50 p-4">
								<h3 class="mb-2 text-sm font-medium text-gray-700">
									Add values to "{productOptions[selectedOptionIndex].key}"
								</h3>
								<div class="flex gap-x-2">
									<input
										type="text"
										bind:value={currentOptionValue}
										placeholder="Enter a value"
										class="block w-full rounded-md bg-white px-3 py-1.5 text-sm text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600"
									/>
									<button
										type="button"
										onclick={addValueToOption}
										disabled={!currentOptionValue.trim()}
										class="rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-xs hover:bg-indigo-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 disabled:bg-gray-300 disabled:hover:bg-gray-300"
									>
										Add
									</button>
								</div>
							</div>
						{/if}
					</div>
				</div>

				<!-- Product Status Section -->
				<div class="border-b border-gray-900/10 pb-12">
					<h2 class="text-base/7 font-semibold text-gray-900">Product Status</h2>
					<p class="mt-1 text-sm/6 text-gray-600">
						Configure visibility and purchase options for this coffee product.
					</p>

					<div class="mt-10 space-y-10">
						<fieldset>
							<div class="mt-6 space-y-6">
								<!-- Active Status -->
								<div class="flex gap-3">
									<div class="flex h-6 shrink-0 items-center">
										<div class="group grid size-4 grid-cols-1">
											<input
												id="active"
												name="active"
												type="checkbox"
												bind:checked={formData.active}
												class="col-start-1 row-start-1 appearance-none rounded-sm border border-gray-300 bg-white checked:border-indigo-600 checked:bg-indigo-600 indeterminate:border-indigo-600 indeterminate:bg-indigo-600 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 disabled:border-gray-300 disabled:bg-gray-100 disabled:checked:bg-gray-100 forced-colors:appearance-auto"
											/>
											<svg
												class="pointer-events-none col-start-1 row-start-1 size-3.5 self-center justify-self-center stroke-white group-has-disabled:stroke-gray-950/25"
												viewBox="0 0 14 14"
												fill="none"
											>
												<path
													class="opacity-0 group-has-checked:opacity-100"
													d="M3 8L6 11L11 3.5"
													stroke-width="2"
													stroke-linecap="round"
													stroke-linejoin="round"
												/>
												<path
													class="opacity-0 group-has-indeterminate:opacity-100"
													d="M3 7H11"
													stroke-width="2"
													stroke-linecap="round"
													stroke-linejoin="round"
												/>
											</svg>
										</div>
									</div>
									<div class="text-sm/6">
										<label for="active" class="font-medium text-gray-900">Active Product</label>
										<p class="text-gray-500">
											When active, this product will be visible to customers and available for purchase.
										</p>
									</div>
								</div>

								<!-- Allow Subscription -->
								<div class="flex gap-3">
									<div class="flex h-6 shrink-0 items-center">
										<div class="group grid size-4 grid-cols-1">
											<input
												id="allow_subscription"
												name="allow_subscription"
												type="checkbox"
												bind:checked={formData.allow_subscription}
												class="col-start-1 row-start-1 appearance-none rounded-sm border border-gray-300 bg-white checked:border-indigo-600 checked:bg-indigo-600 indeterminate:border-indigo-600 indeterminate:bg-indigo-600 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 disabled:border-gray-300 disabled:bg-gray-100 disabled:checked:bg-gray-100 forced-colors:appearance-auto"
											/>
											<svg
												class="pointer-events-none col-start-1 row-start-1 size-3.5 self-center justify-self-center stroke-white group-has-disabled:stroke-gray-950/25"
												viewBox="0 0 14 14"
												fill="none"
											>
												<path
													class="opacity-0 group-has-checked:opacity-100"
													d="M3 8L6 11L11 3.5"
													stroke-width="2"
													stroke-linecap="round"
													stroke-linejoin="round"
												/>
												<path
													class="opacity-0 group-has-indeterminate:opacity-100"
													d="M3 7H11"
													stroke-width="2"
													stroke-linecap="round"
													stroke-linejoin="round"
												/>
											</svg>
										</div>
									</div>
									<div class="text-sm/6">
										<label for="allow_subscription" class="font-medium text-gray-900"
											>Allow Subscription</label
										>
										<p class="text-gray-500">
											When enabled, customers will be able to purchase this coffee as a recurring
											subscription.
										</p>
									</div>
								</div>
							</div>
						</fieldset>
					</div>
				</div>
			</div>

			<div class="mt-6 flex items-center justify-end gap-x-6">
				<button 
					type="button" 
					onclick={handleCancel} 
					class="text-sm/6 font-semibold text-gray-900"
				>
					Cancel
				</button>
				<button
					type="submit"
					disabled={isSubmitting}
					class="rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-xs hover:bg-indigo-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 disabled:bg-indigo-300 disabled:hover:bg-indigo-300"
				>
					<div class="flex items-center">
						<Save size={16} class="mr-2" />
						{isSubmitting ? 'Updating...' : 'Update Product'}
					</div>
				</button>
			</div>
		</form>
	</div>
{/if}