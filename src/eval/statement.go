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
