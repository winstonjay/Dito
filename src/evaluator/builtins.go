package evaluator

import (
	"dito/src/object"
	"io"
	"math"
	"os"
	"time"
)

// ######## TODO: enforce arg lens with a generic function.

// Builtins : map of builtin functions
var Builtins = map[string]*object.Builtin{
	"len":    &object.Builtin{Fn: ditoLen},
	"type":   &object.Builtin{Fn: ditoType},
	"print":  &object.Builtin{Fn: ditoPrint},
	"sqrt":   &object.Builtin{Fn: ditoSqrt},
	"iota":   &object.Builtin{Fn: ditoIota},
	"int":    &object.Builtin{Fn: ditoInt},
	"string": &object.Builtin{Fn: ditoString},
	"sleep":  &object.Builtin{Fn: ditoSleep},
	"abs":    &object.Builtin{Fn: ditoAbs},
}

func ditoInt(args ...object.Object) object.Object {
	switch arg := args[0].(type) {
	case *object.Integer:
		return arg
	case *object.Float:
		return object.NewDitoInteger(int64(arg.Value))
	default:
		return newError("Argument to `Int` not supported, got %s", args[0].Type())
	}
}

func ditoString(args ...object.Object) object.Object {
	return &object.DitoString{Value: args[0].Inspect()}
}

func ditoIota(args ...object.Object) object.Object {
	var arg1, arg2, arg3 int64
	switch len(args) {
	case 1:
		if args[0].Type() != object.IntergerObj {
			return newError("Argument to `iota` not supported, got=%s. want=Int", args[0].Type())
		}
		arg1 = 0
		arg2 = args[0].(*object.Integer).Value
		arg3 = 1
	case 2:
		if args[0].Type() != object.IntergerObj {
			return newError("Argument to `iota` not supported, got=%s. want=Int", args[0].Type())
		}
		if args[1].Type() != object.IntergerObj {
			return newError("Argument to `iota` not supported, got=%s. want=Int", args[1].Type())
		}
		arg1 = args[0].(*object.Integer).Value
		arg2 = args[1].(*object.Integer).Value
		arg3 = 1
	case 3:
		if args[0].Type() != object.IntergerObj {
			return newError("Argument to `iota` not supported, got=%s. want=Int", args[0].Type())
		}
		if args[1].Type() != object.IntergerObj {
			return newError("Argument to `iota` not supported, got=%s. want=Int", args[1].Type())
		}
		if args[2].Type() != object.IntergerObj {
			return newError("Argument to `iota` not supported, got=%s. want=Int", args[2].Type())
		}
		arg1 = args[0].(*object.Integer).Value
		arg2 = args[1].(*object.Integer).Value
		arg3 = args[2].(*object.Integer).Value
	default:
		return newError("`iota`: wrong number of args. want=(1-3) got=%d", len(args))
	}
	var result []object.Object
	if arg1 < arg2 {
		if arg3 <= 0 {
			return newError("`iota`: Invalid loop increment. want= > 0. got=%d", arg3)
		}
		for i := arg1; i < arg2; i += arg3 {
			result = append(result, object.NewDitoInteger(i))
		}
	} else {
		if arg3 >= 0 {
			return newError("`iota`: Invalid loop increment. want= < 0. got=%d", arg3)
		}
		for i := arg1; i > arg2; i += arg3 {
			result = append(result, object.NewDitoInteger(i))
		}
	}
	return &object.Array{Len: (arg2 - arg1) / arg3, Elements: result}
}

func ditoType(args ...object.Object) object.Object {
	return &object.DitoString{Value: args[0].Type()}
}

func ditoPrint(args ...object.Object) object.Object {
	for i, arg := range args {
		io.WriteString(os.Stdout, arg.Inspect())
		if i < len(args)-1 {
			io.WriteString(os.Stdout, ", ")
		}
	}
	io.WriteString(os.Stdout, "\n")
	return nil
}

func ditoSqrt(args ...object.Object) object.Object {
	switch arg := args[0].(type) {
	case *object.Integer:
		return object.NewDitoFloat(math.Sqrt(float64(arg.Value)))
	case *object.Float:
		return object.NewDitoFloat(math.Sqrt(arg.Value))
	default:
		return newError("Argument to `sqrt` not supported, got %s", args[0].Type())
	}
}

func ditoLen(args ...object.Object) object.Object {
	switch arg := args[0].(type) {
	case *object.DitoString:
		return object.NewDitoInteger(int64(len(arg.Value)))
	case *object.Array:
		return object.NewDitoInteger(arg.Len)
	default:
		return newError("Argument to `Len` not supported, got %s", args[0].Type())
	}
}

func ditoAbs(args ...object.Object) object.Object {
	switch arg := args[0].(type) {
	case *object.Integer:
		return object.NewDitoInteger(int64(math.Abs(float64(arg.Value))))
	case *object.Float:
		return object.NewDitoFloat(math.Abs(arg.Value))
	default:
		return newError("Argument to `abs` not supported, got %s", args[0].Type())
	}
}

func ditoSleep(args ...object.Object) object.Object {
	switch arg := args[0].(type) {
	case *object.Integer:
		time.Sleep(time.Duration(arg.Value) * time.Millisecond)
		return object.NONE
	case *object.Float:
		time.Sleep(time.Duration(arg.Value) * time.Millisecond)
		return object.NONE
	default:
		return newError("Argument to `sleep` not supported, got %s", args[0].Type())
	}

}
