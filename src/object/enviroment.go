package object

// Varible : symbol table entry, keeps track whether varible is mutable or not.
type Varible struct {
	value   Object
	mutable bool
}

// IsMutable : can it be changed
func (v *Varible) IsMutable() bool {
	return v.mutable
}

// Unpack : can it be changed
func (v *Varible) Unpack() Object {
	return v.value
}

// Environment : Holds the enviroment varibles created by the user. Pretty much
// a symbol table.
type Environment struct {
	store map[string]Varible
	outer *Environment
}

// NewEnvironment : Define a new enviroment scope.
func NewEnvironment() *Environment {
	s := make(map[string]Varible)
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
	v, ok := e.store[name]
	if !ok && e.outer != nil {
		v, ok = e.outer.GetVar(name)
	}
	return v.value, ok
}

// GetVar :
func (e *Environment) GetVar(name string) (Varible, bool) {
	v, ok := e.store[name]
	if !ok && e.outer != nil {
		v, ok = e.outer.GetVar(name)
	}
	return v, ok
}

// Set : set a varible inside the current scope.
func (e *Environment) Set(name string, val Object, mut bool) Object {
	e.store[name] = Varible{value: val, mutable: mut}
	return val
}

// need to think about the enforment of constants.
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
