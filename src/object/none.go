package object

// None : builtin None type.
type None struct{}

// Type : return objects type as a TypeFlag
func (n *None) Type() TypeFlag { return NoneType }

// Inspect : return a string representation of the objects value.
func (n *None) Inspect() string { return NoneType.String() }

// ConvertType : return the conversion into the specified type
func (n *None) ConvertType(which TypeFlag) Object {
	return NewError("Argument to %s not supported, got %s", n.Type(), which)
}
