package service

import (
    "errors"
    "fmt"
)

// Standard service errors that can be used across the application
var (
    // ErrInsufficientPermissions is returned when user lacks permissions for an operation
    ErrInsufficientPermissions = errors.New("insufficient permissions")
    
    // ErrInvalidInput is returned when the input to a service is invalid
    ErrInvalidInput = errors.New("invalid input")
    
    // ErrServiceUnavailable is returned when a dependent service is unavailable
    ErrServiceUnavailable = errors.New("service unavailable")
    
    // ErrOperationFailed is returned when an operation fails for a generic reason
    ErrOperationFailed = errors.New("operation failed")
)

// PermissionError is a typed error for permission-related issues with context
type PermissionError struct {
    UserID     string
    ResourceID string
    Action     string
}

func (e *PermissionError) Error() string {
    return fmt.Sprintf("user %s does not have permission to %s resource %s", 
        e.UserID, e.Action, e.ResourceID)
}

// Is implements errors.Is interface to check if this error is of type ErrInsufficientPermissions
func (e *PermissionError) Is(target error) bool {
    return target == ErrInsufficientPermissions
}

// NewPermissionError creates a new PermissionError
func NewPermissionError(userID, resourceID, action string) error {
    return &PermissionError{
        UserID:     userID,
        ResourceID: resourceID,
        Action:     action,
    }
}