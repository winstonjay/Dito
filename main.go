// Dito is a Toy interpreted programming language for fun.
// : )
package main

import (
	"dito/src/evaluator"
	"dito/src/lexer"
	"dito/src/object"
	"dito/src/parser"
	"dito/src/repl"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
)

func main() {
	args := os.Args[1:] // args without program.
	// https://gobyexample.com/command-line-flags
	if len(args) > 0 {
		filename := args[0]
		file, err := ioutil.ReadFile("testdata/" + filename)
		if err != nil {
			panic(err)
		}
		execFile(string(file), os.Stdout)
		return
	}
	welcomeMsg(repl.QUIT)
	repl.Start(os.Stdin, os.Stdout)
}

func execFile(file string, out io.Writer) {
	env := object.NewEnvironment()
	l := lexer.Init(file)
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		p.PrintParseErrors(out, p.Errors())
		return
	}
	evaluated := evaluator.Eval(program, env)
	if evaluated != nil {
		io.WriteString(out, evaluated.Inspect())
		io.WriteString(out, "\n")
	}
}

func welcomeMsg(quit string) {
	fmt.Printf("Dito Interactive Shell V0.01 on %s\n", runtime.GOOS)
	fmt.Printf("Enter '%s' to quit. Help is coming soon...\n", quit)
}
