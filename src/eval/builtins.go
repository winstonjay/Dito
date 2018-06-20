package eval

import (
	"dito/src/object"
	"io"
	"os"
)

// Builtins : map of builtin functions
var Builtins = map[string]*object.Builtin{

	// type conversions.
	object.IntType.String():    &object.Builtin{Fn: typeSwitch(object.IntType)},
	object.FloatType.String():  &object.Builtin{Fn: typeSwitch(object.FloatType)},
	object.StringType.String(): &object.Builtin{Fn: typeSwitch(object.StringType)},
	object.BoolType.String():   &object.Builtin{Fn: typeSwitch(object.BoolType)},
	object.ArrayType.String():  &object.Builtin{Fn: typeSwitch(object.ArrayType)},

	"type":       &object.Builtin{Fn: objectType},
	"len":        &object.Builtin{Fn: objectLen},
	"abs":        &object.Builtin{Fn: objectAbs},
	"print":      &object.Builtin{Fn: objectPrint},
	"isIterable": &object.Builtin{Fn: objectIsIterable},
	"range":      &object.Builtin{Fn: objectRange},
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

func objectPrint(args ...object.Object) object.Object {
	for i, arg := range args {
		io.WriteString(os.Stdout, arg.Inspect())
		if i < len(args)-1 {
			io.WriteString(os.Stdout, " ")
		}
	}
	io.WriteString(os.Stdout, "\n")
	return object.NONE
}

func objectIsIterable(args ...object.Object) object.Object {
	if n := len(args); n != 1 {
		return object.NewError(object.InvalidArgLenError, "len", 1, n)
	}
	_, ok := args[0].(object.Iterable)
	return object.NewBool(ok)
}

func objectRange(args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.NewError(object.InvalidArgLenError, "range", 2, len(args))
	}
	min, ok := args[0].(*object.Int)
	if !ok {
		return object.NewError("Invalid args[0] to `range` function. got=%T", args[0])
	}
	max, ok := args[1].(*object.Int)
	if !ok {
		return object.NewError("Invalid args[1] to `range` function. got=%T", args[1])
	}
	iter := make([]object.Object, max.Value-min.Value)
	i := 0
	for v := min.Value; v < max.Value; v++ {
		iter[i] = object.NewInt(v)
		i++
	}
	return object.NewArray(iter, -1)
}
