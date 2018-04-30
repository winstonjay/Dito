package eval

import (
	"dito/src/ast"
	"dito/src/object"
)

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
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

// Infix expressions.
func evalInfixExpression(node *ast.InfixExpression, env *object.Environment) object.Object {
	op := object.BinaryOps[node.Operator]
	if op == nil {
		return object.NewError("Unknown Binary operation: '%s'", node.Operator)
	}
	left := Eval(node.Left, env)
	if isError(left) {
		return left
	}
	right := Eval(node.Right, env)
	if isError(right) {
		return right
	}
	return op.EvalBinary(env, left, right)
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

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object
	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
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
