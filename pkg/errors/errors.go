package errors

import "fmt"

type NotFoundError struct {
	Object  string
	Message string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found", e.Object)
}

// creates a not found custom error.
func NewNotFound(object, msg string) *NotFoundError {
	return &NotFoundError{
		object,
		msg,
	}
}
