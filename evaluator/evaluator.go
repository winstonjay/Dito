package evaluator

import (
	"dito/ast"
	"dito/object"
)

// Eval :
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements, env)

	// Statements
	case *ast.AssignmentStatement:
		return evalAssignment(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.IfStatement:
		return evalIfStatement(node, env)
	case *ast.ForStatement:
		return evalForStatement(node, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	// Expressions

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		return evalInfixEpression(node, env)
	case *ast.IfElseExpression:
		return evalIfElseExpression(node, env)

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
	Atoms:
*/

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	return newError("Identifier not found: '%s'", node.Value)
}
