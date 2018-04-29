package object

import (
	"fmt"
	"strconv"
)

// String : builtin string type
type String struct {
	Value string
}

// Type : return objects type as a TypeFlag
func (s *String) Type() TypeFlag { return StringType }

// Inspect : return a string representation of the objects value.
func (s *String) Inspect() string { return s.Value }

// NewString : return new initialised instance of the object.
func NewString(value string) *String { return &String{Value: value} }

// ConvertType : return the conversion into the specified type
func (s *String) ConvertType(which TypeFlag) Object {
	switch which {
	case FloatType:
		if f, err := strconv.ParseFloat(s.Value, 64); err == nil {
			return NewFloat(f)
		}
	case IntType:
		if i, err := strconv.Atoi(s.Value); err == nil {
			return NewInt(i)
		}
	case StringType:
		return s
	case BoolType:
		return NewBool(s.Value != "")
	}
	return NewError("Argument to %s not supported, got %s", s.Type(), which)
}

// Methods needed to satisfy the Iterable interface:

// Length : return the number of items in the String. An item is an ascii char
// like in C strings.
func (s *String) Length() Object {
	return &Int{Value: len(s.Value)}
}

// GetItem : return the char at index
func (s *String) GetItem(key Object) Object {
	if index, ok := key.(*Int); ok {
		size := len(s.Value)
		if index.Value >= 0 && index.Value < size {
			return &Char{Value: byte(s.Value[index.Value])}
		}
	}
	// TODO : give better error.
	return NewError("index error")
}

// SetItem : set item at index except char or a string.
func (s *String) SetItem(key Object, val Object) Object {
	index, ok := key.(*Int)
	if !ok {
		return NewError("index error")
	}
	size := len(s.Value)
	if index.Value < 0 || index.Value >= size {
		return NewError("index error")
	}
	tmp := []byte(s.Value)
	switch v := val.(type) {
	case *Char:
		tmp[index.Value] = v.Value
		s.Value = string(tmp)
		return nil
	case *String:
		if len(v.Value) > 1 {
			return NewError("Index Assignment Error. Length of value too long")
		}
		tmp[index.Value] = v.Value[0]
		s.Value = string(tmp)
		fmt.Printf(s.Value)
		return nil
	default:
		return NewError("Index Assignment Error")
	}
}
