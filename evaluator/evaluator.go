package evaluator

import (
	"dito/ast"
	"dito/object"
	"dito/token"
)

// Eval :
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements, env)

	// Assignments
	case *ast.AssignmentStatement:
		return evalAssignment(node, env)

	// Expressions
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		return evalInfixEpression(node, env)

	case *ast.IfStatement:
		return evalIfStatement(node, env)
	case *ast.ForStatement:
		return evalForStatement(node, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	// Atoms
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.StringLiteral:
		return &object.DitoString{Value: node.Value}
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}
	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value)
	}
	return nil
}

func evalProgram(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range stmts {
		result = Eval(statement, env)
		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}

/*
	Assignment:
*/

func evalAssignment(node *ast.AssignmentStatement, env *object.Environment) object.Object {
	val := Eval(node.Value, env)
	if isError(val) {
		return val
	}
	if node.Token == token.REASSIGN {
		if ident := evalIdentifier(node.Name, env); isError(ident) {
			return ident
		}
	}
	env.Set(node.Name.Value, val)
	return nil
}

func evalIfStatement(ie *ast.IfStatement, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}
	if isTrue(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return object.NONE
	}
}

func evalForStatement(fs *ast.ForStatement, env *object.Environment) object.Object {
	iterCount := 0
	var body, condition object.Object
	for {
		condition = Eval(fs.Condition, env)
		if !isTrue(condition) {
			break
		}
		if isError(condition) {
			return condition
		}
		if iterCount > 10000 {
			return newError("Max iteration limit reached.")
		}
		iterCount++
		body = Eval(fs.LoopBody, env)
	}
	return body
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

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range block.Statements {
		result = Eval(statement, env)
		if result != nil {
			rt := result.Type()
			if rt == object.ErrorObj {
				return result
			}
		}
	}
	return result
}

/*
	Atoms:
*/

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	return newError("Identifier not found: '%s'", node.Value)
}
