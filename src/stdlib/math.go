package stdlib

// func ditoSin(args ...object.Object) object.Object {
// 	switch arg := args[0].(type) {
// 	case *object.Integer:
// 		return object.NewDitoFloat(math.Sin(float64(arg.Value)))
// 	case *object.Float:
// 		return object.NewDitoFloat(math.Sin(arg.Value))
// 	default:
// 		return newError("Argument to `Sin` not supported, got %s", args[0].Type())
// 	}
// }

// func ditoCos(args ...object.Object) object.Object {
// 	switch arg := args[0].(type) {
// 	case *object.Integer:
// 		return object.NewDitoFloat(math.Cos(float64(arg.Value)))
// 	case *object.Float:
// 		return object.NewDitoFloat(math.Cos(arg.Value))
// 	default:
// 		return newError("Argument to `Cos` not supported, got %s", args[0].Type())
// 	}
// }

// func ditoLog(args ...object.Object) object.Object {
// 	switch arg := args[0].(type) {
// 	case *object.Integer:
// 		return object.NewDitoFloat(math.Log(float64(arg.Value)))
// 	case *object.Float:
// 		return object.NewDitoFloat(math.Log(arg.Value))
// 	default:
// 		return newError("Argument to `log` not supported, got %s", args[0].Type())
// 	}
// }

// func ditoLog10(args ...object.Object) object.Object {
// 	switch arg := args[0].(type) {
// 	case *object.Integer:
// 		return object.NewDitoFloat(math.Log10(float64(arg.Value)))
// 	case *object.Float:
// 		return object.NewDitoFloat(math.Log10(arg.Value))
// 	default:
// 		return newError("Argument to `log10` not supported, got %s", args[0].Type())
// 	}
// }

// func ditoLog2(args ...object.Object) object.Object {
// 	switch arg := args[0].(type) {
// 	case *object.Integer:
// 		return object.NewDitoFloat(math.Log2(float64(arg.Value)))
// 	case *object.Float:
// 		return object.NewDitoFloat(math.Log2(arg.Value))
// 	default:
// 		return newError("Argument to `log2` not supported, got %s", args[0].Type())
// 	}
// }

// func ditoTan(args ...object.Object) object.Object {
// 	switch arg := args[0].(type) {
// 	case *object.Integer:
// 		return object.NewDitoFloat(math.Tan(float64(arg.Value)))
// 	case *object.Float:
// 		return object.NewDitoFloat(math.Tan(arg.Value))
// 	default:
// 		return newError("Argument to `Tan` not supported, got %s", args[0].Type())
// 	}
// }
