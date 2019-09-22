package object

// Object : defines the interface for the objects used in the dito programming
// language.
type Object interface {
	// Type : type which is used internally and available to the user through the
	// builtin 'type' function at runtime.
	Type() TypeFlag
	// Inspect : returns the value of the object. Is used to display values to
	// the user.
	Inspect() string
	// ConvertType : Convert an objects type into another type. The default
	// behavior for incompatible conversions is to return an error object.
	ConvertType(which TypeFlag) Object
}

// TODO think about Numeric.

// Numeric : the requirements needed for a Object to be an Numeric.
type Numeric interface {
	Object
	// Abs : return the absolute value of an number
	Abs() Object
}

// Iterable : the requirements needed for a Object to be an Iterable.
// eg: Array, String, Dict.
type Iterable interface {
	Object
	// Length : return the number of items in the iterable
	Length() Object
	// IterItems : return a iterable channel to loop through the items in
	// order.
	Iter() <-chan Object
	// GetItem : get item at location of the provided key.
	GetItem(Object) Object
	// SetItem : set item at location of the provided key.
	SetItem(Object, Object) Object
	// Concat : concat two iterables together.
	Concat(Object) Object
	// Contains : is item in the iterable. would use 'in' operator.
	Contains(Object) Object
	// Slice : return a slice of an iterables elements
	Slice(Object, Object) Object
}

// // IterObject : key value pair for iterating over a Iterable
// type IterObject struct {
// 	key Object
// 	val Object
// }

// Hashable : the requirements to enable the object to be used as a hash key.
type Hashable interface {
	Object
	// Hash : returns the Hashkey for the objects value
	Hash() HashKey
}

// HashKey :
type HashKey struct {
	Type  TypeFlag
	Value uint64
}

// ''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''
// Singleton objects :  Only one instance of these needs to be created.
var (
	TRUE  = &Bool{Value: true}
	FALSE = &Bool{Value: false}
	NONE  = &None{}
)

// ''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''
// Internal objects : users should not really see these

// ReturnValue : Packages other objects to determine
// the end objects of programs.
type ReturnValue struct{ Value Object }

// Type : return objects type as a string
func (rv *ReturnValue) Type() TypeFlag { return ReturnType }

// Inspect : return a string representation of the objects value.
func (rv *ReturnValue) Inspect() string { return rv.Value.Inspect() }

// ConvertType : return the conversion into the specified type
func (rv *ReturnValue) ConvertType(which TypeFlag) Object {
	return NewError("Argument to %s not supported, got %s", rv.Type(), which)
}
