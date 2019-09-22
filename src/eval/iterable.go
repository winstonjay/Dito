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
		return key
	}
	if iter, ok := left.(object.Iterable); ok {
		return iter.GetItem(key)
	}
	return object.NewError("Item is not Iterable")
}

func evalSliceExpression(node *ast.SliceExpression, env *object.Environment) object.Object {
	left := Eval(node.Left, env)
	if isError(left) {
		return left
	}
	start := Eval(node.S, env)
	if isError(start) {
		return start
	}
	end := Eval(node.E, env)
	if isError(end) {
		return end
	}
	if iter, ok := left.(object.Iterable); ok {
		return iter.Slice(start, end)
	}
	return object.NewError("Item is not Iterable")
}

func evalIndexAssignment(node *ast.IndexAssignmentStatement, env *object.Environment) object.Object {
	right := Eval(node.Value, env)
	if isError(right) {
		return right
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
		return object.NewError("Index assignment error: wrong type")
	}
	if node.Token == token.ASSIGN {
		iter.SetItem(key, right)
		return object.NONE
	}
	// TODO this is just implemented as a quick fix.
	left := iter.GetItem(key)
	opString := node.Token.String()
	op := object.BinaryOps[opString[:len(opString)-1]]
	if op == nil {
		return object.NewError("Unknown in place binary op: '%s'", opString)
	}
	val := op.EvalBinary(env, left, right)
	if isError(val) {
		return val
	}
	iter.SetItem(key, val)
	return object.NONE
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
