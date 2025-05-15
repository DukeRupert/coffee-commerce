// CoffeeServiceClient.ts

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
  authToken?: string;
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

  /**
   * Create a new instance of the Coffee Service API client
   * @param options Client configuration options
   */
  constructor(options: CoffeeServiceClientOptions) {
    this.apiUrl = options.apiUrl.endsWith('/')
      ? options.apiUrl.slice(0, -1)
      : options.apiUrl;
    this.timeout = options.timeout || 10000;
    this.apiVersion = options.apiVersion || 'v1';
    this.authToken = options.authToken;
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
      'Accept': 'application/json',
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
  private async handleResponse<T>(response: Response): Promise<T> {
    if (!response.ok) {
      const errorData = await response.json().catch(() => null);
      throw new Error(
        `API error: ${response.status} ${response.statusText}${errorData ? ` - ${JSON.stringify(errorData)}` : ''
        }`
      );
    }

    return response.json() as Promise<T>;
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
  public async getProducts(
    params?: GetProductsParams
  ): Promise<PaginatedResponse<Product>> {
    const url = this.createUrl('/products/', params);

    try {
      const response = await fetch(url, {
        method: 'GET',
        headers: this.getHeaders(),
        signal: AbortSignal.timeout(this.timeout),
      });

      return this.handleResponse<PaginatedResponse<Product>>(response);
    } catch (error) {
      if (error instanceof Error) {
        throw new Error(`Failed to fetch products: ${error.message}`);
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
      allow_subscription: params.allow_subscription !== undefined
        ? params.allow_subscription === 'true'
        : undefined,
      search: params.search || undefined,

      // Sorting
      sort_by: params.sort_by || undefined,
      sort_dir: params.sort_dir as ('asc' | 'desc') | undefined,
    };
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