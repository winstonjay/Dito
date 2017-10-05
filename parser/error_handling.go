package parser

import (
	"dito/token"
	"fmt"
)

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

/*
Example error output or the aim for what the output should be.
print all error messages then the last traceback.

Parse Error: Expected next token is ':='. got '9' instead.
Parse Error: Expected next token is '('. got '!' instead.
Parse Error: Expected next token is ';'. got '=' instead.
Last traceback @ line 3, col 20:
	x := pasta * 2 / 3 =
						  ^ your problem right there.


how error s are handeled here:

file: 'file:///filname'
severity: 'Error'
message: 'cannot use p.assignmentStatement() (type *ast.AssignmentStatement) as
type ast.Statement in return argument:
	*ast.AssignmentStatement does not implement ast.Statement
	(wrong type for ast.tokenLiteral method)
		have ast.tokenLiteral() string
		want ast.tokenLiteral()'
at: '59,4'
source: ''
*/

// ParseError : store error message
type ParseError struct {
	message   string // what's the problem?
	column    int    // what line it is on?
	lineno    int    // where does the problem start on this line?
	lineTrace string // what does this line look like?
}

func (pe *ParseError) String() string {
	return pe.message
}

func (p *Parser) newError(message string) *ParseError {
	return &ParseError{
		message:   message,
		lineno:    p.scanner.Lineno,
		column:    p.scanner.Column - len(p.currentLiteral),
		lineTrace: p.scanner.TraceLine(),
	}
}
