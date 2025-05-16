export interface ApiErrorResponse {
  status: number;
  message: string;
  validationErrors?: Record<string, string>; // Field-specific errors
  code?: string;                            // Optional error code for programmatic handling
}

// Custom error class that includes the structured error response
export class ApiError extends Error {
  status: number;
  validationErrors?: Record<string, string>;
  code?: string;

  constructor(errorResponse: ApiErrorResponse) {
    super(errorResponse.message);
    this.name = 'ApiError';
    this.status = errorResponse.status;
    this.validationErrors = errorResponse.validationErrors;
    this.code = errorResponse.code;
  }

  hasValidationError(field: string): boolean {
    return !!(this.validationErrors && this.validationErrors[field]);
  }

  getValidationError(field: string): string | null {
    return this.validationErrors?.[field] || null;
  }
}