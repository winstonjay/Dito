package evaluator

import (
	"dito/src/ast"
	"dito/src/object"
	"fmt"
)

func evalDefineFunction(fn *ast.Function, env *object.Environment) object.Object {
	obj := object.NewFunction(fn.Parameters, fn.Name.Value, fn.Body, env)
	env.Set(fn.Name.Value, obj)
	return nil
}

func evalFunctionCall(fn ast.Expression, fnArgs []ast.Expression, env *object.Environment) object.Object {
	function := Eval(fn, env)
	if isError(function) {
		fmt.Println(function, "ERRROR")
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
	// we can check out built in stuff here.
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", fn.Type())
	}
}

func extendLambdaEnv(fn *object.LambdaFn, args []object.Object) (*object.Environment, *object.Error) {
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

func extendFunctionEnv(fn *object.Function, args []object.Object) (*object.Environment, *object.Error) {
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
