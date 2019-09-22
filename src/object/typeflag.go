package object

// TypeFlag : type flag for what type of Dito object it is.
type TypeFlag int

// Define the strings used available to the user to describe objects. Values
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
	FunctionType
	DictType
	FileType
)

func (t TypeFlag) String() string { return typeName[t] }

var typeName = [...]string{
	"Char", "Int", "Float", "Bool", "String", "Array",
	"None", "Error", "Return", "Lambda", "Builtin", "Function", "Dict", "File",
}
