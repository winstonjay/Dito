package object

import (
	"fmt"
	"strings"
)

// BuiltinFunction :
type BuiltinFunction func(args ...Object) Object

// Builtin :
type Builtin struct {
	Fn      BuiltinFunction
	Name    string
	Info    string
	ArgC    int
	ArgT    []string
	ReturnT string
}

// Type : return objects type as a TypeFlag
func (b *Builtin) Type() TypeFlag { return BultinType }

// Inspect : return a string representation of the objects value.
func (b *Builtin) Inspect() string {
	return fmt.Sprintf("%s(%s) : %s",
		b.Name, strings.Join(b.ArgT, ", "), b.ReturnT)
}

// NewBuiltin : return new initialized instance of the object.
// func NewBuiltin(fn BuiltinFunction) *Builtin {
// 	return &Builtin{Fn: fn, argC: }
// }

// ConvertType : return the conversion into the specified type
func (b *Builtin) ConvertType(which TypeFlag) Object {
	return NewError(ConvertTypeError, b.Type(), which)
}

// Doc documentation
func (b *Builtin) Doc() map[string]string {
	return map[string]string{
		"name": b.Name,
		"spec": b.Inspect(),
		"info": b.Info,
	}
}
