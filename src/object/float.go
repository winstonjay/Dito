package object

import (
	"math"
	"strconv"
)

// Float : builtin float type basically go's float64
type Float struct{ Value float64 }

// Type : return objects type as a TypeFlag
func (f *Float) Type() TypeFlag { return FloatType }

// Inspect : return a string representation of the objects value.
func (f *Float) Inspect() string { return strconv.FormatFloat(f.Value, 'f', -1, 64) }

// NewFloat : return new initialised instance of the object.
func NewFloat(value float64) *Float { return &Float{Value: value} }

// ConvertType : return the conversion into the specified type
func (f *Float) ConvertType(which TypeFlag) Object {
	switch which {
	case FloatType:
		return f
	case IntType:
		return NewInt(int(f.Value))
	case StringType:
		return NewString(f.Inspect())
	case BoolType:
		return NewBool(f.Value != 0)
	default:
		return NewError("Argument to %s not supported, got %s", f.Type(), which)
	}
}

// ''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''
// Methods needed to satisfy the Numeric interface:

// Abs : return the absolute value of an number
func (f *Float) Abs() Object {
	return &Float{Value: math.Abs(f.Value)}
}
