package evaluator

import (
	"dito/src/ast"
	"dito/src/object"
	"fmt"
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
		return object.NewDitoInteger(-value)
	}
	if right.Type() == object.FloatObj {
		value := right.(*object.Float).Value
		return object.NewDitoFloat(-value)
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
	case isNumericType(left) && isNumericType(right):
		if left.Type() == object.IntergerObj && right.Type() == object.IntergerObj {
			return evalIntegerInfixExpression(operator, left, right)
		}
		return evalFloatInfixExpression(operator, left, right)
	case left.Type() != right.Type():
		return newError("Type mismatch: %s %s %s", left.Type(), operator, right.Type())
	case left.Type() == object.StringObj && right.Type() == object.StringObj:
		return evalStringExpression(operator, left, right)
	case operator == "==":
		return object.NewDitoBoolean(left == right)
	case operator == "!=":
		return object.NewDitoBoolean(left != right)
	default:
		return newError("Unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value
	switch operator {
	case "+":
		return object.NewDitoInteger(leftVal + rightVal)
	case "-":
		return object.NewDitoInteger(leftVal - rightVal)
	case "*":
		return object.NewDitoInteger(leftVal * rightVal)
	case "**":
		return IntegerObjPow(leftVal, rightVal)
	case "/":
		if rightVal == 0 {
			return newError("Zero division error: %s %s %s",
				left.Inspect(), operator, right.Inspect())
		}
		if leftVal%rightVal == 0 {
			return object.NewDitoInteger(leftVal / rightVal)
		}
		return object.NewDitoFloat(float64(leftVal) / float64(rightVal))
	case "%":
		if rightVal == 0 {
			return newError("Zero division error: %s %s %s",
				left.Inspect(), operator, right.Inspect())
		}
		return object.NewDitoInteger(leftVal % rightVal)
	case "==":
		return object.NewDitoBoolean(leftVal == rightVal)
	case "!=":
		return object.NewDitoBoolean(leftVal != rightVal)
	case "<":
		return object.NewDitoBoolean(leftVal < rightVal)
	case "<=":
		return object.NewDitoBoolean(leftVal <= rightVal)
	case ">":
		return object.NewDitoBoolean(leftVal > rightVal)
	case ">=":
		return object.NewDitoBoolean(leftVal >= rightVal)
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
		return object.NewDitoFloat(leftVal + rightVal)
	case "-":
		return object.NewDitoFloat(leftVal - rightVal)
	case "*":
		return object.NewDitoFloat(leftVal * rightVal)
	case "**":
		return object.NewDitoFloat(math.Pow(leftVal, rightVal))
	case "/":
		if rightVal == 0 {
			return newError("Zero division error: %s %s %s",
				left.Inspect(), operator, right.Inspect())
		}
		return object.NewDitoFloat(leftVal / rightVal)
	case "==":
		return object.NewDitoBoolean(leftVal == rightVal)
	case "!=":
		return object.NewDitoBoolean(leftVal != rightVal)
	case "<":
		return object.NewDitoBoolean(leftVal < rightVal)
	case "<=":
		return object.NewDitoBoolean(leftVal <= rightVal)
	case ">":
		return object.NewDitoBoolean(leftVal > rightVal)
	case ">=":
		return object.NewDitoBoolean(leftVal >= rightVal)
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
	case "==":
		return object.NewDitoBoolean(leftVal == rightVal)
	case "!=":
		return object.NewDitoBoolean(leftVal != rightVal)
	default:
		return newError("Unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
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

func isNumericType(obj object.Object) bool {
	return obj.Type() == object.FloatObj || obj.Type() == object.IntergerObj
}

func evalIndexExpression(node *ast.IndexExpression, env *object.Environment) object.Object {
	left := Eval(node.Left, env)
	if isError(left) {
		return left
	}
	index := Eval(node.Index, env)
	if isError(index) {
		return left
	}
	if left.Type() != object.ArrayObj || index.Type() != object.IntergerObj {
		return newError("Index operator not supported: %s[%s].", left.Type(), index.Type())
	}
	arrayObject := left.(*object.Array)
	idx := index.(*object.Integer).Value
	size := arrayObject.Len - 1
	if idx < 0 {
		idx += (size + 1)
	}
	fmt.Printf("len=%d, index=%d\n", size, idx)
	if idx < 0 || idx > size {
		return newError("Index out of range: len=%d, index=%d.", size, idx)
	}
	return arrayObject.Elements[idx]
}

// x := iota(10)
