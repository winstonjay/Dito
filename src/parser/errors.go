package parser

import (
	"dito/src/token"
	"fmt"
	"io"
	"strings"
)

// TODO: fix issues with column positions, tracebacks etc
// TODO: Work on traceback style.
// Example error output or the aim for what the output aspires to be.
// print all error messages then the last traceback.
//     Parse Error: Expected next token is ':='. got '9' instead.
//     Parse Error: Expected next token is '('. got '!' instead.
//     Parse Error: Expected next token is ';'. got '=' instead.
//     Last traceback @ line 3, col 20:
// 	       x := pasta * 2 / 3 =
//                            ^ your problem right there.

// ParseError : store error message
type ParseError struct {
	message   string // what's the problem?
	column    int    // what line it is on?
	lineno    int    // where does the problem start on this line?
	lineTrace string // what does this line look like?
}

func (pe *ParseError) String() string {
	return fmt.Sprintf(
		"Traceback line %d column %d:\n%*s%s\n%*s%s\n%s\n",
		pe.lineno, pe.column, 4, " ", pe.lineTrace, 4+pe.column, " ", "^ Is your problem here?",
		pe.message)
}

func (p *Parser) newError(message string) *ParseError {
	linecontent := strings.TrimRight(p.scanner.TraceLine(), "\n\t ")
	return &ParseError{
		message:   message,
		lineno:    p.currentLine,
		column:    len(linecontent) - len(p.currentLiteral),
		lineTrace: linecontent,
	}
}

// PrintParseErrors : output all parse errors.
func (p *Parser) PrintParseErrors(out io.Writer, errors []*ParseError) {
	io.WriteString(out, "PARSE ERROR:\n")
	for _, msg := range p.Errors() {
		io.WriteString(out, msg.String())
	}
}

// Errors : return a array o all acculumated errors found when parsing.
func (p *Parser) Errors() []*ParseError {
	return p.errors
}

func (p *Parser) genericError(t string, e error) {
	msg := fmt.Sprintf("Could not parse %q as %s: %v", p.currentLiteral, t, e)
	p.errors = append(p.errors, p.newError(msg))
}

func (p *Parser) peekError(t token.Token) {
	msg := fmt.Sprintf("Expected next token is '%s'. got '%s' instead",
		t, p.peekLiteral)
	p.errors = append(p.errors, p.newError(msg))
}

func (p *Parser) noParseFnError(t token.Token) {
	msg := fmt.Sprintf("No parse function found for '%s'", t)
	p.errors = append(p.errors, p.newError(msg))
}
