import { fail } from '@sveltejs/kit';
import type { PageServerLoad, Actions } from './$types';
import client from '$lib/client'

export const load = (async () => {
    // You could fetch data needed for the form here, such as:
    // - Product categories
    // - Available roast levels
    // - Inventory statuses
    // - etc.

    return {
        roastLevels: ['Light', 'Medium', 'Medium-Dark', 'Dark', 'French'],
        grindOptions: ['Whole Bean', 'Drip Ground', 'French Press', 'Espresso'],
        weightOptions: ['12oz', '3lb', '5lb']
    };
}) satisfies PageServerLoad;

export const actions = {
    createProduct: async ({ request }) => {
        // Get form data
        const formData = await request.formData();

        // Basic validation function
        const validateRequired = (field: FormDataEntryValue | null, fieldName: string) => {
            if (!field || (typeof field === 'string' && field.trim() === '')) {
                return { valid: false, error: `${fieldName} is required` };
            }
            return { valid: true };
        };

        // Extract and validate fields
        const name = formData.get('name');
        const nameValidation = validateRequired(name, 'Product name');
        if (!nameValidation.valid) {
            return fail(400, { error: nameValidation.error, values: Object.fromEntries(formData) });
        }

        const description = formData.get('description');
        const descValidation = validateRequired(description, 'Description');
        if (!descValidation.valid) {
            return fail(400, { error: descValidation.error, values: Object.fromEntries(formData) });
        }

        // Get the rest of the fields
        const origin = formData.get('origin');
        const roastLevel = formData.get('roast_level');
        const flavorNotes = formData.get('flavor_notes');
        const imageUrl = formData.get('image_url');
        const stockLevel = parseInt(formData.get('stock_level')?.toString() || '0', 10);
        const active = formData.get('active') === 'on' || formData.get('active') === 'true';
        const allowSubscription = formData.get('allow_subscription') === 'on' || formData.get('allow_subscription') === 'true';

        // Get options as arrays (these will come from checkboxes or multi-selects)
        const weightOptions = formData.getAll('weight_options').map(option => option.toString());
        const grindOptions = formData.getAll('grind_options').map(option => option.toString());

        // Create options object
        const options = {
            weights: weightOptions,
            grinds: grindOptions
        };

        // Create the product data object
        const productData = {
            name,
            description,
            origin,
            roast_level: roastLevel,
            flavor_notes: flavorNotes,
            image_url: imageUrl,
            stock_level: stockLevel,
            active,
            allow_subscription: allowSubscription,
            options
        };

        try {
      // Call the client method
      const result = await client.createProduct(productData);
      return { success: true, product: result };
    } catch (error) {
      return { success: false, error: error.message };
    }
    }
} satisfies Actions;