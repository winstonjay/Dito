package object

import (
	"dito/src/token"
	"math"
)

type binaryFn func(*Environment, Object, Object) Object

type binaryOp struct {
	token token.Token
	name  string
	fn    map[TypeFlag]binaryFn
}

func (op *binaryOp) whichType(t1, t2 TypeFlag) TypeFlag {
	if t1 == t2 {
		return t1
	}
	if t2 == DictType {
		return DictType
	}
	if t2 == ArrayType {
		return ArrayType
	}
	if t1 > BoolType || t2 > BoolType {
		return ErrorType
	}
	if t1 > t2 {
		return t1
	}
	return t2
}

func (op *binaryOp) EvalBinary(env *Environment, a, b Object) Object {

	if op.token == token.AND || op.token == token.OR {
		which := BoolType
		fn := op.fn[which]
		if fn == nil {
			return NewError("unknown binary function for given types")
		}
		a = a.ConvertType(BoolType)
		b = b.ConvertType(BoolType)
		if a.Type() == ErrorType {
			return NewError("a: cannot convert types: %s, %s (%s)",
				a.Type(), which, a.Inspect())
		}
		if b.Type() == ErrorType {
			return NewError("b: cannot convert types: %s, %s (%s)",
				b.Type(), which, b.Inspect())
		}
		return fn(env, a, b)
	}

	//TODO: this needs cleaning up

	which := op.whichType(a.Type(), b.Type())

	if which == ErrorType {
		return NewError("mis matched types: %s, %s", a.Type(), b.Type())
	}
	if which != ArrayType && which != DictType {
		a = a.ConvertType(which)
		b = b.ConvertType(which)
		if a.Type() == ErrorType {
			return NewError("a: cannot convert types: %s, %s (%s)",
				a.Type(), which, a.Inspect())
		}
		if b.Type() == ErrorType {
			return NewError("b: cannot convert types: %s, %s (%s)",
				b.Type(), which, b.Inspect())
		}
	}
	fn := op.fn[which]
	if fn == nil {
		return NewError("unknown binary function for given types: %s %s %s",
			a.Type(), op.name, b.Type())
	}
	return fn(env, a, b)
}

// BinaryOps : binary operations available to the user.
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
			name: "<<",
			fn: map[TypeFlag]binaryFn{
				CharType: func(env *Environment, a, b Object) Object {
					return NewChar(a.(*Char).Value << b.(*Char).Value)
				},
				IntType: func(env *Environment, a, b Object) Object {
					if shift := b.(*Int).Value; shift > -1 {
						return NewInt(a.(*Int).Value << uint(shift))
					}
					return NewError("shift count type must be unsigned Int")
				},
			},
		},

		{
			name: ">>",
			fn: map[TypeFlag]binaryFn{
				CharType: func(env *Environment, a, b Object) Object {
					return NewChar(a.(*Char).Value >> b.(*Char).Value)
				},
				IntType: func(env *Environment, a, b Object) Object {
					if shift := b.(*Int).Value; shift > -1 {
						return NewInt(a.(*Int).Value >> uint(shift))
					}
					return NewError("shift count type must be unsigned Int")
				},
			},
		},

		{
			name: "&",
			fn: map[TypeFlag]binaryFn{
				IntType: func(env *Environment, a, b Object) Object {
					return NewInt(a.(*Int).Value & b.(*Int).Value)
				},
			},
		},
		{
			name: "|",
			fn: map[TypeFlag]binaryFn{
				IntType: func(env *Environment, a, b Object) Object {
					return NewInt(a.(*Int).Value | b.(*Int).Value)
				},
			},
		},

		{
			name: "^",
			fn: map[TypeFlag]binaryFn{
				IntType: func(env *Environment, a, b Object) Object {
					return NewInt(a.(*Int).Value ^ b.(*Int).Value)
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
				StringType: func(env *Environment, a, b Object) Object {
					return NewBool(a.(*String).Value == b.(*String).Value)
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
				StringType: func(env *Environment, a, b Object) Object {
					return NewBool(a.(*String).Value != b.(*String).Value)
				},
			},
		},

		{
			name: "++",
			fn: map[TypeFlag]binaryFn{
				StringType: func(env *Environment, a, b Object) Object {
					return a.(*String).Concat(b.(*String))
				},
				ArrayType: func(env *Environment, a, b Object) Object {
					// TODO: set up recovery. this causes a panic when used wrong.
					return a.(*Array).Concat(b.(*Array))
				},
				DictType: func(env *Environment, a, b Object) Object {
					return a.(*Dict).Concat(b.(*Dict))
				},
			},
		},

		{
			name: "in",
			fn: map[TypeFlag]binaryFn{
				StringType: func(env *Environment, a, b Object) Object {
					return b.(*String).Contains(a)
				},
				ArrayType: func(env *Environment, a, b Object) Object {
					return b.(*Array).Contains(a)
				},
				DictType: func(env *Environment, a, b Object) Object {
					return b.(*Dict).Contains(a)
				},
			},
		},

		{
			name: "and",
			fn: map[TypeFlag]binaryFn{
				BoolType: func(env *Environment, a, b Object) Object {
					return NewBool(a.(*Bool).Value && b.(*Bool).Value)
				},
			},
		},

		{
			name: "or",
			fn: map[TypeFlag]binaryFn{
				BoolType: func(env *Environment, a, b Object) Object {
					return NewBool(a.(*Bool).Value || b.(*Bool).Value)
				},
			},
		},
	}

	for _, op := range ops {
		BinaryOps[op.name] = op
	}
}
