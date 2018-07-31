package object

import "fmt"

// Bool : builtin bool type.
type Bool struct{ Value bool }

// Type : return objects type as a TypeFlag
func (b *Bool) Type() TypeFlag { return BoolType }

// Inspect : return a string representation of the objects value.
func (b *Bool) Inspect() string { return fmt.Sprintf("%v", b.Value) }

// NewBool : return a reference to either the true or false objects.
func NewBool(input bool) *Bool {
	if input {
		return TRUE
	}
	return FALSE
}

// ConvertType : return the conversion into the specified type
func (b *Bool) ConvertType(which TypeFlag) Object {
	switch which {
	case BoolType:
		return b
	case IntType:
		if b.Value {
			return NewInt(1)
		}
		return NewInt(0)
	case FloatType:
		if b.Value {
			return NewFloat(1.0)
		}
		return NewFloat(0.0)
	}
	return NewError("Argument to %s not supported, got %s", which, b.Type())
}
