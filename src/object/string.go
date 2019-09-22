package object

import (
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"
)

// String : builtin string type
type String struct {
	Value string
}

// Type : return objects type as a TypeFlag
func (s *String) Type() TypeFlag { return StringType }

// Inspect : return a string representation of the objects value.
func (s *String) Inspect() string { return s.Value }

// NewString : return new initialized instance of the object.
func NewString(value string) *String { return &String{Value: value} }

// ConvertType : return the conversion into the specified type
func (s *String) ConvertType(which TypeFlag) Object {
	switch which {
	case StringType:
		return s
	case FloatType:
		if f, err := strconv.ParseFloat(s.Value, 64); err == nil {
			return NewFloat(f)
		}
	case IntType:
		if i, err := strconv.Atoi(s.Value); err == nil {
			return NewInt(i)
		}
	case ArrayType:
		n := len(s.Value)
		a := &Array{Elements: make([]Object, n), Len: n}
		for i, s := range s.Value {
			a.Elements[i] = &Char{Value: byte(s)}
		}
		return a
	case BoolType:
		return NewBool(s.Value != "")
	}
	return NewError("Argument to %s not supported, got %s", which, s.Type())
}

// ''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''
// Methods needed to satisfy the Iterable interface:

// Length : return the number of items in the String. An item is an ascii char
// like in C strings.
func (s *String) Length() Object { return &Int{Value: len(s.Value)} }

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

// Slice : return a slice of an arrays elements.
func (s *String) Slice(start Object, end Object) Object {
	startInt, ok := start.(*Int)
	if !ok {
		return NewError("slice start index type error.")
	}
	endInt, ok := end.(*Int)
	if !ok {
		return NewError("slice end index type error.")
	}
	if startInt.Value > endInt.Value {
		return NewError("slice index error. start must be less than end index")
	}
	if endInt.Value < 0 || endInt.Value > len(s.Value) {
		return NewError("slice end index out of bounds error")
	}
	if startInt.Value < 0 || startInt.Value >= len(s.Value) {
		return NewError("slice start index out of bounds error")
	}
	slice := s.Value[startInt.Value:endInt.Value]
	return NewString(slice)
}

// Concat : Add item to the current string creating a new string.
func (s *String) Concat(other Object) Object {
	return NewString(s.Value + other.(*String).Value)
}

// Contains :
func (s *String) Contains(sub Object) Object {
	return NewBool(strings.Contains(s.Value, sub.(*String).Value))
}

// Iter : loop through items elements in order.
func (s *String) Iter() <-chan Object {
	ch := make(chan Object)
	go func() {
		for _, item := range s.Value {
			ch <- NewString(string(item))
		}
		close(ch)
	}()
	return ch
}

// ''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''
// Methods needed to satisfy the Hashable interface:

// Hash : hash value of string
func (s *String) Hash() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}
