// src/routes/api/products/+server.ts
import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import type { ApiError } from '$lib/coffeServiceClient/error';
import client from '$lib/client'

export const POST: RequestHandler = async ({ request }) => {
    try {
        // Parse the JSON body
        const productData = await request.json();

        // Call the coffee service API
        const product = await client.createProduct(productData);

        // Return success response
        return json({
            success: true,
            message: 'Product created successfully',
            product
        });
    } catch (error) {
        console.error('Error creating product:', error);
        const err = error as ApiError
        return json({
            success: false,
            message: err.message || 'Validation failed',
            validationErrors: err.validationErrors
        }, { status: err.status || 400 });
    }
};