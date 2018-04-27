package object

// // Array : array object.
// type Array struct {
// 	Elements []Object
// 	Len      int64
// }

// func (a *Array) Type() string { return ArrayType }

// func (a *Array) Inspect() string {
// 	var out bytes.Buffer
// 	out.WriteString("[")
// 	for i, el := range a.Elements {
// 		out.WriteString(el.Inspect())
// 		if i < len(a.Elements)-1 {
// 			out.WriteString(", ")
// 		}
// 	}
// 	out.WriteString("]")
// 	return out.String()
// }

// // NewArray :
// func NewArray(elements []Object, length int64) *Array {
// 	if length == -1 {
// 		length = int64(len(elements))
// 	}
// 	return &Array{Elements: elements, Len: length}
// }
