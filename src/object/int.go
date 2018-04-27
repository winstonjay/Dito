package object

import (
	"fmt"
)

// Int : builtin integer type.
// -9223372036854775807 and 9223372036854775807
type Int struct{ Value int64 }

// Type : return objects type as a TypeFlag
func (i *Int) Type() TypeFlag { return IntType }

// Inspect : return a string representation of the objects value.
func (i *Int) Inspect() string { return fmt.Sprintf("%d", i.Value) }

// NewInt : return new initialised instance of the object.
func NewInt(value int64) *Int { return &Int{Value: value} }

// ConvertType : return the conversion into the specified type
func (i *Int) ConvertType(which TypeFlag) Object {
	switch which {
	case IntType:
		return i
	case FloatType:
		return NewFloat(float64(i.Value))
	case StringType:
		return NewString(i.Inspect())
	case BoolType:
		return NewBool(i.Value != 0)
	default:
		return NewError("Argument to %s not supported, got %s", i.Type(), which)
	}
}
