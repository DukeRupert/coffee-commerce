<script lang="ts">
	import type { PageProps } from './$types';
	import { enhance } from '$app/forms';

	// Data provided by the server
	let { data, form }: PageProps = $props();

	// Destructure data for convenience
	let { roastLevels, grindOptions, weightOptions } = $derived(data);

	// Form state
	let isSubmitting = $state(false);
	let showSuccessMessage = $state(false);
</script>

<svelte:head>
	<title>Create New Coffee Product</title>
</svelte:head>

<div class="container">
	<h1>Create New Coffee Product</h1>

	{#if showSuccessMessage}
		<div class="alert success">
			<p>Product created successfully!</p>
		</div>
	{/if}

	{#if form?.error}
		<div class="alert error">
			<p>{form.error}</p>
		</div>
	{/if}

	<form
		method="POST"
		use:enhance={({ formElement, formData, action, cancel, result }) => {
			isSubmitting = true;

			// Function runs when submission is complete
			return async ({ result, update }) => {
				isSubmitting = false;

				if (result.type === 'redirect') {
					// We're redirecting, meaning success
					showSuccessMessage = true;
					setTimeout(() => {
						showSuccessMessage = false;
					}, 5000);
				} else if (result.type === 'failure') {
					// Error case - let the form validation handle this
					console.error('Form submission failed:', result.data?.error);
				}

				// Update the form
				await update();
			};
		}}
		action="?/createProduct"
	>
		<div class="form-grid">
			<!-- Basic Information Section -->
			<div class="form-section">
				<h2>Basic Information</h2>

				<div class="form-field">
					<label for="name">Product Name <span class="required">*</span></label>
					<input
						type="text"
						id="name"
						name="name"
						required
						value={form?.values?.name || ''}
						placeholder="Ethiopia Yirgacheffe"
					/>
				</div>

				<div class="form-field">
					<label for="description">Description <span class="required">*</span></label>
					<textarea
						id="description"
						name="description"
						required
						rows="5"
						placeholder="Describe the coffee's taste profile, origins, etc."
						>{form?.values?.description || ''}</textarea
					>
				</div>

				<div class="form-field">
					<label for="image_url">Image URL</label>
					<input
						type="url"
						id="image_url"
						name="image_url"
						value={form?.values?.image_url || ''}
						placeholder="https://example.com/coffee-image.jpg"
					/>
				</div>
			</div>

			<!-- Coffee Details Section -->
			<div class="form-section">
				<h2>Coffee Details</h2>

				<div class="form-field">
					<label for="origin">Origin</label>
					<input
						type="text"
						id="origin"
						name="origin"
						value={form?.values?.origin || ''}
						placeholder="Ethiopia, Yirgacheffe region"
					/>
				</div>

				<div class="form-field">
					<label for="roast_level">Roast Level</label>
					<select id="roast_level" name="roast_level">
						<option value="" disabled selected={!form?.values?.roast_level}
							>Select a roast level</option
						>
						{#each roastLevels as level}
							<option value={level} selected={form?.values?.roast_level === level}>{level}</option>
						{/each}
					</select>
				</div>

				<div class="form-field">
					<label for="flavor_notes">Flavor Notes</label>
					<input
						type="text"
						id="flavor_notes"
						name="flavor_notes"
						value={form?.values?.flavor_notes || ''}
						placeholder="Blueberry, Chocolate, Citrus"
					/>
				</div>
			</div>

			<!-- Inventory Section -->
			<div class="form-section">
				<h2>Inventory & Options</h2>

				<div class="form-field">
					<label for="stock_level">Stock Level</label>
					<input
						type="number"
						id="stock_level"
						name="stock_level"
						min="0"
						value={form?.values?.stock_level || '0'}
					/>
				</div>

				<div class="form-field checkbox">
					<input
						type="checkbox"
						id="active"
						name="active"
						checked={form?.values?.active === 'on' || form?.values?.active === 'true'}
					/>
					<label for="active">Active (visible to customers)</label>
				</div>

				<div class="form-field checkbox">
					<input
						type="checkbox"
						id="allow_subscription"
						name="allow_subscription"
						checked={form?.values?.allow_subscription === 'on' ||
							form?.values?.allow_subscription === 'true'}
					/>
					<label for="allow_subscription">Allow Subscription</label>
				</div>

				<div class="form-field">
					<label>Available Weight Options</label>
					<div class="checkbox-group">
						{#each weightOptions as weight}
							<div class="checkbox">
								<input
									type="checkbox"
									id="weight_{weight}"
									name="weight_options"
									value={weight}
									checked={form?.values?.weight_options?.includes(weight)}
								/>
								<label for="weight_{weight}">{weight}</label>
							</div>
						{/each}
					</div>
				</div>

				<div class="form-field">
					<label>Available Grind Options</label>
					<div class="checkbox-group">
						{#each grindOptions as grind}
							<div class="checkbox">
								<input
									type="checkbox"
									id="grind_{grind.replace(' ', '_')}"
									name="grind_options"
									value={grind}
									checked={form?.values?.grind_options?.includes(grind)}
								/>
								<label for="grind_{grind.replace(' ', '_')}">{grind}</label>
							</div>
						{/each}
					</div>
				</div>
			</div>
		</div>

		<div class="form-actions">
			<a href="/products" class="button secondary">Cancel</a>
			<button type="submit" class="button primary" disabled={isSubmitting}>
				{isSubmitting ? 'Creating...' : 'Create Product'}
			</button>
		</div>
	</form>
</div>

<style>
	.container {
		max-width: 1200px;
		margin: 0 auto;
		padding: 2rem;
	}

	h1 {
		margin-bottom: 2rem;
		color: #333;
	}

	.form-grid {
		display: grid;
		grid-template-columns: 1fr;
		gap: 2rem;
	}

	@media (min-width: 768px) {
		.form-grid {
			grid-template-columns: repeat(2, 1fr);
		}

		.form-section:first-child {
			grid-column: 1 / 3;
		}
	}

	.form-section {
		background: #f9f9f9;
		padding: 1.5rem;
		border-radius: 8px;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
	}

	.form-section h2 {
		margin-top: 0;
		margin-bottom: 1.5rem;
		font-size: 1.2rem;
		color: #555;
		border-bottom: 1px solid #ddd;
		padding-bottom: 0.5rem;
	}

	.form-field {
		margin-bottom: 1.25rem;
	}

	.form-field label {
		display: block;
		margin-bottom: 0.5rem;
		font-weight: 500;
		color: #333;
	}

	.form-field input[type='text'],
	.form-field input[type='url'],
	.form-field input[type='number'],
	.form-field select,
	.form-field textarea {
		width: 100%;
		padding: 0.75rem;
		border: 1px solid #ddd;
		border-radius: 4px;
		font-size: 1rem;
	}

	.form-field textarea {
		resize: vertical;
	}

	.checkbox {
		display: flex;
		align-items: center;
	}

	.checkbox input {
		margin-right: 0.5rem;
	}

	.checkbox-group {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
		gap: 0.5rem;
	}

	.required {
		color: #e53e3e;
	}

	small {
		display: block;
		margin-top: 0.25rem;
		color: #666;
		font-size: 0.75rem;
	}

	.form-actions {
		margin-top: 2rem;
		display: flex;
		justify-content: flex-end;
		gap: 1rem;
	}

	.button {
		padding: 0.75rem 1.5rem;
		border-radius: 4px;
		font-size: 1rem;
		cursor: pointer;
		text-decoration: none;
		text-align: center;
		transition:
			background-color 0.2s,
			color 0.2s;
		display: inline-flex;
		align-items: center;
		justify-content: center;
		min-width: 120px;
	}

	.button.primary {
		background-color: #3b82f6;
		color: white;
		border: none;
	}

	.button.primary:hover {
		background-color: #2563eb;
	}

	.button.primary:disabled {
		background-color: #93c5fd;
		cursor: not-allowed;
	}

	.button.secondary {
		background-color: #e5e7eb;
		color: #4b5563;
		border: none;
	}

	.button.secondary:hover {
		background-color: #d1d5db;
	}

	.alert {
		margin-bottom: 2rem;
		padding: 1rem;
		border-radius: 4px;
	}

	.alert.error {
		background-color: #fee2e2;
		color: #b91c1c;
		border: 1px solid #fecaca;
	}

	.alert.success {
		background-color: #dcfce7;
		color: #166534;
		border: 1px solid #bbf7d0;
	}
</style>
