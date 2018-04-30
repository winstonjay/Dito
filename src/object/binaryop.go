package object

import (
	"math"
)

type binaryFn func(*Environment, Object, Object) Object

type binaryOp struct {
	name string
	fn   map[TypeFlag]binaryFn
}

func (op *binaryOp) whichType(t1, t2 TypeFlag) TypeFlag {
	if t1 > 3 || t2 > 3 {
		return ErrorType
	}
	if t1 > t2 {
		return t1
	}
	return t2
}

func (op *binaryOp) EvalBinary(env *Environment, a, b Object) Object {
	which := op.whichType(a.Type(), b.Type())
	if which == ErrorType {
		return NewError("mis matched types: %s, %s", a.Type(), b.Type())
	}
	a = a.ConvertType(which)
	b = b.ConvertType(which)
	if a.Type() == ErrorType {
		return NewError("a: cannot convert types: %s, %s (%s)", a.Type(), which, a.Inspect())
	}
	if b.Type() == ErrorType {
		return NewError("b: cannot convert types: %s, %s (%s)", b.Type(), which, b.Inspect())
	}
	fn := op.fn[which]
	if fn == nil {
		return NewError("unknown binary function")
	}
	return fn(env, a, b)
}

// BinaryOps : binary operations availible to the user.
var BinaryOps = make(map[string]*binaryOp)

func init() {

	ops := []*binaryOp{
		{
			name: "+",
			fn: map[TypeFlag]binaryFn{
				CharType: func(env *Environment, a, b Object) Object {
					return NewChar(a.(*Char).Value + b.(*Char).Value)
				},
				IntType: func(env *Environment, a, b Object) Object {
					return NewInt(a.(*Int).Value + b.(*Int).Value)
				},
				FloatType: func(env *Environment, a, b Object) Object {
					return NewFloat(a.(*Float).Value + b.(*Float).Value)
				},
			},
		},

		{
			name: "-",
			fn: map[TypeFlag]binaryFn{
				CharType: func(env *Environment, a, b Object) Object {
					return NewChar(a.(*Char).Value - b.(*Char).Value)
				},
				IntType: func(env *Environment, a, b Object) Object {
					return NewInt(a.(*Int).Value - b.(*Int).Value)
				},
				FloatType: func(env *Environment, a, b Object) Object {
					return NewFloat(a.(*Float).Value - b.(*Float).Value)
				},
			},
		},

		{
			name: "*",
			fn: map[TypeFlag]binaryFn{
				CharType: func(env *Environment, a, b Object) Object {
					return NewChar(a.(*Char).Value * b.(*Char).Value)
				},
				IntType: func(env *Environment, a, b Object) Object {
					return NewInt(a.(*Int).Value * b.(*Int).Value)
				},
				FloatType: func(env *Environment, a, b Object) Object {
					return NewFloat(a.(*Float).Value * b.(*Float).Value)
				},
			},
		},

		{
			name: "/",
			fn: map[TypeFlag]binaryFn{
				CharType: func(env *Environment, a, b Object) Object {
					return NewFloat(float64(a.(*Char).Value) / float64(b.(*Char).Value))
				},
				IntType: func(env *Environment, a, b Object) Object {
					return NewFloat(float64(a.(*Int).Value) / float64(b.(*Int).Value))
				},
				FloatType: func(env *Environment, a, b Object) Object {
					return NewFloat(a.(*Float).Value / b.(*Float).Value)
				},
			},
		},

		{
			name: "%",
			fn: map[TypeFlag]binaryFn{
				CharType: func(env *Environment, a, b Object) Object {
					return NewChar(a.(*Char).Value % b.(*Char).Value)
				},
				IntType: func(env *Environment, a, b Object) Object {
					return NewInt(a.(*Int).Value % b.(*Int).Value)
				},
			},
		},

		{
			name: "//",
			fn: map[TypeFlag]binaryFn{
				CharType: func(env *Environment, a, b Object) Object {
					return NewChar(a.(*Char).Value / b.(*Char).Value)
				},
				IntType: func(env *Environment, a, b Object) Object {
					return NewInt(a.(*Int).Value / b.(*Int).Value)
				},
			},
		},

		{
			name: "**",
			fn: map[TypeFlag]binaryFn{
				CharType: func(env *Environment, a, b Object) Object {
					v := math.Pow(float64(a.(*Int).Value), float64(b.(*Int).Value))
					return NewInt(int(v))
				},
				IntType: func(env *Environment, a, b Object) Object {
					v := math.Pow(float64(a.(*Int).Value), float64(b.(*Int).Value))
					return NewInt(int(v))
				},
				FloatType: func(env *Environment, a, b Object) Object {
					v := math.Pow(a.(*Float).Value, b.(*Float).Value)
					return NewFloat(v)
				},
			},
		},

		{
			name: "<",
			fn: map[TypeFlag]binaryFn{
				CharType: func(env *Environment, a, b Object) Object {
					return NewBool(a.(*Char).Value < b.(*Char).Value)
				},
				IntType: func(env *Environment, a, b Object) Object {
					return NewBool(a.(*Int).Value < b.(*Int).Value)
				},
				FloatType: func(env *Environment, a, b Object) Object {
					return NewBool(a.(*Float).Value < b.(*Float).Value)
				},
			},
		},

		{
			name: ">",
			fn: map[TypeFlag]binaryFn{
				CharType: func(env *Environment, a, b Object) Object {
					return NewBool(a.(*Char).Value > b.(*Char).Value)
				},
				IntType: func(env *Environment, a, b Object) Object {
					return NewBool(a.(*Int).Value > b.(*Int).Value)
				},
				FloatType: func(env *Environment, a, b Object) Object {
					return NewBool(a.(*Float).Value > b.(*Float).Value)
				},
			},
		},

		{
			name: "<=",
			fn: map[TypeFlag]binaryFn{
				CharType: func(env *Environment, a, b Object) Object {
					return NewBool(a.(*Char).Value <= b.(*Char).Value)
				},
				IntType: func(env *Environment, a, b Object) Object {
					return NewBool(a.(*Int).Value <= b.(*Int).Value)
				},
				FloatType: func(env *Environment, a, b Object) Object {
					return NewBool(a.(*Float).Value <= b.(*Float).Value)
				},
			},
		},

		{
			name: ">=",
			fn: map[TypeFlag]binaryFn{
				CharType: func(env *Environment, a, b Object) Object {
					return NewBool(a.(*Char).Value >= b.(*Char).Value)
				},
				IntType: func(env *Environment, a, b Object) Object {
					return NewBool(a.(*Int).Value >= b.(*Int).Value)
				},
				FloatType: func(env *Environment, a, b Object) Object {
					return NewBool(a.(*Float).Value >= b.(*Float).Value)
				},
			},
		},

		{
			name: "==",
			fn: map[TypeFlag]binaryFn{
				BoolType: func(env *Environment, a, b Object) Object {
					return NewBool(a.(*Bool).Value == b.(*Bool).Value)
				},
				CharType: func(env *Environment, a, b Object) Object {
					return NewBool(a.(*Char).Value == b.(*Char).Value)
				},
				IntType: func(env *Environment, a, b Object) Object {
					return NewBool(a.(*Int).Value == b.(*Int).Value)
				},
				FloatType: func(env *Environment, a, b Object) Object {
					return NewBool(a.(*Float).Value == b.(*Float).Value)
				},
			},
		},

		{
			name: "!=",
			fn: map[TypeFlag]binaryFn{
				BoolType: func(env *Environment, a, b Object) Object {
					return NewBool(a.(*Bool).Value != b.(*Bool).Value)
				},
				CharType: func(env *Environment, a, b Object) Object {
					return NewBool(a.(*Char).Value != b.(*Char).Value)
				},
				IntType: func(env *Environment, a, b Object) Object {
					return NewBool(a.(*Int).Value != b.(*Int).Value)
				},
				FloatType: func(env *Environment, a, b Object) Object {
					return NewBool(a.(*Float).Value != b.(*Float).Value)
				},
			},
		},
	}

	for _, op := range ops {
		BinaryOps[op.name] = op
	}
}
