package object

import (
	"fmt"
)

// Int : builtin integer type.
// -9223372036854775807 and 9223372036854775807
type Int struct{ Value int }

// Type : return objects type as a TypeFlag
func (i *Int) Type() TypeFlag { return IntType }

// Inspect : return a string representation of the objects value.
func (i *Int) Inspect() string { return fmt.Sprintf("%d", i.Value) }

// NewInt : return new initialised instance of the object.
func NewInt(value int) *Int { return &Int{Value: value} }

// ConvertType : return the conversion into the specified type
func (i *Int) ConvertType(which TypeFlag) Object {
	switch which {
	case IntType:
		return i
	case CharType:
		return NewChar(byte(i.Value))
	case FloatType:
		return NewFloat(float64(i.Value))
	case StringType:
		return NewString(i.Inspect())
	case BoolType:
		return NewBool(i.Value != 0)
	case ArrayType:
		elements := make([]Object, i.Value)
		for j := 0; j < i.Value; j++ {
			elements[j] = ZERO
		}
		return NewArray(elements, i.Value)
	default:
		return NewError("Argument to %s not supported, got %s", which, i.Type())
	}
}

var (
	// ZERO : is the number zero
	ZERO = NewInt(0)
)

// ''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''
// Methods needed to satisfy the Numeric interface:

// Abs : return the absolute value of an number
func (i *Int) Abs() Object {
	if i.Value >= 0 {
		return i
	}
	return &Int{Value: -i.Value}
}
