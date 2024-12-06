package api

import (
	"fmt"
	"reflect"
)

// MissingFieldError is returned when a required field in a map is not found.
type MissingFieldError struct {
	FieldName string
}

func (err MissingFieldError) Error() string {
	return fmt.Sprintf("unable to find field '%q' in map", err.FieldName)
}

// TypeAssertionError is returned when a required type assertion fails.
type TypeAssertionError struct {
	Variable     any
	ExpectedType string
}

func (err TypeAssertionError) Error() string {
	return fmt.Sprintf("unable to convert %v to %s", err.Variable, err.ExpectedType)
}

// InterfaceTypeError is returned when a provided type for a generic function does not match the expected kind.
type InterfaceTypeError struct {
	Want     reflect.Kind
	Received reflect.Type
}

func (err InterfaceTypeError) Error() string {
	return fmt.Sprintf(
		"interface type must be a '%s', received '%s'", err.Want.String(), err.Received.String())
}
