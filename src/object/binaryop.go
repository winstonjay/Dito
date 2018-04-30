package object

type binaryFn func(Environment, Object, Object) Object

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

func (op *binaryOp) EvalBinary(env Environment, a, b Object) Object {
	which := op.whichType(a.Type(), b.Type())
	a = a.ConvertType(which)
	b = a.ConvertType(which)
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
				CharType: func(env Environment, a, b Object) Object {
					return NewChar(a.(*Char).Value + b.(*Char).Value)
				},
				IntType: func(env Environment, a, b Object) Object {
					return NewInt(a.(*Int).Value + b.(*Int).Value)
				},
				FloatType: func(env Environment, a, b Object) Object {
					return NewFloat(a.(*Float).Value + b.(*Float).Value)
				},
			},
		},

		{
			name: "-",
			fn: map[TypeFlag]binaryFn{
				CharType: func(env Environment, a, b Object) Object {
					return NewChar(a.(*Char).Value - b.(*Char).Value)
				},
				IntType: func(env Environment, a, b Object) Object {
					return NewInt(a.(*Int).Value - b.(*Int).Value)
				},
				FloatType: func(env Environment, a, b Object) Object {
					return NewFloat(a.(*Float).Value - b.(*Float).Value)
				},
			},
		},

		{
			name: "*",
			fn: map[TypeFlag]binaryFn{
				CharType: func(env Environment, a, b Object) Object {
					return NewChar(a.(*Char).Value * b.(*Char).Value)
				},
				IntType: func(env Environment, a, b Object) Object {
					return NewInt(a.(*Int).Value * b.(*Int).Value)
				},
				FloatType: func(env Environment, a, b Object) Object {
					return NewFloat(a.(*Float).Value * b.(*Float).Value)
				},
			},
		},

		{
			name: "/",
			fn: map[TypeFlag]binaryFn{
				CharType: func(env Environment, a, b Object) Object {
					return NewChar(a.(*Char).Value / b.(*Char).Value)
				},
				IntType: func(env Environment, a, b Object) Object {
					return NewInt(a.(*Int).Value / b.(*Int).Value)
				},
				FloatType: func(env Environment, a, b Object) Object {
					return NewFloat(a.(*Float).Value / b.(*Float).Value)
				},
			},
		},
	}

	for _, op := range ops {
		BinaryOps[op.name] = op
	}
}
