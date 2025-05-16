package postgres

import (
    "errors"
    "fmt"
)

// Standard repository errors that can be used across the application
var (
    // ErrDuplicateName is returned when trying to create a resource with a name that already exists
    ErrDuplicateName = errors.New("resource with this name already exists")
    
    // ErrDatabaseConnection is returned when there's an issue connecting to the database
    ErrDatabaseConnection = errors.New("database connection error")
    
    // ErrResourceNotFound is returned when a requested resource cannot be found
    ErrResourceNotFound = errors.New("resource not found")
    
    // ErrTransactionFailed is returned when a database transaction fails
    ErrTransactionFailed = errors.New("database transaction failed")
)

// DuplicateNameError is a typed error for duplicate name scenarios with additional context
type DuplicateNameError struct {
    ResourceType string
    Name         string
}

func (e *DuplicateNameError) Error() string {
    return fmt.Sprintf("%s with name '%s' already exists", e.ResourceType, e.Name)
}

// Is implements errors.Is interface to check if this error is of type ErrDuplicateName
func (e *DuplicateNameError) Is(target error) bool {
    return target == ErrDuplicateName
}

// NewDuplicateNameError creates a new DuplicateNameError
func NewDuplicateNameError(resourceType, name string) error {
    return &DuplicateNameError{
        ResourceType: resourceType,
        Name:         name,
    }
}