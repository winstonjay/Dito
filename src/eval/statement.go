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
	env.Set(node.Name.Value, val, node.Token != token.LET)
	return nil
}

func evalReAssign(node *ast.ReAssignStatement, env *object.Environment) object.Object {
	v, ok := env.GetVar(node.Name.Value)
	if !ok {
		return object.NewError("identifier not found: '%s'", node.Name.Value)
	}
	if !v.IsMutable() {
		return object.NewError("identifier '%s' has a immutable value.", node.Name.Value)
	}
	right := Eval(node.Value, env)
	if isError(right) {
		return right
	}
	// check to see if we are just doing a simple assignment.
	if node.Token == token.ASSIGN {
		env.Set(node.Name.Value, right, true)
		return nil
	}

	// TODO: This is just using regular binary ops instead of in place operations.
	// would it be quick to do this otherwise.
	left := v.Unpack()
	opString := node.Token.String()

	op := object.BinaryOps[opString[:len(opString)-1]]
	if op == nil {
		return object.NewError("Unknown inplace binary op: '%s'", opString)
	}
	val := op.EvalBinary(env, left, right)
	if isError(val) {
		return val
	}
	env.Set(node.Name.Value, val, true)
	return nil
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
