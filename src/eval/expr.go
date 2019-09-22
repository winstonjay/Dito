package eval

import (
	"dito/src/ast"
	"dito/src/object"
)

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "not":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return object.NewError("Unknown operator: %s %s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case object.TRUE:
		return object.FALSE
	case object.FALSE:
		return object.TRUE
	case object.NONE:
		return object.TRUE
	default:
		if right.Type() == object.IntType {
			val := right.(*object.Int).Value
			if val == 0 || val == -0 {
				return object.TRUE
			}
		}
		return object.FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() == object.IntType {
		value := right.(*object.Int).Value
		return object.NewInt(-value)
	}
	if right.Type() == object.FloatType {
		value := right.(*object.Float).Value
		return object.NewFloat(-value)
	}
	return object.NewError("Unknown operator: -%s", right.Type())
}

func evalIfElseExpression(node *ast.IfElseExpression, env *object.Environment) object.Object {
	condition := Eval(node.Condition, env)
	if isError(condition) {
		return condition
	}
	if isTrue(condition) {
		return Eval(node.Initial, env)
	}
	return Eval(node.Alternative, env)
}

func evalExpressions(expr []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object
	for _, e := range expr {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

func evalDictExpression(dl *ast.DictLiteral, env *object.Environment) object.Object {
	d := object.NewDict()
	for key, value := range dl.Items {
		d.SetItem(Eval(key, env), Eval(value, env))
	}
	return d
}

func isTrue(obj object.Object) bool {
	switch obj {
	case object.NONE:
		return false
	case object.TRUE:
		return true
	case object.FALSE:
		return false
	default:
		return true
	}
}
