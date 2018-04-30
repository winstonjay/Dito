package eval

import (
	"dito/src/ast"
	"dito/src/object"
	"dito/src/token"
)

func evalAssignment(node *ast.AssignmentStatement, env *object.Environment) object.Object {
	val := Eval(node.Value, env)
	if isError(val) {
		return val
	}
	// just ordinary assignement
	if node.Token == token.REASSIGN || node.Token == token.NEWASSIGN {
		env.Set(node.Name.Value, val)
		return nil
	}
	return object.NewError("Assignement Error: unknown")
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range block.Statements {
		result = Eval(statement, env)
		if result != nil {
			rt := result.Type()
			if rt == object.ErrorType || rt == object.ReturnType {
				return result
			}
		}
	}
	return result
}
