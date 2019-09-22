package object

import (
	"fmt"
	"hash/fnv"
)

// Char :
type Char struct{ Value byte }

// Type : return objects type as a TypeFlag
func (c *Char) Type() TypeFlag { return CharType }

// Inspect : return a string representation of the objects value.
func (c *Char) Inspect() string { return fmt.Sprintf("'%c'", c.Value) }

// NewChar : return new initialized instance of the object.
func NewChar(value byte) *Char { return &Char{Value: value} }

// ConvertType : return the conversion into the specified type
func (c *Char) ConvertType(which TypeFlag) Object {
	switch which {
	case CharType:
		return c
	case IntType:
		return NewInt(int(c.Value))
	case FloatType:
		return NewFloat(float64(c.Value))
	case StringType:
		return NewString(c.Inspect())
	case BoolType:
		return NewBool(c.Value != 0)
	default:
		return NewError("Argument to %s not supported, got %s", which, c.Type())
	}
}

// ''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''
// Methods needed to satisfy the Numeric interface:

// Abs : return the absolute value of an number
func (c *Char) Abs() *Char { return c }

// ''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''
// Methods needed to satisfy the Hashable interface:

// Hash : hash value of char
func (c *Char) Hash() HashKey {
	h := fnv.New64a()
	h.Write([]byte{c.Value})
	return HashKey{Type: c.Type(), Value: h.Sum64()}
}
