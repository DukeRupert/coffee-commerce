package dto

// SyncResult represents the result of syncing a single product
type SyncResult struct {
	ProductID      string `json:"product_id"`
	ProductName    string `json:"product_name"`
	StoredStripeID string `json:"stored_stripe_id"`
	ActualStripeID string `json:"actual_stripe_id,omitempty"`
	Status         string `json:"status"` // "ok", "mismatch", "not_found", "error"
	Updated        bool   `json:"updated"`
	Error          string `json:"error,omitempty"`
}

// SyncStripeProductIDsResult represents the overall sync results
type SyncStripeProductIDsResult struct {
	TotalProducts int          `json:"total_products"`
	Results       []SyncResult `json:"results"`
	Summary       struct {
		OK         int `json:"ok"`
		Mismatches int `json:"mismatches"`
		NotFound   int `json:"not_found"`
		Errors     int `json:"errors"`
		Updated    int `json:"updated"`
	} `json:"summary"`
}
