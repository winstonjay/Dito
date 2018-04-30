package object

// Object : defines the interface for the objects used in the dito programming
// language.
type Object interface {
	// Type : type which is used internaly and avalible to the user through the
	// builtin 'type' function at runtime.
	Type() TypeFlag
	// Inspect : returns the value of the object. Is used to display values to
	// the user.
	Inspect() string

	// ConvertType : Convert an objects type into another type. The defualt
	// behaviour for inconpaterble conversions is to return an error object.
	ConvertType(which TypeFlag) Object
}

// Iterable : the requirements needed for a Object to be an Iterable.
// eg: Array's, String's.
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

	// - Possible future candidates:

	// Merge : merege/concatenate two items together.
	// Merge(Object) Object
	// Contains : is item in the iterable. would use 'in' operator.
	// Contains(Object) Object
}

// // IterObject : key value pair for iterating over a Iterable
// type IterObject struct {
// 	key Object
// 	val Object
// }

// Numeric : the requirements needed for a Object to be an Numeric.
type Numeric interface {
	Object
	// Abs : return the absolute value of an number
	Abs() Object
}

// TypeFlag : type flag for what type of Dito object it is.
type TypeFlag int

// Define the strings used availible to the user to describe objects. Values
// here will be returned when an objects type method is called.
const (
	CharType TypeFlag = iota
	IntType
	FloatType
	BoolType
	StringType
	ArrayType
	NoneType
	ErrorType
	ReturnType
	LambdaType
	BultinType
	// FunctionType = "Function" // not implemented
)

func (t TypeFlag) String() string { return typeName[t] }

var typeName = [...]string{
	"Char", "Int", "Float", "Bool", "String", "Array", "None", "None", "Return", "Lambda", "Builtin",
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
