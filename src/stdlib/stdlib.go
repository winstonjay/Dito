package stdlib

import (
	"dito/src/object"
	"io"
	"os"
)

// Print outputs to the terminal a varible amount of args
// it is ended by a '\n' character.
func Print(args ...object.Object) object.Object {
	for _, arg := range args {
		io.WriteString(os.Stdout, arg.Inspect())
	}
	io.WriteString(os.Stdout, "\n")
	return nil
}
