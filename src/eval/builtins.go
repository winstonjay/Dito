package eval

import (
	"dito/src/object"
	"io"
	"math"
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
	"error": &object.Builtin{
		Name:    "error",
		Fn:      objectError,
		Info:    "create a new error message with a string.",
		ArgC:    1,
		ArgT:    []string{"String"},
		ReturnT: "Error",
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
	"sin": &object.Builtin{
		Name:    "sin",
		Fn:      objectSin,
		Info:    "Return the sine of x radians of an `Atom`",
		ArgC:    1,
		ArgT:    []string{"Atom"},
		ReturnT: "Float",
	},
	"tan": &object.Builtin{
		Name:    "tan",
		Fn:      objectTan,
		Info:    "Return the tangent of x radians of an `Atom`",
		ArgC:    1,
		ArgT:    []string{"Atom"},
		ReturnT: "Float",
	},
	"cos": &object.Builtin{
		Name:    "cos",
		Fn:      objectCos,
		Info:    "Return the cosine of x radians of an `Atom`",
		ArgC:    1,
		ArgT:    []string{"Atom"},
		ReturnT: "Float",
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
	"random": &object.Builtin{
		Name:    "random",
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
	"sleep": &object.Builtin{
		Name:    "sleep",
		Fn:      objectSleep,
		Info:    "Pause execution for x milliseconds",
		ArgC:    1,
		ArgT:    []string{"Atom"},
		ReturnT: "None",
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

func objectCos(args ...object.Object) object.Object {
	if n := len(args); n != 1 {
		return object.NewError(object.InvalidArgLenError, "cos", 1, n)
	}
	switch v := args[0].(type) {
	case *object.Float:
		return object.NewFloat(math.Cos(v.Value))
	case *object.Int:
		return object.NewFloat(math.Cos(float64(v.Value)))
	default:
		return object.NewError("Argument '%s' to func 'cos' is invalid", args[0].Inspect())
	}
}

func objectTan(args ...object.Object) object.Object {
	if n := len(args); n != 1 {
		return object.NewError(object.InvalidArgLenError, "tan", 1, n)
	}
	switch v := args[0].(type) {
	case *object.Float:
		return object.NewFloat(math.Tan(v.Value))
	case *object.Int:
		return object.NewFloat(math.Tan(float64(v.Value)))
	default:
		return object.NewError("Argument '%s' to func 'tan' is invalid", args[0].Inspect())
	}
}

func objectSin(args ...object.Object) object.Object {
	if n := len(args); n != 1 {
		return object.NewError(object.InvalidArgLenError, "sin", 1, n)
	}
	switch v := args[0].(type) {
	case *object.Float:
		return object.NewFloat(math.Sin(v.Value))
	case *object.Int:
		return object.NewFloat(math.Sin(float64(v.Value)))
	default:
		return object.NewError("Argument '%s' to func 'sin' is invalid", args[0].Inspect())
	}
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

func objectError(args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError(object.InvalidArgLenError, "error", 2, len(args))
	}
	if args[0].Type() != object.StringType {
		return object.NewError("Argument to `error` not supported, got %s", args[0].Type())
	}
	return object.NewError(args[0].Inspect())
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
	if max.Value-min.Value < 0 {
		return object.NewError("Invalid args to `range` function. a > b.")
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
		return object.NewError(object.InvalidArgLenError, "random", 0, len(args))
	}
	return object.NewFloat(rand.Float64())
}

func objectTime(args ...object.Object) object.Object {
	if len(args) != 0 {
		return object.NewError(object.InvalidArgLenError, "time", 0, len(args))
	}
	return object.NewInt(int(time.Now().Unix()))
}

func objectSleep(args ...object.Object) object.Object {
	switch arg := args[0].(type) {
	case *object.Int:
		time.Sleep(time.Duration(arg.Value) * time.Millisecond)
		return object.NONE
	case *object.Float:
		time.Sleep(time.Duration(arg.Value) * time.Millisecond)
		return object.NONE
	default:
		return object.NewError("Argument to `sleep` not supported, got %s", args[0].Type())
	}
}
