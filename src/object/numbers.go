package object

import (
	"fmt"
	"strconv"
)

// Integer : builtin integer type.
// -9223372036854775807 and 9223372036854775807
type Integer struct{ Value int64 }

// Type :
func (i *Integer) Type() string    { return IntergerObj }
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

func NewDitoInteger(value int64) *Integer { return &Integer{Value: value} }

// Float : builtin float type.
type Float struct{ Value float64 }

// Type :
func (f *Float) Type() string    { return FloatObj }
func (f *Float) Inspect() string { return strconv.FormatFloat(f.Value, 'f', -1, 64) }

func NewDitoFloat(value float64) *Float { return &Float{Value: value} }
