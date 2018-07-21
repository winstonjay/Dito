package eval

import (
	"dito/src/ast"
	"dito/src/object"
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
	var body object.Object
	iter, ok := Eval(fs.Iter, env).(object.Iterable)
	if !ok {
		return object.NewError("Loop error")
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
