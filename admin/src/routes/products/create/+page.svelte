<script lang="ts">
  import type { PageData } from "./$types";
  import { X } from '@lucide/svelte';
  let { data }: { data: PageData } = $props();

  // Define the product interface to match the API schema
  interface CreateProductRequest {
    name: string;
    description: string;
    image_url?: string;
    origin?: string;
    roast_level?: string;
    flavor_notes?: string;
    stock_level?: number;
    active?: boolean;
    allow_subscription?: boolean;
    options?: {
      weights?: string[];
      grinds?: string[];
    };
  }

  // State for form values
  let product = $state<CreateProductRequest>({
    name: '',
    description: '',
    image_url: '',
    origin: '',
    roast_level: 'Medium',
    flavor_notes: '',
    stock_level: 0,
    active: true,
    allow_subscription: false,
    options: {
      weights: [],
      grinds: []
    }
  });

  // State for form options
  let weightInput = $state('');
  let grindInput = $state('');
  
  // Predefined options
  const roastLevels = ['Light', 'Medium', 'Medium-Dark', 'Dark'];
  
  // Handle adding a weight option
  function addWeight() {
    if (weightInput.trim()) {
      product.options = product.options || {};
      product.options.weights = [...(product.options.weights || []), weightInput.trim()];
      weightInput = '';
    }
  }
  
  // Handle adding a grind option
  function addGrind() {
    if (grindInput.trim()) {
      product.options = product.options || {};
      product.options.grinds = [...(product.options.grinds || []), grindInput.trim()];
      grindInput = '';
    }
  }
  
  // Handle removing a weight option
  function removeWeight(index: number) {
    if (product.options?.weights) {
      product.options.weights = product.options.weights.filter((_, i) => i !== index);
    }
  }
  
  // Handle removing a grind option
  function removeGrind(index: number) {
    if (product.options?.grinds) {
      product.options.grinds = product.options.grinds.filter((_, i) => i !== index);
    }
  }
  
  // Handle form submission
  function handleSubmit(event: SubmitEvent ) {
    // Clone the product object to avoid any reactivity issues
    const productData = JSON.parse(JSON.stringify(product));
    
    // If weights or grinds are empty arrays, remove them
    if (productData.options?.weights?.length === 0) {
      delete productData.options.weights;
    }
    if (productData.options?.grinds?.length === 0) {
      delete productData.options.grinds;
    }
    // If options object is empty, remove it
    if (Object.keys(productData.options || {}).length === 0) {
      delete productData.options;
    }
    
    // Submit the data (implement actual API call)
    console.log('Submitting product data:', productData);
    // Here you would typically call your API endpoint
    // For example: fetch('/api/products', { method: 'POST', body: JSON.stringify(productData) })
  }
  
  // Handle cancel button
  function handleCancel() {
    // Reset form or navigate away
    console.log('Form canceled');
  }
</script>
<div class="md:flex md:items-center md:justify-between">
  <div class="min-w-0 flex-1">
    <h2 class="text-2xl/7 font-bold text-gray-900 sm:truncate sm:text-3xl sm:tracking-tight">Create Product</h2>
  </div>
  <div class="mt-4 flex md:mt-0 md:ml-4">
    <a href="/products" class="inline-flex items-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-xs ring-1 ring-gray-300 ring-inset hover:bg-gray-50"><X size={16}/><span class="ml-3">Cancel</span></a>
  </div>
</div>
<form onsubmit={handleSubmit}>
  <div class="space-y-12">
    <div class="border-b border-gray-900/10 pb-12">
      <h2 class="text-base/7 font-semibold text-gray-900">Coffee Product Information</h2>
      <p class="mt-1 text-sm/6 text-gray-600">Provide details about the coffee product you want to add to your inventory.</p>

      <div class="mt-10 grid grid-cols-1 gap-x-6 gap-y-8 sm:grid-cols-6">
        <!-- Product Name -->
        <div class="sm:col-span-4">
          <label for="name" class="block text-sm/6 font-medium text-gray-900">Product Name</label>
          <div class="mt-2">
            <input 
              type="text" 
              name="name" 
              id="name" 
              required
              bind:value={product.name}
              class="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
              placeholder="Ethiopian Yirgacheffe"
            >
          </div>
        </div>

        <!-- Product Description -->
        <div class="col-span-full">
          <label for="description" class="block text-sm/6 font-medium text-gray-900">Description</label>
          <div class="mt-2">
            <textarea 
              name="description" 
              id="description" 
              rows="3" 
              required
              bind:value={product.description}
              class="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
              placeholder="Describe the flavor profile, processing method, and other notable characteristics"
            ></textarea>
          </div>
          <p class="mt-3 text-sm/6 text-gray-600">Write a detailed description that helps customers understand what makes this coffee special.</p>
        </div>

        <!-- Image URL -->
        <div class="sm:col-span-4">
          <label for="image_url" class="block text-sm/6 font-medium text-gray-900">Image URL</label>
          <div class="mt-2">
            <input 
              type="url" 
              name="image_url" 
              id="image_url" 
              bind:value={product.image_url}
              class="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
              placeholder="https://example.com/images/coffee.jpg"
            >
          </div>
          <p class="mt-1 text-sm/6 text-gray-600">Provide a URL to an image of the coffee product.</p>
        </div>

        <!-- Coffee Origin -->
        <div class="sm:col-span-3">
          <label for="origin" class="block text-sm/6 font-medium text-gray-900">Origin</label>
          <div class="mt-2">
            <input 
              type="text" 
              name="origin" 
              id="origin" 
              bind:value={product.origin}
              class="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
              placeholder="Ethiopia, Yirgacheffe"
            >
          </div>
        </div>

        <!-- Roast Level -->
        <div class="sm:col-span-3">
          <label for="roast_level" class="block text-sm/6 font-medium text-gray-900">Roast Level</label>
          <div class="mt-2 grid grid-cols-1">
            <select 
              id="roast_level" 
              name="roast_level" 
              bind:value={product.roast_level}
              class="col-start-1 row-start-1 w-full appearance-none rounded-md bg-white py-1.5 pr-8 pl-3 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
            >
              {#each roastLevels as level}
                <option value={level}>{level}</option>
              {/each}
            </select>
            <svg class="pointer-events-none col-start-1 row-start-1 mr-2 size-5 self-center justify-self-end text-gray-500 sm:size-4" viewBox="0 0 16 16" fill="currentColor" aria-hidden="true" data-slot="icon">
              <path fill-rule="evenodd" d="M4.22 6.22a.75.75 0 0 1 1.06 0L8 8.94l2.72-2.72a.75.75 0 1 1 1.06 1.06l-3.25 3.25a.75.75 0 0 1-1.06 0L4.22 7.28a.75.75 0 0 1 0-1.06Z" clip-rule="evenodd" />
            </svg>
          </div>
        </div>

        <!-- Flavor Notes -->
        <div class="sm:col-span-4">
          <label for="flavor_notes" class="block text-sm/6 font-medium text-gray-900">Flavor Notes</label>
          <div class="mt-2">
            <input 
              type="text" 
              name="flavor_notes" 
              id="flavor_notes" 
              bind:value={product.flavor_notes}
              class="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
              placeholder="Blueberry, Chocolate, Citrus"
            >
          </div>
          <p class="mt-1 text-sm/6 text-gray-600">Separate flavor notes with commas.</p>
        </div>

        <!-- Stock Level -->
        <div class="sm:col-span-2">
          <label for="stock_level" class="block text-sm/6 font-medium text-gray-900">Stock Level</label>
          <div class="mt-2">
            <input 
              type="number" 
              name="stock_level" 
              id="stock_level" 
              min="0"
              bind:value={product.stock_level}
              class="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
            >
          </div>
        </div>
      </div>
    </div>

    <!-- Product Options Section -->
    <div class="border-b border-gray-900/10 pb-12">
      <h2 class="text-base/7 font-semibold text-gray-900">Product Options</h2>
      <p class="mt-1 text-sm/6 text-gray-600">Configure available weights and grind options for this coffee.</p>

      <div class="mt-10 grid grid-cols-1 gap-x-6 gap-y-8 sm:grid-cols-6">
        <!-- Weight Options -->
        <div class="sm:col-span-3">
          <label for="weight_option" class="block text-sm/6 font-medium text-gray-900">Weight Options</label>
          <div class="mt-2 flex gap-x-2">
            <input 
              type="text" 
              id="weight_option" 
              bind:value={weightInput}
              class="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
              placeholder="e.g. 12oz"
            >
            <button 
              type="button" 
              onclick={addWeight}
              class="rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-xs hover:bg-indigo-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
            >
              Add
            </button>
          </div>
          <div class="mt-2 flex flex-wrap gap-2">
            {#if product.options?.weights && product.options.weights.length > 0}
              {#each product.options.weights as weight, index}
                <div class="inline-flex items-center rounded-md bg-indigo-50 px-2 py-1 text-sm font-medium text-indigo-700">
                  {weight}
                  <button 
                    type="button" 
                    class="ml-1 inline-flex h-4 w-4 flex-shrink-0 items-center justify-center rounded-full text-indigo-400 hover:bg-indigo-200 hover:text-indigo-500 focus:outline-none focus:bg-indigo-500 focus:text-white"
                    onclick={() => removeWeight(index)}
                  >
                    <span class="sr-only">Remove {weight}</span>
                    <svg class="h-2 w-2" stroke="currentColor" fill="none" viewBox="0 0 8 8">
                      <path stroke-linecap="round" stroke-width="1.5" d="M1 1l6 6m0-6L1 7" />
                    </svg>
                  </button>
                </div>
              {/each}
            {/if}
          </div>
        </div>

        <!-- Grind Options -->
        <div class="sm:col-span-3">
          <label for="grind_option" class="block text-sm/6 font-medium text-gray-900">Grind Options</label>
          <div class="mt-2 flex gap-x-2">
            <input 
              type="text" 
              id="grind_option" 
              bind:value={grindInput}
              class="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
              placeholder="e.g. Whole Bean"
            >
            <button 
              type="button" 
              onclick={addGrind}
              class="rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-xs hover:bg-indigo-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
            >
              Add
            </button>
          </div>
          <div class="mt-2 flex flex-wrap gap-2">
            {#if product.options?.grinds && product.options.grinds.length > 0}
              {#each product.options.grinds as grind, index}
                <div class="inline-flex items-center rounded-md bg-indigo-50 px-2 py-1 text-sm font-medium text-indigo-700">
                  {grind}
                  <button 
                    type="button" 
                    class="ml-1 inline-flex h-4 w-4 flex-shrink-0 items-center justify-center rounded-full text-indigo-400 hover:bg-indigo-200 hover:text-indigo-500 focus:outline-none focus:bg-indigo-500 focus:text-white"
                    onclick={() => removeGrind(index)}
                  >
                    <span class="sr-only">Remove {grind}</span>
                    <svg class="h-2 w-2" stroke="currentColor" fill="none" viewBox="0 0 8 8">
                      <path stroke-linecap="round" stroke-width="1.5" d="M1 1l6 6m0-6L1 7" />
                    </svg>
                  </button>
                </div>
              {/each}
            {/if}
          </div>
        </div>
      </div>
    </div>

    <!-- Product Status Section -->
    <div class="border-b border-gray-900/10 pb-12">
      <h2 class="text-base/7 font-semibold text-gray-900">Product Status</h2>
      <p class="mt-1 text-sm/6 text-gray-600">Configure visibility and purchase options for this coffee product.</p>

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
                    checked={product.active}
                    onchange={() => product.active = !product.active}
                    class="col-start-1 row-start-1 appearance-none rounded-sm border border-gray-300 bg-white checked:border-indigo-600 checked:bg-indigo-600 indeterminate:border-indigo-600 indeterminate:bg-indigo-600 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 disabled:border-gray-300 disabled:bg-gray-100 disabled:checked:bg-gray-100 forced-colors:appearance-auto"
                  >
                  <svg class="pointer-events-none col-start-1 row-start-1 size-3.5 self-center justify-self-center stroke-white group-has-disabled:stroke-gray-950/25" viewBox="0 0 14 14" fill="none">
                    <path class="opacity-0 group-has-checked:opacity-100" d="M3 8L6 11L11 3.5" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
                    <path class="opacity-0 group-has-indeterminate:opacity-100" d="M3 7H11" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
                  </svg>
                </div>
              </div>
              <div class="text-sm/6">
                <label for="active" class="font-medium text-gray-900">Active Product</label>
                <p class="text-gray-500">When active, this product will be visible to customers and available for purchase.</p>
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
                    checked={product.allow_subscription}
                    onchange={() => product.allow_subscription = !product.allow_subscription}
                    class="col-start-1 row-start-1 appearance-none rounded-sm border border-gray-300 bg-white checked:border-indigo-600 checked:bg-indigo-600 indeterminate:border-indigo-600 indeterminate:bg-indigo-600 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 disabled:border-gray-300 disabled:bg-gray-100 disabled:checked:bg-gray-100 forced-colors:appearance-auto"
                  >
                  <svg class="pointer-events-none col-start-1 row-start-1 size-3.5 self-center justify-self-center stroke-white group-has-disabled:stroke-gray-950/25" viewBox="0 0 14 14" fill="none">
                    <path class="opacity-0 group-has-checked:opacity-100" d="M3 8L6 11L11 3.5" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
                    <path class="opacity-0 group-has-indeterminate:opacity-100" d="M3 7H11" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
                  </svg>
                </div>
              </div>
              <div class="text-sm/6">
                <label for="allow_subscription" class="font-medium text-gray-900">Allow Subscription</label>
                <p class="text-gray-500">When enabled, customers will be able to purchase this coffee as a recurring subscription.</p>
              </div>
            </div>
          </div>
        </fieldset>
      </div>
    </div>
  </div>

  <div class="mt-6 flex items-center justify-end gap-x-6">
    <button type="button" onclick={handleCancel} class="text-sm/6 font-semibold text-gray-900">Cancel</button>
    <button type="submit" class="rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-xs hover:bg-indigo-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600">Create Product</button>
  </div>
</form>