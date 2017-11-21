package repl

import (
	"bufio"
	"dito/evaluator"
	"dito/lexer"
	"dito/object"
	"dito/parser"
	"dito/token"
	"fmt"
	"io"
)

// command prompt constants
const (
	PROMPT = "(dito)> "
	QUIT   = "QQ"
)

// Start : run repl for the dito interpeter.
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		// semi colons are being put here just for now in the repl.
		line := scanner.Text() + ";"
		if line == QUIT {
			return
		}
		l := lexer.Init(line)
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			PrintParserErrors(out, p.Errors())
			continue
		}
		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

// PrintParserErrors :
func PrintParserErrors(out io.Writer, errors []*parser.ParseError) {
	io.WriteString(out, "PARSE ERROR:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg.String()+"\n")
	}
}

// For running inspections of the lexical scanner.
func printTokens(l *lexer.Scanner) {
	toke, literal, _ := l.NextToken()
	for toke != token.EOF {
		fmt.Printf("Token(type='%s', literal='%s')\n", toke.String(), literal)
		toke, literal, _ = l.NextToken()
	}
}
