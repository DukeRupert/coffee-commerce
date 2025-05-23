// admin/src/routes/products/[id]/edit/+page.server.ts
import type { PageServerLoad, Actions } from './$types';
import client from "$lib/client";

interface ProductUpdateRequest {
  name?: string;
  description?: string;
  image_url?: string;
  active?: boolean;
  stock_level?: number;
  weight?: number;
  origin?: string;
  roast_level?: string;
  flavor_notes?: string;
  options?: Record<string, string[]>;
  allow_subscription?: boolean;
}

export const load: PageServerLoad = async ({ params }) => {
    try {
        const productId = params.id;
        
        if (!productId) {
            throw new Error('Product ID is required');
        }

        // Call the API to get the product details
        const productData = await client.getProduct(productId);

        return {
            product: productData.product,
            variants: productData.variants,
            productId
        };
    } catch (error) {
        console.error('Error loading product for edit:', error);
        
        return {
            product: null,
            variants: [],
            productId: params.id,
            error: error instanceof Error ? error.message : 'Unknown error occurred'
        };
    }
};

export const actions = {
    update: async ({ request, params }) => {
        try {
            const productId = params.id;
            
            if (!productId) {
                return {
                    success: false,
                    message: 'Product ID is required'
                };
            }

            const formData = await request.formData();
            
            // Parse form data to match ProductUpdateDTO structure
            const updateData: any = {};

            // Only include fields that were actually submitted and have values
            const name = formData.get('name') as string;
            if (name !== null && name.trim() !== '') {
                updateData.name = name.trim();
            }

            const description = formData.get('description') as string;
            if (description !== null && description.trim() !== '') {
                updateData.description = description.trim();
            }

            const imageUrl = formData.get('image_url') as string;
            if (imageUrl !== null && imageUrl.trim() !== '') {
                updateData.image_url = imageUrl.trim();
            }

            const origin = formData.get('origin') as string;
            if (origin !== null && origin.trim() !== '') {
                updateData.origin = origin.trim();
            }

            const roastLevel = formData.get('roast_level') as string;
            if (roastLevel !== null && roastLevel.trim() !== '') {
                updateData.roast_level = roastLevel.trim();
            }

            const flavorNotes = formData.get('flavor_notes') as string;
            if (flavorNotes !== null && flavorNotes.trim() !== '') {
                updateData.flavor_notes = flavorNotes.trim();
            }

            const stockLevel = formData.get('stock_level') as string;
            if (stockLevel !== null && stockLevel.trim() !== '') {
                updateData.stock_level = parseInt(stockLevel, 10);
            }

            const weight = formData.get('weight') as string;
            if (weight !== null && weight.trim() !== '') {
                updateData.weight = parseInt(weight, 10);
            }

            // Handle boolean fields
            const active = formData.get('active') === 'on';
            updateData.active = active;

            const allowSubscription = formData.get('allow_subscription') === 'on';
            updateData.allow_subscription = allowSubscription;

            // Handle options - parse from JSON string
            const optionsJson = formData.get('options') as string;
            if (optionsJson && optionsJson.trim() !== '') {
                try {
                    const options = JSON.parse(optionsJson);
                    updateData.options = options;
                } catch (e) {
                    console.error('Failed to parse options JSON:', e);
                    return {
                        success: false,
                        message: 'Invalid options format',
                        validationErrors: { options: 'Invalid JSON format' }
                    };
                }
            }

            // Call the API to update the product
            const updatedProduct = await client.updateProduct(productId, updateData);

            // Return success response instead of throwing redirect
            return {
                success: true,
                message: 'Product updated successfully',
                redirect: `/products/${productId}`
            };

        } catch (error) {
            console.error('Error updating product:', error);
            
            // Handle redirect as success (when it comes from our own redirect)
            if (error && typeof error === 'object' && 'status' in error && error.status === 303) {
                return {
                    success: true,
                    message: 'Product updated successfully',
                    redirect: error.location
                };
            }
            
            // Handle validation errors from API
            if (error && typeof error === 'object' && 'validationErrors' in error) {
                return {
                    success: false,
                    message: error.message || 'Validation failed',
                    validationErrors: error.validationErrors
                };
            }

            return {
                success: false,
                message: error instanceof Error ? error.message : 'An unknown error occurred'
            };
        }
    }
} satisfies Actions;