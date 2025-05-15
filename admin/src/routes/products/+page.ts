import client from '$lib/client';
import type { PageLoad } from './$types';

/**
 * Page load function for the products page
 * 
 * This function will:
 * 1. Extract filter parameters from the URL query string
 * 2. Call the API with those parameters
 * 3. Return the products and meta data for the page
 */
export const load: PageLoad = async ({ url }) => {
  try {
    // Get all query parameters from the URL
    const queryParams = Object.fromEntries(url.searchParams);
    
    // Convert URL parameters to properly typed filter parameters
    const filterParams = client.extractProductParams(queryParams);
    
    // Call the API with the extracted parameters
    const response = await client.getProducts(filterParams);
    
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