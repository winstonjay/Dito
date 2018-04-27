package object

// Environment : Holds the enviroment varibles created by the user.
type Environment struct {
	store map[string]Object
	outer *Environment
}

// NewEnvironment : Define a new enviroment scope.
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

// NewEnclosedEnviroment : Define a new enviroment scope within another.
func NewEnclosedEnviroment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// Get : get a varible inside the current scope
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

// Set : set a varible inside the current scope.
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
