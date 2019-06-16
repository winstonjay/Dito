package imports

import (
	"crypto/sha1"
	"dito/src/ast"
	"dito/src/eval"
	"dito/src/object"
	"dito/src/parser"
	"dito/src/scanner"
	"fmt"
	"io/ioutil"
	"os"
)

// Importer :
type Importer struct {
	env    *object.Environment
	hashes map[string]bool
}

func (i *Importer) loadFile(name string) (string, error) {
	filepath := "lib/" + name + ".dito"
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	hash := fmt.Sprintf("%x", sha1.Sum(file))
	if i.hashes[hash] != false {
		return "", nil
	}
	i.hashes[hash] = true
	return string(file), nil
}

func evalImportStatement(node *ast.ImportStatement, env *object.Environment) object.Object {
	filename := node.Value
	filepath := "lib/" + filename + ".dito"
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return object.NewError("Import File %s could not be opened.", filepath)
	}
	il := scanner.Init(string(file))
	ip := parser.New(il)
	iprogram := ip.ParseProgram()

	if len(ip.Errors()) != 0 {
		ip.PrintParseErrors(os.Stderr, ip.Errors())
		return object.NewError("Could not import file due to parse errors.")
	}
	eval.Eval(iprogram, env)
	// return object.NONE
	return nil
}
