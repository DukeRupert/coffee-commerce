/**
 * Create product request data
 */
export interface CreateProductRequest {
  /**
   * Product name
   */
  name: string;
  
  /**
   * Product description
   */
  description: string;
  
  /**
   * URL to product image
   */
  image_url?: string;
  
  /**
   * Coffee origin (e.g., "Ethiopia, Yirgacheffe")
   */
  origin?: string;
  
  /**
   * Roast level (e.g., "Light", "Medium", "Dark")
   */
  roast_level?: string;
  
  /**
   * Flavor profile (e.g., "Blueberry, Chocolate, Citrus")
   */
  flavor_notes?: string;
  
  /**
   * Current stock level
   */
  stock_level?: number;
  
  /**
   * Whether the product is active and visible to customers
   */
  active?: boolean;
  
  /**
   * Whether this product can be purchased as a subscription
   */
  allow_subscription?: boolean;
  
  /**
   * Product options like available weights and grinds
   */
  options?: {
    /**
     * Available weight options (e.g., ["12oz", "1lb", "5lb"])
     */
    weights?: string[];
    
    /**
     * Available grind options (e.g., ["Whole Bean", "Drip", "Espresso"])
     */
    grinds?: string[];
  };
}

/**
 * Create product request data
 */
export interface CreateProductRequest {
  /**
   * Product name
   */
  name: string;
  
  /**
   * Product description
   */
  description: string;
  
  /**
   * URL to product image
   */
  image_url?: string;
  
  /**
   * Coffee origin (e.g., "Ethiopia, Yirgacheffe")
   */
  origin?: string;
  
  /**
   * Roast level (e.g., "Light", "Medium", "Dark")
   */
  roast_level?: string;
  
  /**
   * Flavor profile (e.g., "Blueberry, Chocolate, Citrus")
   */
  flavor_notes?: string;
  
  /**
   * Current stock level
   */
  stock_level?: number;
  
  /**
   * Whether the product is active and visible to customers
   */
  active?: boolean;
  
  /**
   * Whether this product can be purchased as a subscription
   */
  allow_subscription?: boolean;
  
  /**
   * Product options like available weights and grinds
   */
  options?: {
    /**
     * Available weight options (e.g., ["12oz", "1lb", "5lb"])
     */
    weights?: string[];
    
    /**
     * Available grind options (e.g., ["Whole Bean", "Drip", "Espresso"])
     */
    grinds?: string[];
  };
}

/**
 * Product variant pricing information
 */
export interface ProductPrice {
  amount: number;
  currency: string;
  type: string;
}

/**
 * Product variant interface
 */
export interface ProductVariant {
  id: string;
  active: boolean;
  options: Record<string, any>;
  price: ProductPrice;
  price_id: string;
  product_id: string;
  stock_level: number;
  stripe_price_id: string;
  created_at: string;
  updated_at: string;
}

/**
 * Enhanced product response with variants
 */
export interface ProductWithVariants {
  product: {
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
    options: Record<string, any>;
    allow_subscription: boolean;
    stripe_id: string;
    created_at: string;
    updated_at: string;
  };
  variants: ProductVariant[];
}

/**
 * Product update request data that matches the backend ProductUpdateDTO
 */
export interface ProductUpdateRequest {
  /**
   * Product name
   */
  name?: string;
  
  /**
   * Product description
   */
  description?: string;
  
  /**
   * URL to product image
   */
  image_url?: string;
  
  /**
   * Whether the product is active and visible to customers
   */
  active?: boolean;
  
  /**
   * Current stock level
   */
  stock_level?: number;
  
  /**
   * Weight in grams
   */
  weight?: number;
  
  /**
   * Coffee origin (e.g., "Ethiopia, Yirgacheffe")
   */
  origin?: string;
  
  /**
   * Roast level (e.g., "Light", "Medium", "Dark")
   */
  roast_level?: string;
  
  /**
   * Flavor profile (e.g., "Blueberry, Chocolate, Citrus")
   */
  flavor_notes?: string;
  
  /**
   * Product options like available weights and grinds
   */
  options?: Record<string, string[]>;
  
  /**
   * Whether this product can be purchased as a subscription
   */
  allow_subscription?: boolean;
}