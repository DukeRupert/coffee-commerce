import { createCoffeeServiceClient, type CoffeeServiceClient } from './coffeServiceClient/main';
import { PUBLIC_API_URL, PUBLIC_API_VERSION, PUBLIC_API_TIMEOUT } from '$env/static/public';
/**
 * Parse and validate a timeout value from environment variables
 * 
 * @param timeoutStr The timeout string to parse
 * @param defaultValue The default timeout value if parsing fails
 * @param minValue The minimum allowed timeout value
 * @param maxValue The maximum allowed timeout value
 * @returns A valid timeout value
 */
function parseTimeout(
    timeoutStr: string | undefined,
    defaultValue = 5000,
    minValue = 1000,
    maxValue = 30000
): number {
    // If no timeout string provided, return default
    if (!timeoutStr) {
        return defaultValue;
    }

    // Try to parse the timeout string
    const parsedTimeout = parseInt(timeoutStr, 10);

    // Check if parsing resulted in a valid number
    if (isNaN(parsedTimeout)) {
        console.warn(`Invalid timeout value: "${timeoutStr}". Using default: ${defaultValue}ms`);
        return defaultValue;
    }

    // Enforce minimum timeout
    if (parsedTimeout < minValue) {
        console.warn(`Timeout value ${parsedTimeout}ms is too low. Using minimum: ${minValue}ms`);
        return minValue;
    }

    // Enforce maximum timeout
    if (parsedTimeout > maxValue) {
        console.warn(`Timeout value ${parsedTimeout}ms is too high. Using maximum: ${maxValue}ms`);
        return maxValue;
    }

    // Return the valid timeout value
    return parsedTimeout;
}

// Set default API values
const DEFAULT_API_URL = 'http://localhost:8080';
const DEFAULT_API_VERSION = 'v1';
const DEFAULT_TIMEOUT = 5000;

// Parse and validate the timeout
const timeout = parseTimeout(
    PUBLIC_API_TIMEOUT,
    DEFAULT_TIMEOUT,
    1000,    // Minimum 1 second
    30000    // Maximum 30 seconds
);

// Create the API client
const client: CoffeeServiceClient = createCoffeeServiceClient({
    apiUrl: PUBLIC_API_URL || DEFAULT_API_URL,
    apiVersion: PUBLIC_API_VERSION || DEFAULT_API_VERSION,
    timeout: timeout
});

export default client;