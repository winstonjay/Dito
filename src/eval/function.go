package eval

import (
	"dito/src/ast"
	"dito/src/object"
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

func extendLambdaEnv(fn *object.Lambda, args []object.Object) (*object.Environment, *object.Error) {
	env := object.NewEnclosedEnviroment(fn.Env)
	if len(fn.Parameters) != len(args) {
		return nil, object.NewError("Wrong number of function args. Want=%d, Got=%d. in %s",
			len(fn.Parameters), len(args), fn.Inspect())
	}
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx], false)
	}
	return env, nil
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Lambda:
		extendedEnv, err := extendLambdaEnv(fn, args)
		if err != nil {
			return err
		}
		evaluated := Eval(fn.Expr, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Function:
		extendedEnv, err := extendFunctionEnv(fn, args)
		if err != nil {
			return err
		}
		evaluated := Eval(fn.Stmts, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return object.NewError("not a function: %s", fn.Type())
	}
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func extendFunctionEnv(fn *object.Function, args []object.Object) (*object.Environment, *object.Error) {
	env := object.NewEnclosedEnviroment(fn.Env)
	if len(fn.Parameters) != len(args) {
		return nil, object.NewError("Wrong number of function args. Want=%d, Got=%d.",
			len(fn.Parameters), len(args))
	}
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx], false)
	}
	return env, nil
}
