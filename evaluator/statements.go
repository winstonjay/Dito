package evaluator

import (
	"dito/ast"
	"dito/object"
	"dito/token"
)

/*
	Statements:
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
