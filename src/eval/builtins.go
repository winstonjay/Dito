package eval

import (
	"dito/src/object"
)

// Builtins : map of builtin functions
var Builtins = map[string]*object.Builtin{

	// type conversions.
	object.IntType.String():    &object.Builtin{Fn: typeSwitch(object.IntType)},
	object.FloatType.String():  &object.Builtin{Fn: typeSwitch(object.FloatType)},
	object.StringType.String(): &object.Builtin{Fn: typeSwitch(object.StringType)},
	object.BoolType.String():   &object.Builtin{Fn: typeSwitch(object.BoolType)},
	object.ArrayType.String():  &object.Builtin{Fn: typeSwitch(object.ArrayType)},

	"type": &object.Builtin{Fn: objectType},
	"len":  &object.Builtin{Fn: objectLen},
	"abs":  &object.Builtin{Fn: objectAbs},
}

func typeSwitch(which object.TypeFlag) object.BuiltinFunction {

	return func(args ...object.Object) object.Object {
		n := len(args)
		switch n {
		case 1:
			return args[0].ConvertType(which)
		default:
			return object.NewError(object.InvalidArgLenError, which.String(), 1, n)
		}
	}
}

func objectType(args ...object.Object) object.Object {
	if n := len(args); n > 1 {
		return object.NewError(object.InvalidArgLenError, "type", 1, n)
	}
	return object.NewString(args[0].Type().String())
}

func objectLen(args ...object.Object) object.Object {
	if n := len(args); n != 1 {
		return object.NewError(object.InvalidArgLenError, "len", 1, n)
	}
	iter, ok := args[0].(object.Iterable)
	if ok {
		return iter.Length()
	}
	return object.NewError("Argument '%s' to func 'len' is not a Iterable", args[0].Inspect())
}

func objectAbs(args ...object.Object) object.Object {
	if n := len(args); n != 1 {
		return object.NewError(object.InvalidArgLenError, "abs", 1, n)
	}
	num, ok := args[0].(object.Numeric)
	if ok {
		return num.Abs()
	}
	return object.NewError("Argument '%s' to func 'abs' is not a Numeric", args[0].Inspect())
}
