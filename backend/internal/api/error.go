package api

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
    Status          int                    `json:"-"`             // HTTP status code (not serialized)
    Message         string                 `json:"message"`       // Human-readable error message
    ValidationErrors map[string]string     `json:"validationErrors,omitempty"` // Field validation errors
    Code            string                 `json:"code,omitempty"` // Error code for client handling
    Details         interface{}            `json:"details,omitempty"` // Additional error details
}

func SendErrorResponse(w http.ResponseWriter, err ErrorResponse) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(err.Status)
    json.NewEncoder(w).Encode(err)
}