package evaluator

import (
	"dito/ast"
	"dito/object"
)

func evalFunctionCall(fn ast.Expression, fnArgs []ast.Expression, env *object.Environment) object.Object {
	function := Eval(fn, env)
	if isError(function) {
		return function
	}
	args := evalExpressions(fnArgs, env)
	if len(args) == 1 && isError(args[0]) {
		return args[0]
	}
	return applyFunction(function, args)
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.LambdaFn:
		extendedEnv, err := extendFunctionEnv(fn, args)
		if err != nil {
			return err
		}
		evaluated := Eval(fn.Expr, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.LambdaFn, args []object.Object) (*object.Environment, *object.Error) {
	env := object.NewEnclosedEnviroment(fn.Env)
	if len(fn.Parameters) != len(args) {
		return nil, newError("Wrong number of function args. Want=%d, Got=%d.",
			len(fn.Parameters), len(args))
	}
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}
	return env, nil
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}
