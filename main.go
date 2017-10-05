package main

import (
	"dito/evaluator"
	"dito/lexer"
	"dito/object"
	"dito/parser"
	"dito/repl"
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
	l := lexer.New(file)
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		repl.PrintParserErrors(out, p.Errors())
		return
	}
	evaluated := evaluator.Eval(program, env)
	if evaluated != nil {
		io.WriteString(out, evaluated.Inspect())
		io.WriteString(out, "\n")
	}
}

func welcomeMsg(username, quit string) {
	fmt.Printf("Hello %s! Welcome to the Dito programing language\n", username)
	fmt.Printf("Feel free to type commands. Type '%s' to quit.\n", quit)
}
