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