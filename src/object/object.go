package object

import (
	"bytes"
	"dito/src/ast"
	"fmt"
)

// Object : defines the interface for the objects used in the dito programming language.
type Object interface {
	// Type : type which is used internaly and avalible to the user through
	// the builtin 'type' function at runtime.
	Type() string
	// Inspect : returns the value of the object.
	// Is used to display values to the user.
	Inspect() string
}

// Define the strings used availible to the user
// to describe objects.
const (
	ArrayObj    = "Array"
	IntergerObj = "Int"
	FloatObj    = "Float"
	BooleanObj  = "Bool"
	StringObj   = "String"
	NoneObj     = "None"
	ErrorObj    = "Error"
	ReturnObj   = "Return"   // not implemented
	FunctionObj = "Function" // not implemented
	LambdaObj   = "Lambda function"
	BultinObj   = "Builtin" // not implemented
)

// ######## singleton objects.
var (
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
	NONE  = &None{}
)

// ReturnValue : Packages other objects to determine
// the end objects of programs.
type ReturnValue struct{ Value Object }

func (rv *ReturnValue) Type() string    { return ReturnObj }
func (rv *ReturnValue) Inspect() string { return rv.Value.Inspect() }

/*
######## Primitive Type Objects.
*/

// Array : array object.
type Array struct {
	Elements []Object
	Len      int64
}

func (a *Array) Type() string { return ArrayObj }
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

// NewDitoArray :
func NewDitoArray(elements []Object, length int64) *Array {
	if length == -1 {
		length = int64(len(elements))
	}
	return &Array{Elements: elements, Len: length}
}

// DitoString : builtin integer type.
type DitoString struct {
	Value string
}

func (s *DitoString) Type() string    { return StringObj }
func (s *DitoString) Inspect() string { return s.Value }

// Boolean : builtin bool type.
type Boolean struct{ Value bool }

func (b *Boolean) Type() string    { return BooleanObj }
func (b *Boolean) Inspect() string { return fmt.Sprintf("%v", b.Value) }

func NewDitoBoolean(input bool) *Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

// None : builtin None type.
type None struct{}

func (n *None) Type() string    { return NoneObj }
func (n *None) Inspect() string { return NoneObj }

// Error : builtin Error type.
type Error struct{ Message string }

func (e *Error) Type() string    { return ErrorObj }
func (e *Error) Inspect() string { return "Evaluation " + ErrorObj + ": " + e.Message }

// Function :
type Function struct {
	Parameters []*ast.Identifier
	Name       string
	Stmts      *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() string    { return FunctionObj }
func (f *Function) Inspect() string { return fmt.Sprintf("<%s>", f.Name) }

func NewFunction(params []*ast.Identifier, name string, stmts *ast.BlockStatement, env *Environment) *Function {
	return &Function{
		Parameters: params,
		Name:       name,
		Stmts:      stmts,
		Env:        env,
	}
}

// LambdaFn :
type LambdaFn struct {
	Parameters []*ast.Identifier
	Expr       ast.Expression
	Env        *Environment
}

func (lf *LambdaFn) Type() string { return LambdaObj }
func (lf *LambdaFn) Inspect() string {
	var out bytes.Buffer
	out.WriteString("<")
	out.WriteString(LambdaObj)
	out.WriteString(": func(")
	for i, param := range lf.Parameters {
		out.WriteString(param.String())
		if i < len(lf.Parameters)-1 {
			out.WriteString(", ")
		}
	}
	out.WriteString(")")
	// out.WriteString(lf.Expr.String())
	out.WriteString(">")
	return out.String()
}

func NewDitoLambdaFn(params []*ast.Identifier, expr ast.Expression, env *Environment) *LambdaFn {
	return &LambdaFn{Parameters: params, Env: env, Expr: expr}
}

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

func (b *Builtin) Type() string    { return BultinObj }
func (b *Builtin) Inspect() string { return "<builtin function>" }

func NewBuiltinFn(fn BuiltinFunction, argsMax, argsMin int, argType []string) *Builtin {
	return &Builtin{Fn: fn, ArgsMax: argsMax, ArgsMin: argsMin, ArgType: argType}
}
