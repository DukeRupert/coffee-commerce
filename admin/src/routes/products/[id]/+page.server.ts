// admin/src/routes/products/[id]/+page.server.ts
import type { PageServerLoad } from './$types';
import client from "$lib/client";

export const load: PageServerLoad = async ({ params }) => {
    try {
        const productId = params.id;
        
        if (!productId) {
            throw new Error('Product ID is required');
        }

        // Call the API to get the product details
        const productData = await client.getProduct(productId);

        return {
            product: productData,
            productId
        };
    } catch (error) {
        console.error('Error loading product:', error);
        
        // Return error state - you might want to handle this differently
        // based on your error handling strategy
        return {
            product: null,
            productId: params.id,
            error: error instanceof Error ? error.message : 'Unknown error occurred'
        };
    }
};