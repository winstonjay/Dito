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
	"os/user"
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
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	welcomeMsg(user.Username, repl.QUIT)
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
	evaluator.Eval(program, env)
}

func welcomeMsg(username, quit string) {
	fmt.Printf("Hello %s! Welcome to the Dito programing language\n", username)
	fmt.Printf("Feel free to type commands. Type '%s' to quit.\n", quit)
}
