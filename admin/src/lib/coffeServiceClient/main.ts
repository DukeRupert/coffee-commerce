// CoffeeServiceClient.ts
import type { CreateProductRequest, ProductWithVariants, ProductUpdateRequest } from './interface';
import { ApiError, type ApiErrorResponse } from './error';

/**
 * Response structure for paginated API responses
 */
export interface PaginatedResponse<T> {
	data: T[];
	meta: Meta;
}

/**
 * Metadata for paginated responses
 */
export interface Meta {
	page: number;
	per_page: number;
	total: number;
	total_pages: number;
	has_next: boolean;
	has_prev: boolean;
}

/**
 * Product interface based on the backend model
 */
export interface Product {
	id: string;
	name: string;
	description: string;
	image_url: string;
	active: boolean;
	stock_level: number;
	weight: number;
	origin: string;
	roast_level: string;
	flavor_notes: string;
	options: Record<string, string[]>;
	allow_subscription: boolean;
	stripe_id: string;
	created_at: string;
	updated_at: string;
}

/**
 * Options for configuring the CoffeeServiceClient
 */
export interface CoffeeServiceClientOptions {
	apiUrl: string;
	/**
	 * Default request timeout in milliseconds
	 * @default 10000 (10 seconds)
	 */
	timeout?: number;
	/**
	 * API version path segment
	 * @default 'v1'
	 */
	apiVersion?: string;
	/**
	 * Authorization token for authenticated requests
	 */
	authToken?: string /**
	 * Custom fetch implementation (for proxying, testing, etc.)
	 */;
	customFetch?: typeof fetch;
}

/**
 * Parameters for the getProducts method
 */
export interface GetProductsParams {
	/**
	 * Page number for pagination
	 * @default 1
	 */
	page?: number;
	/**
	 * Number of items per page
	 * @default 20
	 */
	per_page?: number;
	/**
	 * Filter products by active status
	 */
	active?: boolean;
	/**
	 * Filter products by subscription capability
	 */
	allow_subscription?: boolean;
	/**
	 * Search term for product name and description
	 */
	search?: string;
	/**
	 * Sort field
	 * @example 'name', 'created_at', etc.
	 */
	sort_by?: string;
	/**
	 * Sort direction
	 * @example 'asc' or 'desc'
	 */
	sort_dir?: 'asc' | 'desc';
}

/**
 * Client for interacting with the Coffee Service API
 */
export class CoffeeServiceClient {
	private apiUrl: string;
	private timeout: number;
	private apiVersion: string;
	private authToken?: string;
	private fetchImpl: typeof fetch;

	/**
	 * Create a new instance of the Coffee Service API client
	 * @param options Client configuration options
	 */
	constructor(options: CoffeeServiceClientOptions) {
		this.apiUrl = options.apiUrl.endsWith('/') ? options.apiUrl.slice(0, -1) : options.apiUrl;
		this.timeout = options.timeout || 10000;
		this.apiVersion = options.apiVersion || 'v1';
		this.authToken = options.authToken;
		this.fetchImpl = options.customFetch || fetch;
	}

	/**
	 * Set the authorization token for authenticated requests
	 * @param token JWT token or other authorization token
	 */
	public setAuthToken(token: string): void {
		this.authToken = token;
	}

	/**
	 * Clear the authorization token
	 */
	public clearAuthToken(): void {
		this.authToken = undefined;
	}

	/**
	 * Get the base URL for API requests
	 * @returns The base URL with API version
	 */
	private getBaseUrl(): string {
		return `${this.apiUrl}/api/${this.apiVersion}`;
	}

	/**
	 * Get common headers for all requests
	 * @returns HTTP headers object
	 */
	private getHeaders(): HeadersInit {
		const headers: HeadersInit = {
			'Content-Type': 'application/json',
			Accept: 'application/json'
		};

		if (this.authToken) {
			headers['Authorization'] = `Bearer ${this.authToken}`;
		}

		return headers;
	}

	/**
	 * Helper method to handle API responses
	 * @param response Fetch Response object
	 * @returns Parsed response data
	 * @throws Error if the response is not successful
	 */
	public async handleResponse<T>(response: Response): Promise<T> {
		const contentType = response.headers.get('content-type');
		const isJson = contentType && contentType.includes('application/json');

		if (!response.ok) {
			let errorResponse: ApiErrorResponse;

			if (isJson) {
				try {
					const errorData = await response.json();
					errorResponse = {
						status: response.status,
						message: errorData.message || response.statusText,
						validationErrors: errorData.validationErrors,
						code: errorData.code
					};
				} catch (e) {
					errorResponse = {
						status: response.status,
						message: `Request failed with status ${response.status}`
					};
				}
			} else {
				const text = await response.text();
				errorResponse = {
					status: response.status,
					message: text || response.statusText
				};
			}

			throw new ApiError(errorResponse);
		}

		// Handle successful response
		if (isJson) {
			return (await response.json()) as T;
		} else if (response.status === 204) {
			return {} as T; // No content
		} else {
			const text = await response.text();
			return text as unknown as T;
		}
	}

	/**
	 * Create a URL with query parameters
	 * @param endpoint API endpoint
	 * @param params Query parameters
	 * @returns URL with query parameters
	 */
	private createUrl(endpoint: string, params?: Record<string, any>): string {
		const url = new URL(`${this.getBaseUrl()}${endpoint}`);

		if (params) {
			Object.entries(params).forEach(([key, value]) => {
				if (value !== undefined && value !== null) {
					url.searchParams.append(key, String(value));
				}
			});
		}

		return url.toString();
	}

	/**
	 * Get a list of products with optional filtering and pagination
	 * @param params Query parameters for filtering and pagination
	 * @returns Promise with paginated product data
	 */
	public async getProducts(params?: GetProductsParams): Promise<PaginatedResponse<Product>> {
		const url = this.createUrl('/products', params);

		try {
			const response = await this.fetchImpl(url, {
				method: 'GET',
				headers: this.getHeaders(),
				signal: AbortSignal.timeout(this.timeout)
			});

			return this.handleResponse<PaginatedResponse<Product>>(response);
		} catch (error) {
			if (error instanceof Error) {
				throw new Error(`Failed to fetch products: ${error.message}`);
			}
			throw error;
		}
	}

	/**
	 * Create a new product
	 * @param productData Data for the new product
	 * @returns Promise with the created product data
	 */
	public async createProduct(productData: CreateProductRequest): Promise<Product> {
		const url = this.createUrl('/products');

		try {
			const response = await this.fetchImpl(url, {
				method: 'POST',
				headers: this.getHeaders(),
				body: JSON.stringify(productData),
				signal: AbortSignal.timeout(this.timeout)
			});

			return this.handleResponse<Product>(response);
		} catch (error) {
			if (error instanceof ApiError) {
				// Already formatted, just rethrow
				throw error;
			}
			if (error instanceof Error) {
				throw new ApiError({
					status: 500,
					message: `Failed to create product: ${error.message}`
				});
			}
			throw new ApiError({
				status: 500,
				message: 'An unknown error occurred'
			});
		}
	}

	/**
	 * Delete a product
	 * @param productId ID of the product to delete
	 * @returns Promise with the deletion status
	 */
	public async deleteProduct(productId: string): Promise<void> {
		const url = this.createUrl(`/products/${productId}`);

		try {
			const response = await this.fetchImpl(url, {
				method: 'DELETE',
				headers: this.getHeaders(),
				signal: AbortSignal.timeout(this.timeout)
			});

			// For DELETE operations, often there's no response body
			if (response.status === 204) {
				return;
			}

			return this.handleResponse<void>(response);
		} catch (error) {
			if (error instanceof Error) {
				throw new Error(`Failed to delete product: ${error.message}`);
			}
			throw error;
		}
	}

	/**
	 * Update a product's stock level
	 * @param productId ID of the product to update
	 * @param stockLevel New stock level
	 * @returns Promise with the updated product data
	 */
	public async updateProductStock(productId: string, stockLevel: number): Promise<Product> {
		const url = this.createUrl(`/products/${productId}/stock`);

		try {
			const response = await this.fetchImpl(url, {
				method: 'PATCH',
				headers: this.getHeaders(),
				body: JSON.stringify({ stock_level: stockLevel }),
				signal: AbortSignal.timeout(this.timeout)
			});

			return this.handleResponse<Product>(response);
		} catch (error) {
			if (error instanceof Error) {
				throw new Error(`Failed to update product stock: ${error.message}`);
			}
			throw error;
		}
	}

	// Additional methods can be added here for other API endpoints
	/**
	 * Helper function to extract product filter parameters from URL parameters
	 * Designed for use with SvelteKit or similar frameworks
	 *
	 * @param params URL parameters object (from route or query params)
	 * @returns Properly formatted parameters for getProducts()
	 */
	public extractProductParams(params: Record<string, string>): GetProductsParams {
		return {
			// Pagination
			page: params.page ? parseInt(params.page, 10) : undefined,
			per_page: params.per_page ? parseInt(params.per_page, 10) : undefined,

			// Filters
			active: params.active !== undefined ? params.active === 'true' : undefined,
			allow_subscription:
				params.allow_subscription !== undefined ? params.allow_subscription === 'true' : undefined,
			search: params.search || undefined,

			// Sorting
			sort_by: params.sort_by || undefined,
			sort_dir: params.sort_dir as ('asc' | 'desc') | undefined
		};
	}

	/**
	 * Get a single product by ID with variants
	 * @param productId ID of the product to retrieve
	 * @returns Promise with the product data including variants
	 */
	public async getProduct(productId: string): Promise<ProductWithVariants> {
		const url = this.createUrl(`/products/${productId}`);

		try {
			const response = await this.fetchImpl(url, {
				method: 'GET',
				headers: this.getHeaders(),
				signal: AbortSignal.timeout(this.timeout)
			});

			return this.handleResponse<ProductWithVariants>(response);
		} catch (error) {
			if (error instanceof Error) {
				throw new Error(`Failed to fetch product: ${error.message}`);
			}
			throw error;
		}
	}

	/**
	 * Update an existing product
	 * @param productId ID of the product to update
	 * @param productData New data for the product (only fields to update)
	 * @returns Promise with the updated product data
	 */
	public async updateProduct(
		productId: string,
		productData: ProductUpdateRequest
	): Promise<Product> {
		const url = this.createUrl(`/products/${productId}`);

		try {
			const response = await this.fetchImpl(url, {
				method: 'PUT', // Use PUT for partial updates
				headers: this.getHeaders(),
				body: JSON.stringify(productData),
				signal: AbortSignal.timeout(this.timeout)
			});

			return this.handleResponse<Product>(response);
		} catch (error) {
			if (error instanceof ApiError) {
				// Already formatted, just rethrow
				throw error;
			}
			if (error instanceof Error) {
				throw new ApiError({
					status: 500,
					message: `Failed to update product: ${error.message}`
				});
			}
			throw new ApiError({
				status: 500,
				message: 'An unknown error occurred'
			});
		}
	}
}

/**
 * Create a new Coffee Service API client
 * @param options Client configuration options
 * @returns A new CoffeeServiceClient instance
 */
export function createCoffeeServiceClient(
	options: CoffeeServiceClientOptions
): CoffeeServiceClient {
	return new CoffeeServiceClient(options);
}
