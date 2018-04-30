package object

import (
	"bytes"
)

// Array : array object. TODO: Need to decide if arrays will allow mixed items.
// I dont like this tbh so maybe enforce a type system.
type Array struct {
	Elements []Object
	Len      int
}

// Type : return objects type as a TypeFlag
func (a *Array) Type() TypeFlag { return ArrayType }

// Inspect : return a string representation of the objects value.
func (a *Array) Inspect() string {
	var out bytes.Buffer
	out.WriteString("[")
	for i, el := range a.Elements {
		out.WriteString(el.Inspect())
		if i < len(a.Elements)-1 {
			out.WriteString(", ")
		}
	}
	out.WriteString("]")
	return out.String()
}

// NewArray :
func NewArray(elements []Object, length int) *Array {
	if length == -1 {
		length = len(elements)
	}
	return &Array{Elements: elements, Len: length}
}

// ConvertType : return the conversion into the specified type
func (a *Array) ConvertType(which TypeFlag) Object {
	if which == ArrayType {
		return a
	}
	return NewError("Argument to %s not supported, got %s", a.Type(), which)
}

// ''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''
// Methods needed to satisfy the Iterable interface:

// Length : return the number of items in the Array.
func (a *Array) Length() Object { return &Int{Value: a.Len} }

// GetItem : return the item at the position provided by the index key
func (a *Array) GetItem(key Object) Object {
	if idx, ok := key.(*Int); ok {
		if idx.Value >= 0 && idx.Value < a.Len {
			return a.Elements[idx.Value]
		}
	}
	// TODO : give better error.
	return NewError("index error")
}

// SetItem : set item at index except char or a string.
func (a *Array) SetItem(key Object, val Object) Object {
	idx, ok := key.(*Int)
	if !ok {
		return NewError("index error")
	}
	if idx.Value < 0 || idx.Value >= a.Len {
		return NewError("index error")
	}
	a.Elements[idx.Value] = val
	return nil
}

// Iter : loop through items elements in order.
func (a *Array) Iter() <-chan Object {
	ch := make(chan Object)
	go func() {
		for _, item := range a.Elements {
			ch <- item
		}
		close(ch)
	}()
	return ch
}
