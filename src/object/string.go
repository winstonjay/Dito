package object

import "strconv"

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
		if i, err := strconv.ParseInt(s.Value, 10, 64); err == nil {
			return NewInt(i)
		}
	case StringType:
		return s
	case BoolType:
		return NewBool(s.Value != "")
	}
	return NewError("Argument to %s not supported, got %s", s.Type(), which)
}
