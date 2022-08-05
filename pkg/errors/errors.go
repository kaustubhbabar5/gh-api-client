package errors

import (
	"encoding/json"
	"fmt"
)

type JSONErrs []error

func (errs JSONErrs) MarshalJSON() ([]byte, error) {
	res := make([]string, len(errs))
	for i, e := range errs {
		res[i] = e.Error()
	}
	return json.Marshal(res)
}

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
