package eval

import (
	"dito/src/ast"
	"dito/src/object"
	"dito/src/token"
)

func evalIndexExpression(node *ast.IndexExpression, env *object.Environment) object.Object {
	left := Eval(node.Left, env)
	if isError(left) {
		return left
	}
	key := Eval(node.Index, env)
	if isError(key) {
		return left
	}
	if iter, ok := left.(object.Iterable); ok {
		return iter.GetItem(key)
	}
	return object.NewError("Item is does not satisfy Iterable type.")
}

func evalIndexAssignment(node *ast.IndexAssignmentStatement, env *object.Environment) object.Object {
	val := Eval(node.Value, env)
	if isError(val) {
		return val
	}
	maybeIter := Eval(node.IdxExp.Left, env)
	if isError(maybeIter) {
		return maybeIter
	}
	key := Eval(node.IdxExp.Index, env)
	if isError(key) {
		return key
	}

	// TODO INDEX ASSIGNMENT.

	if iter, ok := maybeIter.(object.Iterable); ok {
		// returns object.Error or nil
		return iter.SetItem(key, val)
	}
	return object.NewError("Index assignement error: wrong type")
}

func evalForStatement(fs *ast.ForStatement, env *object.Environment) object.Object {
	if fs.Token == token.IN {
		return evalForIn(fs, env)
	}
	return evalForWhile(fs, env)
}

func evalForIn(fs *ast.ForStatement, env *object.Environment) object.Object {
	var body object.Object
	iter, ok := Eval(fs.Iter, env).(object.Iterable)
	if !ok {
		return iter
	}
	for item := range iter.Iter() {
		env.Set(fs.ID.Value, item, false)
		body = Eval(fs.LoopBody, env)
		if body != nil {
			// the surrounding if is duplicated in isError fn.
			rt := body.Type()
			if rt == object.ErrorType || rt == object.ReturnType {
				return body
			}
		}
	}
	return body
}

func evalForWhile(fs *ast.ForStatement, env *object.Environment) object.Object {
	var body, condition object.Object
	for {
		condition = Eval(fs.Condition, env)
		if !isTrue(condition) {
			break
		}
		if isError(condition) {
			return condition
		}
		body = Eval(fs.LoopBody, env)
		if body != nil {
			rt := body.Type()
			if rt == object.ErrorType || rt == object.ReturnType {
				return body
			}
		}
	}
	return body
}
