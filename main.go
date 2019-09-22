// Dito is a Toy interpreted programming language for fun.
// : )
package main

import (
	"dito/src/eval"
	"dito/src/object"
	"dito/src/parser"
	"dito/src/repl"
	"dito/src/scanner"
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
		filepath := args[0]
		file, err := ioutil.ReadFile(filepath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		}
		execFile(string(file), os.Stdout)
		return
	}
	welcomeMsg(repl.QUIT)
	repl.Start(os.Stdin, os.Stdout)
}

func execFile(file string, out io.Writer) {
	env := object.InitialEnvironment()
	l := scanner.Init(file + "\n\nmain()")
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		p.PrintParseErrors(out, p.Errors())
		return
	}
	evaluated := eval.Eval(program, env)
	if evaluated != nil && evaluated != object.NONE {
		io.WriteString(out, evaluated.Inspect())
		io.WriteString(out, "\n")
	}
	// quick hack for automatic main calls.

}

func welcomeMsg(quit string) {
	fmt.Printf("\033[33mDito Interactive Shell V0.01\033[m on %s\n", runtime.GOOS)
	fmt.Printf("Enter '%s' to quit. Help is coming soon...\n", quit)
}
