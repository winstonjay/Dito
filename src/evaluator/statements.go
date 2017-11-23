package evaluator

import (
	"dito/src/ast"
	"dito/src/object"
)

/*
	Statements:
*/

func evalAssignment(node *ast.AssignmentStatement, env *object.Environment) object.Object {
	val := Eval(node.Value, env)
	if isError(val) {
		return val
	}
	// * Currently don't enforce scope initialation assignment `:=`
	// * as we are only having single statement functions right now.
	// if node.Token == token.REASSIGN {
	// 	if ident := evalIdentifier(node.Name, env); isError(ident) {
	// 		return ident
	// 	}
	// }
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
	if fs.ID == nil {
		goto whileStmt
	}
	goto forInStmt

whileStmt:
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
		if body != nil {
			if body.Type() == object.ReturnObj {
				return body
			}
		}
	}
	return body
forInStmt:
	iter := Eval(fs.Iter, env).(*object.Array)
	for _, item := range iter.Elements {
		env.Set(fs.ID.Value, item)
		body = Eval(fs.LoopBody, env)
		if body != nil {
			// the surrounding if is duplicated in isError fn.
			if body.Type() == object.ReturnObj || isError(body) {
				return body
			}
		}
	}
	return body
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range block.Statements {
		result = Eval(statement, env)
		if result != nil {
			rt := result.Type()
			if rt == object.ErrorObj || rt == object.ReturnObj {
				return result
			}
		}
	}
	return result
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
