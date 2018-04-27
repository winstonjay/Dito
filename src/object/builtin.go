package object

// BuiltinFunction :
type BuiltinFunction func(args ...Object) Object

// Builtin :
type Builtin struct {
	Fn         func(args ...Object) Object
	ArgsMax    int
	ArgsMin    int
	ArgType    []string
	ReturnType string
}

// Type : return objects type as a TypeFlag
func (b *Builtin) Type() TypeFlag { return BultinType }

// Inspect : return a string representation of the objects value.
func (b *Builtin) Inspect() string { return "<builtin function>" }

// NewBuiltin : return new initialised instance of the object.
func NewBuiltin(fn BuiltinFunction, argsMax, argsMin int, argType []string) *Builtin {
	return &Builtin{Fn: fn, ArgsMax: argsMax, ArgsMin: argsMin, ArgType: argType}
}

// ConvertType : return the conversion into the specified type
func (b *Builtin) ConvertType(which TypeFlag) Object {
	return NewError(ConvertTypeError, b.Type(), which)
}
