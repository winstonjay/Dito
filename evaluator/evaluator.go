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

/*
	Atoms:
*/

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	return newError("Identifier not found: '%s'", node.Value)
}
