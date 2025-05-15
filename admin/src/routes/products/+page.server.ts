import type { PageServerLoad, Actions } from './$types';
import client from "$lib/client"
import { fail, redirect } from '@sveltejs/kit';
import type { CreateProductRequest } from '$lib/coffeServiceClient/interface';

/**
 * Page load function for the products page
 * 
 * This function will:
 * 1. Extract filter parameters from the URL query string
 * 2. Call the API with those parameters
 * 3. Return the products and meta data for the page
 */
export const load: PageServerLoad = async ({ url }) => {
    try {
        // Get all query parameters from the URL
        const queryParams = Object.fromEntries(url.searchParams);

        // Convert URL parameters to properly typed filter parameters
        const filterParams = client.extractProductParams(queryParams);

        // Call the API with the extracted parameters
        const response = await client.getProducts(filterParams);
        console.log(response)

        return {
            // Return the products array
            products: response?.data,
            // Return pagination metadata
            meta: response.meta,
            // Return the original filter parameters for UI state
            filters: filterParams
        };
    } catch (error) {
        console.error('Error loading products:', error);
        // Return empty data in case of error
        return {
            products: [],
            meta: {
                page: 1,
                per_page: 20,
                total: 0,
                total_pages: 0,
                has_next: false,
                has_prev: false
            },
            filters: {},
            error: error instanceof Error ? error.message : 'Unknown error'
        };
    }
};

export const actions = {
    default: async ({ request }) => {
        // Get form data
        const formData = await request.formData();

        // Validate required fields
        const name = formData.get('name')?.toString();
        if (!name || name.trim() === '') {
            return fail(400, {
                error: 'Product name is required',
                values: Object.fromEntries(formData)
            });
        }

        const description = formData.get('description')?.toString();
        if (!description || description.trim() === '') {
            return fail(400, {
                error: 'Description is required',
                values: Object.fromEntries(formData)
            });
        }

        const stripeId = formData.get('stripe_id')?.toString();
        if (!stripeId || stripeId.trim() === '') {
            return fail(400, {
                error: 'Stripe ID is required',
                values: Object.fromEntries(formData)
            });
        }

        // Extract other form fields
        const weightOptions = formData.getAll('weight_options').map(option => option.toString());
        const grindOptions = formData.getAll('grind_options').map(option => option.toString());

        // Build product data object
        const productData: CreateProductRequest = {
            name,
            description,
            stripe_id: stripeId,
            origin: formData.get('origin')?.toString(),
            roast_level: formData.get('roast_level')?.toString(),
            flavor_notes: formData.get('flavor_notes')?.toString(),
            image_url: formData.get('image_url')?.toString(),
            stock_level: parseInt(formData.get('stock_level')?.toString() || '0', 10),
            active: formData.get('active') === 'on' || formData.get('active') === 'true',
            allow_subscription: formData.get('allow_subscription') === 'on' || formData.get('allow_subscription') === 'true',
            options: {
                weights: weightOptions.length > 0 ? weightOptions : undefined,
                grinds: grindOptions.length > 0 ? grindOptions : undefined
            }
        };

        try {
            // Call the client createProduct method
            const result = await client.createProduct(productData);

            // Redirect to the product list page with a success notification
            return redirect(303, '/products?created=true');
        } catch (error) {
            console.error('Failed to create product:', error);

            // Return error information to the form
            return fail(500, {
                error: error instanceof Error
                    ? `API Error: ${error.message}`
                    : 'An unexpected error occurred',
                values: Object.fromEntries(formData)
            });
        }
    }
} satisfies Actions;