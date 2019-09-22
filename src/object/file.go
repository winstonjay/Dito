package object

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"syscall"
)

// File : wrapper around goes file system
type File struct {
	file *os.File
}

// Type : return objects type as a TypeFlag
func (f *File) Type() TypeFlag { return FileType }

// Inspect : return a string representation of the objects value.
func (f *File) Inspect() string { return fmt.Sprintf("<File: %s>", f.file.Name()) }

// ConvertType : return the conversion into the specified type
func (f *File) ConvertType(which TypeFlag) Object {
	return NewError("Argument to %s not supported, got %s", which, f.Type())
}

// Write :
func (f *File) Write(str string) Object {
	_, err := f.file.WriteString(str)
	if err != nil {
		return NewError("Could not write file %s", f.Inspect())
	}
	return NONE
}

// Read :
func (f *File) Read() Object {
	reader := bufio.NewReader(f.file)
	fileBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return NewError("Could not read file %s", f.Inspect())
	}
	return NewString(string(fileBytes))
}

// Close :
func (f *File) Close() Object {
	if err := f.file.Close(); err != nil {
		return NewError("Could not close file %s", f.Inspect())
	}
	return NONE
}

// FileHandler :
type FileHandler struct{}

// Open : return a new readable file
func (fh *FileHandler) Open(path string) Object {
	fp, err := os.Open(path)
	if err != nil {
		return NewError("Could not open file '%s' (%s)", path, err)
	}
	return &File{fp}
}

// Create : return a new writable file
func (fh *FileHandler) Create(path string) Object {
	fp, err := os.Create(path)
	if err != nil {
		return NewError("Could not open file '%s' (%s)", path, err)
	}
	return &File{fp}
}

// OS VARS
var (
	FILE   = &FileHandler{}
	STDIN  = &File{os.NewFile(uintptr(syscall.Stdin), "/dev/stdin")}
	STDOUT = &File{os.NewFile(uintptr(syscall.Stdout), "/dev/stdout")}
	STDERR = &File{os.NewFile(uintptr(syscall.Stderr), "/dev/stderr")}
)
