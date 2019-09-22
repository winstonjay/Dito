package object

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"math"
	"strconv"
)

// Float : builtin float type basically go's float64
type Float struct{ Value float64 }

// Type : return objects type as a TypeFlag
func (f *Float) Type() TypeFlag { return FloatType }

// Inspect : return a string representation of the objects value.
func (f *Float) Inspect() string { return strconv.FormatFloat(f.Value, 'f', -1, 64) }

// NewFloat : return new initialized instance of the object.
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
		return NewError("Argument to %s not supported, got %s", which, f.Type())
	}
}

// ''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''
// Methods needed to satisfy the Numeric interface:

// Abs : return the absolute value of an number
func (f *Float) Abs() Object {
	return &Float{Value: math.Abs(f.Value)}
}

// ''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''
// Methods needed to satisfy the Hashable interface:

// Hash : hash value of float
func (f *Float) Hash() HashKey {
	h := fnv.New64a()
	h.Write(float64ToByte(f.Value))
	return HashKey{Type: f.Type(), Value: h.Sum64()}
}

func float64ToByte(f float64) []byte {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, f)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	return buf.Bytes()
}
