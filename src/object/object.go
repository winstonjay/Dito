package object

import (
	"bytes"
	"dito/src/ast"
	"fmt"
	"strconv"
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
	LambdaObj   = "Lambda"
	BultinObj   = "Builtin" // not implemented
)

// singleton objects.
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
	Primitive Type Objects.
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

//

// String : builtin integer type.
type DitoString struct {
	Value string
}

func (s *DitoString) Type() string    { return StringObj }
func (s *DitoString) Inspect() string { return s.Value }

//

// Integer : builtin integer type.
// -9223372036854775807 and 9223372036854775807
type Integer struct{ Value int64 }

// Type :
func (i *Integer) Type() string    { return IntergerObj }
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

func NewDitoInteger(value int64) *Integer { return &Integer{Value: value} }

//

// Float : builtin float type.
type Float struct{ Value float64 }

// Type :
func (f *Float) Type() string    { return FloatObj }
func (f *Float) Inspect() string { return strconv.FormatFloat(f.Value, 'f', -1, 64) }

func NewDitoFloat(value float64) *Float { return &Float{Value: value} }

//

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

//

// None : builtin None type.
type None struct{}

func (n *None) Type() string    { return NoneObj }
func (n *None) Inspect() string { return NoneObj }

//

// Error : builtin Error type.
type Error struct{ Message string }

func (e *Error) Type() string    { return ErrorObj }
func (e *Error) Inspect() string { return ErrorObj + ": " + e.Message }

//

// LambdaFn :
type LambdaFn struct {
	Parameters []*ast.Identifier
	Expr       ast.Expression
	Env        *Environment
}

func (lf *LambdaFn) Type() string    { return LambdaObj }
func (lf *LambdaFn) Inspect() string { return "<Lambda function>" }

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
