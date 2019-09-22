package object

// Variable : symbol table entry, keeps track whether varible is mutable or not.
type Variable struct {
	value   Object
	mutable bool
}

// IsMutable : can it be changed
func (v *Variable) IsMutable() bool {
	return v.mutable
}

// Unpack : can it be changed
func (v *Variable) Unpack() Object {
	return v.value
}

// Environment : Holds the environment variables created by the user. Pretty much
// a symbol table.
type Environment struct {
	store map[string]Variable
	outer *Environment
}

// InitialEnvironment : Define the initial environment scope. with system variables etc.
func InitialEnvironment() *Environment {
	return &Environment{
		store: map[string]Variable{
			"STDIN":  Variable{value: STDIN, mutable: false},
			"STDOUT": Variable{value: STDOUT, mutable: false},
			"STDERR": Variable{value: STDERR, mutable: false},
		},
	}
}

// NewEnvironment : Define a new environment scope.
func NewEnvironment() *Environment {
	s := make(map[string]Variable)
	return &Environment{store: s, outer: nil}
}

// NewEnclosedEnvironment : Define a new environment scope within another.
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// Get : get a variable inside the current scope
func (e *Environment) Get(name string) (Object, bool) {
	v, ok := e.store[name]
	if !ok && e.outer != nil {
		v, ok = e.outer.GetVar(name)
	}
	return v.value, ok
}

// GetVar :
func (e *Environment) GetVar(name string) (Variable, bool) {
	v, ok := e.store[name]
	if !ok && e.outer != nil {
		v, ok = e.outer.GetVar(name)
	}
	return v, ok
}

// Set : set a variable inside the current scope.
func (e *Environment) Set(name string, val Object, mut bool) Object {
	e.store[name] = Variable{value: val, mutable: mut}
	return val
}

// need to think about the enforcement of constants.
func (e *Environment) existsAndMutable(name string) (bool, bool) {
	v, ok := e.store[name]
	if !ok && e.outer != nil {
		v, ok = e.outer.GetVar(name)
	}
	if !ok {
		return ok, false
	}
	return ok, v.mutable
}
