package object

import (
	"bytes"
	"dito/src/ast"
)

// Lambda :
type Lambda struct {
	Parameters []*ast.Identifier
	Expr       ast.Expression
	Env        *Environment
}

// Type : return objects type as a TypeFlag
func (l *Lambda) Type() TypeFlag { return LambdaType }

// Inspect : return a string representation of the objects value.
func (l *Lambda) Inspect() string {
	var out bytes.Buffer
	out.WriteString("<")
	// out.WriteString(LambdaType.String())
	out.WriteString("Lambda(")
	for i, param := range l.Parameters {
		out.WriteString(param.String())
		if i < len(l.Parameters)-1 {
			out.WriteString(", ")
		}
	}
	out.WriteString(")")
	out.WriteString(l.Expr.String())
	out.WriteString(">")
	return out.String()
}

// NewLambda : return new initialised instance of the object.
func NewLambda(params []*ast.Identifier, expr ast.Expression, env *Environment) *Lambda {
	return &Lambda{Parameters: params, Env: env, Expr: expr}
}

// ConvertType : return the conversion into the specified type
func (l *Lambda) ConvertType(which TypeFlag) Object {
	return NewError("Argument to %s not supported, got %s", l.Type(), which)
}
