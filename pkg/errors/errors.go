package errors

import "fmt"

type NotFound struct {
	Object  string
	Message string
}

func (e *NotFound) Error() string {
	return fmt.Sprintf("%s not found", e.Object)
}

// creates a not found custom error
func NewNotFound(object, msg string) *NotFound {
	return &NotFound{
		object,
		msg,
	}
}
