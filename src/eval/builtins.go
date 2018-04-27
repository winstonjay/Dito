package eval

import "dito/src/object"

// Builtins : map of builtin functions
var Builtins = map[string]*object.Builtin{

	// type conversions.
	object.IntType.String():    &object.Builtin{Fn: typeSwitch(object.IntType)},
	object.FloatType.String():  &object.Builtin{Fn: typeSwitch(object.FloatType)},
	object.StringType.String(): &object.Builtin{Fn: typeSwitch(object.StringType)},
	object.BoolType.String():   &object.Builtin{Fn: typeSwitch(object.BoolType)},

	"type": &object.Builtin{Fn: objectType},
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
