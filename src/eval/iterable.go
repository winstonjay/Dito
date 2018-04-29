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
	iter, ok := maybeIter.(object.Iterable)
	if !ok {
		return object.NewError("Index assignement error: wrong type")
	}
	// returns object.Error or nil
	return iter.SetItem(key, val)
}
