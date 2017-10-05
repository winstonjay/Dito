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
	if right.Type() == object.IntergerObj {
		value := right.(*object.Integer).Value
		return &object.Integer{Value: -value}
	}
	if right.Type() == object.FloatObj {
		value := right.(*object.Float).Value
		return &object.Float{Value: -value}
	}
	return newError("Unknown operator: -%s", right.Type())
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
	case left.Type() == object.FloatObj || left.Type() == object.IntergerObj &&
		right.Type() == object.FloatObj || right.Type() == object.IntergerObj:

		if left.Type() == object.IntergerObj && right.Type() == object.IntergerObj {
			return evalIntegerInfixExpression(operator, left, right)
		}
		return evalFloatInfixExpression(operator, left, right)

	case left.Type() != right.Type():
		return newError("Type mismatch: %s %s %s", left.Type(), operator, right.Type())

	case left.Type() == object.StringObj && right.Type() == object.StringObj:
		return evalStringExpression(operator, left, right)

	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)

	default:
		return newError("Unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value
	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "**":
		return IntegerObjPow(leftVal, rightVal)

	case "/":
		if rightVal == 0 {
			return newError("Zero division error: %s %s %s",
				left.Inspect(), operator, right.Inspect())
		}
		if leftVal%rightVal == 0 {
			return &object.Integer{Value: leftVal / rightVal}
		}
		return &object.Float{Value: float64(leftVal) / float64(rightVal)}
	case "%":
		if rightVal == 0 {
			return newError("Zero division error: %s %s %s",
				left.Inspect(), operator, right.Inspect())
		}
		return &object.Integer{Value: leftVal % rightVal}
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	default:
		return newError("Unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalFloatInfixExpression(operator string, left, right object.Object) object.Object {
	var leftVal, rightVal float64

	// Handle the type promotion if there is oone.
	if left.Type() == object.FloatObj {
		leftVal = left.(*object.Float).Value
	} else {
		leftVal = float64(left.(*object.Integer).Value)
	}
	if right.Type() == object.FloatObj {
		rightVal = right.(*object.Float).Value
	} else {
		rightVal = float64(right.(*object.Integer).Value)
	}

	switch operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "**":
		return &object.Float{Value: math.Pow(leftVal, rightVal)}

	case "/":
		if rightVal == 0 {
			return newError("Zero division error: %s %s %s",
				left.Inspect(), operator, right.Inspect())
		}
		return &object.Float{Value: leftVal / rightVal}
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
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
