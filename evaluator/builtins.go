package evaluator

import (
	"dito/object"
	"fmt"
	"math"
)

// Builtins : map of builtin functions
var Builtins = map[string]*object.Builtin{
	"len":   &object.Builtin{Fn: ditoLen},
	"type":  &object.Builtin{Fn: ditoType},
	"print": &object.Builtin{Fn: ditoPrint},
	"sqrt":  &object.Builtin{Fn: ditoSqrt},
	"iota":  &object.Builtin{Fn: ditoIota},
	"int":   &object.Builtin{Fn: ditoInt},
	// "abs": &object.Builtin{Fn: validDitoAbs},
}

func ditoInt(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("Wrong number of arguments. got=%d, want=1", len(args))
	}
	switch arg := args[0].(type) {
	case *object.Integer:
		return arg
	case *object.Float:
		return object.NewDitoInteger(int64(arg.Value))
	default:
		return newError("Argument to `Int` not supported, got %s", args[0].Type())
	}
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
	for i := arg1; i < arg2; i += arg3 {
		result = append(result, object.NewDitoInteger(i))
	}
	return &object.Array{Len: (arg2 - arg1) / arg3, Elements: result}
}

func ditoType(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("Wrong number of arguments. got=%d, want=1", len(args))
	}
	return &object.DitoString{Value: args[0].Type()}
}

func ditoPrint(args ...object.Object) object.Object {
	for _, arg := range args {
		fmt.Printf("%s ", arg.Inspect())
	}
	fmt.Println()
	return nil
}

func ditoSqrt(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("Wrong number of arguments. got=%d, want=1", len(args))
	}
	switch arg := args[0].(type) {
	case *object.Integer:
		return object.NewDitoFloat(math.Sqrt(float64(arg.Value)))
	case *object.Float:
		return object.NewDitoFloat(math.Sqrt(arg.Value))
	case *object.Array:
		for i, item := range arg.Elements {
			arg.Elements[i] = ditoSqrt(item)
		}
		return arg
	default:
		return newError("Argument to `sqrt` not supported, got %s", args[0].Type())
	}
}

func ditoLen(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("Wrong number of arguments. got=%d, want=1", len(args))
	}
	switch arg := args[0].(type) {
	case *object.DitoString:
		return object.NewDitoInteger(int64(len(arg.Value)))
	case *object.Array:
		return object.NewDitoInteger(arg.Len)
	default:
		return newError("Argument to `Len` not supported, got %s", args[0].Type())
	}
}
