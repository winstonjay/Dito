package object

import (
	"dito/src/ast"
	"fmt"
)

// Function :
type Function struct {
	Parameters []*ast.Identifier
	Name       string
	Stmts      *ast.BlockStatement
	Env        *Environment
}

// Type : return objects type as a TypeFlag
func (f *Function) Type() TypeFlag { return FunctionType }

// Inspect : return a string representation of the objects value.
func (f *Function) Inspect() string { return fmt.Sprintf("<function %s>", f.Name) }

// NewFunction : ...
func NewFunction(fn *ast.Function, env *Environment) *Function {
	obj := &Function{
		Parameters: fn.Parameters,
		Name:       fn.Name.Value,
		Stmts:      fn.Body,
		Env:        env,
	}
	env.Set(obj.Name, obj, false)
	return obj
}

// ConvertType : return the conversion into the specified type
func (f *Function) ConvertType(which TypeFlag) Object {
	return NewError("Argument to %s not supported, got %s", f.Type(), which)
}
