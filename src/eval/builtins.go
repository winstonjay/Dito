package eval

import (
	"dito/src/object"
	"io"
	"math/rand"
	"os"
	"time"
)

// Builtins : map of builtin functions
var Builtins = map[string]*object.Builtin{

	// type conversions.
	"int": &object.Builtin{
		Name:    "int",
		Fn:      typeSwitch(object.IntType),
		Info:    "Convert value to `Int`",
		ArgC:    1,
		ArgT:    []string{"Atom"},
		ReturnT: "Int",
	},
	"float": &object.Builtin{
		Name:    "float",
		Fn:      typeSwitch(object.FloatType),
		Info:    "Convert value to `Float`",
		ArgC:    1,
		ArgT:    []string{"Atom"},
		ReturnT: "Float",
	},
	"string": &object.Builtin{
		Name:    "string",
		Fn:      typeSwitch(object.StringType),
		Info:    "Convert value to `String`",
		ArgC:    1,
		ArgT:    []string{"Any"},
		ReturnT: "String",
	},
	"bool": &object.Builtin{
		Name:    "bool",
		Fn:      typeSwitch(object.BoolType),
		Info:    "Convert value to `Bool`",
		ArgC:    1,
		ArgT:    []string{"Any"},
		ReturnT: "Bool",
	},
	"array": &object.Builtin{
		Name:    "array",
		Fn:      typeSwitch(object.ArrayType),
		Info:    "Convert value to `Array`",
		ArgC:    1,
		ArgT:    []string{"Iter"},
		ReturnT: "Array",
	},
	"type": &object.Builtin{
		Name:    "type",
		Fn:      objectType,
		Info:    "Reflect a values type",
		ArgC:    1,
		ArgT:    []string{"Any"},
		ReturnT: "String",
	},
	"len": &object.Builtin{
		Name:    "len",
		Fn:      objectLen,
		Info:    "Return the length of an `Iter`",
		ArgC:    1,
		ArgT:    []string{"Iter"},
		ReturnT: "Int",
	},
	"abs": &object.Builtin{
		Name:    "abs",
		Fn:      objectAbs,
		Info:    "Return the absolute value of an `Atom`",
		ArgC:    1,
		ArgT:    []string{"Atom"},
		ReturnT: "Atom",
	},
	"print": &object.Builtin{
		Name:    "print",
		Fn:      objectPrint,
		Info:    "Print a varible number of arguments to the std out.",
		ArgC:    -1,
		ArgT:    []string{"Any..."},
		ReturnT: "None",
	},
	"range": &object.Builtin{
		Name:    "range",
		Fn:      objectRange,
		Info:    "Generate an Array of Int's within a given range.",
		ArgC:    2,
		ArgT:    []string{"Int", "Int"},
		ReturnT: "Array",
	},
	"rand": &object.Builtin{
		Name:    "rand",
		Fn:      objectRand,
		Info:    "Generate a random Float between 0-1.",
		ArgC:    0,
		ArgT:    []string{},
		ReturnT: "Float",
	},
	"time": &object.Builtin{
		Name:    "time",
		Fn:      objectTime,
		Info:    "Return current Unix timestamp as an Int",
		ArgC:    0,
		ArgT:    []string{},
		ReturnT: "Int",
	},
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

// func objectIsIterable(args ...object.Object) object.Object {
// 	if n := len(args); n != 1 {
// 		return object.NewError(object.InvalidArgLenError, "len", 1, n)
// 	}
// 	_, ok := args[0].(object.Iterable)
// 	return object.NewBool(ok)
// }

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

func objectRand(args ...object.Object) object.Object {
	if len(args) != 0 {
		return object.NewError(object.InvalidArgLenError, "rand", 0, len(args))
	}
	return object.NewFloat(rand.Float64())
}

func objectTime(args ...object.Object) object.Object {
	if len(args) != 0 {
		return object.NewError(object.InvalidArgLenError, "rand", 0, len(args))
	}
	return object.NewInt(int(time.Now().Unix()))
}
