package builtins

import "dito/object"

var builtins = map[string]*object.Builtin{
	// "len": &object.Builtin{Fn: validDitoLength},
	"type": &object.Builtin{Fn: ditoType},
	// "abs": &object.Builtin{Fn: validDitoAbs},
}

func ditoType(item ...object.Object) object.Object {
	return &object.DitoString{Value: item[0].Type()}
}

// func validDitoLength(args ...object.Object) object.Object {
// 	if len(args) != 1 {
// 		return newError("wrong number of arguments. got=%d, want=1", len(args))
// 	}
// 	switch arg := args[0].(type) {
// 	case *object.String:
// 		return &object.Integer{Value: arg.Len}
// 	case *object.Array:
// 		return &object.Integer{Value: arg.Len}
// 	default:
// 		return newError("argument to `len` not supported, got %s", args[0].Type())
// 	}
// }

// // DitoAbs :
// func validDitoAbs(args ...object.Object) object.Object {
// 	if len(args) != 1 {
// 		return newError("wrong number of arguments. got=%d, want=1", len(args))
// 	}
// 	switch arg := args[0].(type) {
// 	case *object.Integer:
// 		return &object.Integer{Value: builtin.DitoIntAbs(arg.Value)}
// 	case *object.Float:
// 		return &object.Float{Value: builtin.DitoFloatAbs(arg.Value)}
// 	default:
// 		return newError("argument to `len` not supported, got %s", args[0].Type())
// 	}
// }
