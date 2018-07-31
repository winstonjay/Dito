package repl

import (
	"bufio"
	"dito/src/eval"
	"dito/src/object"
	"dito/src/parser"
	"dito/src/scanner"
	"dito/src/token"
	"fmt"
	"io"
)

// command prompt constants
const (
	PROMPT = "dito: "
	QUIT   = "QQ"
)

// Start : run repl for the dito interpeter.
func Start(in io.Reader, out io.Writer) {
	b := bufio.NewScanner(in)
	env := object.NewEnvironment()
	for {
		fmt.Printf(PROMPT)
		scanned := b.Scan()
		if !scanned {
			return
		}
		// semi colons are being put here just for now in the repl.
		line := b.Text()
		if line == QUIT {
			return
		}
		l := scanner.Init(line + "\n")
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			p.PrintParseErrors(out, p.Errors())
			continue
		}
		evaluated := eval.Eval(program, env)
		if evaluated != nil && evaluated != object.NONE {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

// For running inspections of the lexical scanner.
func printTokens(l *scanner.Scanner) {
	toke, literal, _ := l.NextToken()
	for toke != token.EOF {
		fmt.Printf("Token(type='%s', literal='%s')\n", toke.String(), literal)
		toke, literal, _ = l.NextToken()
	}
}
