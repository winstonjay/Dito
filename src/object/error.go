package object

import (
	"fmt"
)

// Error : builtin Error type.
type Error struct{ Message string }

// Type : return objects type as a TypeFlag
func (e *Error) Type() TypeFlag { return ErrorType }

// Inspect : return a string representation of the objects value.
func (e *Error) Inspect() string { return ErrorType.String() + ", " + e.Message }

// NewError : return new initialized instance of the object.
func NewError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

// ConvertType : return the conversion into the specified type
func (e *Error) ConvertType(which TypeFlag) Object {
	return NewError(ConvertTypeError, e.Type(), which)
}

// some pre defined errors for consistency.
const (
	InvalidArgLenError = "Wrong number of args to function %s. Want=%d. Got=%d."
	ConvertTypeError   = "Cannot convet type '%s' to type '%s'"
)
