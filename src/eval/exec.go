package eval

import (
	"dito/src/ast"
	"dito/src/object"
)

// Eval :
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements, env)

	// // Statements
	// case *ast.AssignmentStatement:
	// 	return evalAssignment(node, env)
	// case *ast.IndexAssignmentStatement:
	// 	return evalIndexAssignment(node, env)
	// case *ast.ReturnStatement:
	// 	return &object.ReturnValue{Value: Eval(node.Value, env)}
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	// case *ast.IfStatement:
	// 	return evalIfStatement(node, env)
	// case *ast.ForStatement:
	// 	return evalForStatement(node, env)
	// case *ast.BlockStatement:
	// 	return evalBlockStatement(node, env)
	// case *ast.ImportStatement:
	// 	return evalImportStatement(node, env)

	// Expressions
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		return evalInfixExpression(node, env)
	case *ast.IfElseExpression:
		return evalIfElseExpression(node, env)
	// case *ast.IndexExpression:
	// 	return evalIndexExpression(node, env)

	// // Functions
	// case *ast.Function:
	// 	return evalDefineFunction(node, env)
	// case *ast.LambdaFunction:
	// 	return object.NewDitoLambdaFn(node.Parameters, node.Expr, env)
	case *ast.CallExpression:
		return evalFunctionCall(node.Function, node.Arguments, env)

	// Atoms
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.StringLiteral:
		return object.NewString(node.Value)
	case *ast.IntegerLiteral:
		return object.NewInt(node.Value)
	case *ast.FloatLiteral:
		return object.NewFloat(node.Value)
	case *ast.BooleanLiteral:
		return object.NewBool(node.Value)
		// case *ast.ArrayLiteral:
		// 	return evalArray(node, env)
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

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	if builtin, ok := Builtins[node.Value]; ok {
		return builtin
	}
	return object.NewError("Identifier not found: '%s'", node.Value)
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ErrorType
	}
	return false
}
