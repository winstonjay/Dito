package evaluator

import (
	"dito/ast"
	"dito/object"
	"math"
)

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("Unknown operator: %s %s", operator, right.Type())
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
		if right.Type() == object.IntergerObj {
			val := right.(*object.Integer).Value
			if val == 0 || val == -0 {
				return object.TRUE
			}
		}
		return object.FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.IntergerObj {
		return newError("Unknown operator: -%s", right.Type())
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

// Infix expressions.
func evalInfixEpression(node *ast.InfixExpression, env *object.Environment) object.Object {
	operator := node.Operator
	left := Eval(node.Left, env)
	if isError(left) {
		return left
	}
	right := Eval(node.Right, env)
	if isError(right) {
		return right
	}
	switch {
	case left.Type() != right.Type():
		return newError("Type mismatch: %s %s %s",
			left.Type(), operator, right.Type())
	case left.Type() == object.StringObj && right.Type() == object.StringObj:
		return evalStringExpression(operator, left, right)
	case left.Type() == object.FloatObj && right.Type() == object.FloatObj:
		if operator == "/" || operator == "%" {
			if right.(*object.Float).Value == float64(0) {
				goto ZeroDivisionError
			}
		}
		return evalFloatInfixExpression(operator, left, right)
	case left.Type() == object.IntergerObj && right.Type() == object.IntergerObj:
		if operator == "/" || operator == "%" {
			if right.(*object.Integer).Value == 0 {
				goto ZeroDivisionError
			}
		}
		return evalIntegerInfixExpression(operator, left, right)
	default:
		return newError("Unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
ZeroDivisionError:
	return newError("Zero division error: %s %s %s", left.Inspect(), operator, right.Inspect())
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value
	switch operator {
	case "+":
		if Int64SumOverflows(leftVal, rightVal) {
			goto OverflowError
		}
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		if Int64DiffOverflows(leftVal, rightVal) {
			goto OverflowError
		}
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		if Int64MulOverflows(leftVal, rightVal) {
			goto OverflowError
		}
		return &object.Integer{Value: leftVal * rightVal}
	case "**":
		// overflow handled inside function.
		return IntegerObjPow(leftVal, rightVal)
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "%":
		return &object.Integer{Value: leftVal % rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("Unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}

OverflowError:
	return newError("Overlow/Underlow in operation: %d %s %d",
		leftVal, operator, rightVal)
}

func evalFloatInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Float).Value
	rightVal := right.(*object.Float).Value
	switch operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "**":
		return &object.Float{Value: math.Pow(float64(leftVal), float64(rightVal))}
	case "/":
		return &object.Float{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("Unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalStringExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.DitoString).Value
	rightVal := right.(*object.DitoString).Value
	switch operator {
	case "+":
		return &object.DitoString{Value: leftVal + rightVal}
	case "!=":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	default:
		return newError("Unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}

}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return object.TRUE
	}
	return object.FALSE
}
