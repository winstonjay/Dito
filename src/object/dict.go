package object

import (
	"bytes"
	"fmt"
	"strings"
)

// Dict :
type Dict struct {
	Items map[HashKey]DictItem
	Len   int
}

// DictItem :
type DictItem struct {
	Key   Object
	Value Object
}

// Type : return objects type as a TypeFlag
func (d *Dict) Type() TypeFlag { return DictType }

// Inspect : return a string representation of the objects value.
func (d *Dict) Inspect() string {
	var out bytes.Buffer
	items := []string{}
	for _, item := range d.Items {
		if item.Key.Type() == StringType {
			items = append(items, fmt.Sprintf("\"%s\": %s",
				item.Key.Inspect(), item.Value.Inspect()))
		} else {
			items = append(items, fmt.Sprintf("%s: %s",
				item.Key.Inspect(), item.Value.Inspect()))
		}
	}
	out.WriteString("{")
	out.WriteString(strings.Join(items, ", "))
	out.WriteString("}")
	return out.String()
}

// NewDict : return new initialized instance of the object.
func NewDict() *Dict {
	return &Dict{
		Items: make(map[HashKey]DictItem),
	}
}

// ConvertType : return the conversion into the specified type
func (d *Dict) ConvertType(which TypeFlag) Object {
	if which == DictType {
		return d
	}
	return NewError("Argument to %s not supported, got %s", which, d.Type())
}

// ''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''
// Methods needed to satisfy the Iterable interface:

// Length : return the number of items in the Array.
func (d *Dict) Length() Object { return &Int{Value: d.Len} }

// GetItem : return the item at the position provided by the hashkey
func (d *Dict) GetItem(key Object) Object {
	hashKey, ok := key.(Hashable)
	if !ok {
		return NewError("type '%s' is not hashable", key.Type())
	}
	if item, ok := d.Items[hashKey.Hash()]; ok {
		return item.Value
	}
	return NewError("key '%s' doesn't exist", key.Inspect())
}

// SetItem : set item at index except char or a string.
func (d *Dict) SetItem(key Object, val Object) Object {
	hashKey, ok := key.(Hashable)
	if !ok {
		return NewError("type '%s' is not hashable", key.Type())
	}
	if item, ok := d.Items[hashKey.Hash()]; ok {
		item.Value = val
	} else {
		d.Items[hashKey.Hash()] = DictItem{Key: key, Value: val}
		d.Len++
	}
	return d
}

// Slice : return a slice of an arrays elements.
func (d *Dict) Slice(start Object, end Object) Object {
	return NewError("TypeError: cannot make slice of Dict")
}

// Concat : Add item to the current string creating a new string.
func (d *Dict) Concat(other Object) Object {
	otherDict, ok := other.(*Dict)
	if !ok {
		return NewError("TypeError: cannot conncat type %s to Dict", other.Type())
	}
	for _, item := range otherDict.Items {
		d.SetItem(item.Key, item.Value)
	}
	return d
}

// Contains :
func (d *Dict) Contains(key Object) Object {
	hashKey, ok := key.(Hashable)
	if !ok {
		return NewError("type '%s' is not hashable", key.Type())
	}
	if _, ok := d.Items[hashKey.Hash()]; ok {
		return TRUE
	}
	return FALSE
}

// Iter : loop through items elements in order.
func (d *Dict) Iter() <-chan Object {
	ch := make(chan Object)
	go func() {
		for _, item := range d.Items {
			ch <- item.Key
		}
		close(ch)
	}()
	return ch
}
