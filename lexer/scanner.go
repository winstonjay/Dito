// Package lexer implements Dito's Lexical Scanner.
// Scanner scans file for tokens and passes them to parser.
// This is implemented by the NextToken function, which the
// parser will call till it reaches an token.EOF.
package lexer

/*
TODO:
	* Implement pattern for reading newlines without bugs.
*/

import (
	"dito/token"
	"fmt"
)

// Scanner is Lexical scanner analises a given program string.
// It must be created/initalised by using the 'New' function.
type Scanner struct {
	// Fixed state from new.
	input        string
	char         rune // current char under examination
	pos          int  // current char in input
	peekPos      int  // current next position (index of peek char)
	linePos      int  // index of the start of the current line.
	lineno       int  // current line under examination
	column       int  // current column position.
	collectspace bool
}

// Init return an intialised lexical scanner. values not
// initalised are set to go's default type values. so 0 for
// all integers.
func Init(input string) *Scanner {
	input += " " // add a buffer space at the end.
	s := &Scanner{input: input}
	s.advance()
	return s
}

// NextToken : Returns the next token encountered by the lexical scanner.
func (s *Scanner) NextToken() (tok token.Token, literal string, line int) {

	// Sorting the issue of new lines.
	// in repl there is a semi colon being added so we dont have to type it.
	// tests and other things have not been set for this.

	// if s.char == '\n' && !s.collectspace {
	// 	s.collectspace = true
	// 	tok = token.NEWLINE
	// }
	// First make sure all comments and spaces are skipped.
	s.skipWhitespace()
	for s.char == '#' {
		s.skipComment()
		s.skipWhitespace()
	}
	// Run through all the possible operators and delimeters that are
	// included in dito's grammar.
	switch s.char {
	// Tokens that can be one of 2 values.
	case ':': // ':', ':='
		tok = s.switch2(token.COLON, token.REASSIGN, token.NEWASSIGN)
	case '=': // '=', '=='
		tok = s.switch2(token.REASSIGN, token.REASSIGN, token.EQUALS)
	case '*': // '*', '**'
		tok = s.switch2(token.MUL, token.MUL, token.POW)
	case '!': // '!', '!='
		tok = s.switch2(token.NOT, token.REASSIGN, token.NEQUALS)
	case '>': // '>', '>='
		tok = s.switch2(token.GTHAN, token.REASSIGN, token.GEQUALS)
	case '<': // '<', '<='
		tok = s.switch2(token.LTHAN, token.REASSIGN, token.LEQUALS)
	case '/': // '/', '//'
		tok = s.switch2(token.DIV, token.DIV, token.IDIV)
	case '-': // '-', '->'
		tok = s.switch2(token.SUB, token.GTHAN, token.RARROW)
	case '+':
		tok = token.ADD
	case '%':
		tok = token.MOD
	case '(':
		tok = token.LPAREN
	case ')':
		tok = token.RPAREN
	case ';':
		tok = token.SEMI
	case ',':
		tok = token.COMMA
	case '{':
		tok = token.LBRACE
	case '}':
		tok = token.RBRACE
	case '[':
		tok = token.LBRACKET
	case ']':
		tok = token.RBRACKET
	case '"':
		return s.readString()
	case 0:
		tok = token.EOF
	default:
		if isDigit(s.char) {
			return s.readNumber()
		}
		if isLetter(s.char) {
			return s.readIdentifer()
		}
		// We shouldnt have got to this point irl. Error is currently
		// handed by the parser as a no parse function found error.
		tok = token.ILLEGAL
	}
	s.advance() // Always advance.
	return tok, tok.String(), s.lineno
}

// TraceLine : Returns last line up to current column.
// eg. at index 8 of "alpha := 100" we would get: 'alpha :=' <-.
func (s *Scanner) TraceLine() string {
	return s.input[s.linePos : s.linePos+s.column]
}

func (s *Scanner) switch2(current, expected, alt token.Token) token.Token {
	if string(s.peek()) == expected.String() {
		s.advance()
		return alt
	}
	return current
}

func (s *Scanner) readString() (token.Token, string, int) {
	start := s.pos + 1
	for {
		s.advance()
		if s.char == '"' || s.char == 0 {
			break
		}
	}
	literal := s.input[start:s.pos]
	s.advance()
	return token.STRING, literal, s.lineno
}

func (s *Scanner) readIdentifer() (token.Token, string, int) {
	start := s.pos
	for isLetter(s.char) || isDigit(s.char) {
		s.advance()
	}
	literal := s.input[start:s.pos]
	return token.LookUpIDVal(literal), literal, s.lineno
}

// readNumber : Return either an integer or a float.
// TODO: Add support for hex and mabye binary.
func (s *Scanner) readNumber() (token.Token, string, int) {
	start := s.pos
	for isDigit(s.char) {
		s.advance()
	}
	if s.char != '.' {
		return token.INT, s.input[start:s.pos], s.lineno
	}
	s.advance()
	for isDigit(s.char) {
		s.advance()
	}
	return token.FLOAT, s.input[start:s.pos], s.lineno
}

func (s *Scanner) advance() {
	if s.peekPos >= len(s.input) {
		s.char = 0
		return
	}
	s.char = rune(s.input[s.peekPos])
	s.pos = s.peekPos
	s.peekPos++
	s.column++
}

// TODO: currently the line position is passed on
// the parser and then differences in lines are used
// to determine the end of expressions and statements.
// this is not working right and in some cases semi
// colons need to be used to stop errors. main example
// is if statements geting confused with if experessions.
func (s *Scanner) advanceLine() {
	s.advance()
	s.linePos = s.pos
	s.column = 0
	s.lineno++
}

func (s *Scanner) peek() rune {
	if s.peekPos >= len(s.input) {
		return 0
	}
	return rune(s.input[s.peekPos])
}

func (s *Scanner) skipWhitespace() {
	for isSpace(s.char) {
		if s.char == '\n' || s.char == '\r' {
			s.advanceLine()
		} else {
			s.advance()
		}
	}
}

func (s *Scanner) skipComment() {
	for s.char != 0 && !(s.char == '\n' || s.char == '\r') {
		s.advance()
	}
	s.advanceLine()
}

func isDigit(char rune) bool {
	return '0' <= char && char <= '9'
}

func isLetter(char rune) bool {
	return ('a' <= char && char <= 'z' ||
		'A' <= char && char <= 'Z' || char == '_')
}

func isSpace(char rune) bool {
	return (char == ' ' || char == '\t' ||
		char == '\n' || char == '\r')
}

// PrintScan : print out the entire lexical analysis of an input
// in one go.
func (s *Scanner) printScan() {

	tok, literal, _ := s.NextToken()
	tokenCount := 0
	fmt.Printf("input:\n\n%s\n\n", s.input)
	fmt.Printf("| line | col  | Token        | Literal     |\n")
	fmt.Printf("-----------------------------------------\n")
	for tok != token.EOF {
		fmt.Printf("| %4d | %4d | %12s | %12s |\n",
			s.lineno+1, s.column-len(literal), tok.String(), literal)
		tokenCount++
		tok, literal, _ = s.NextToken()
	}
	fmt.Printf("\nTotal Tokens: %d, \n", tokenCount)
	fmt.Printf("Total Chars: %d, \n", s.pos)
	fmt.Printf("Total Lines: %d, \n", s.lineno+1)
}
