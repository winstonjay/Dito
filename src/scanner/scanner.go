// Package scanner implements Dito's Lexical Scanner.
// Scanner scans file for tokens and passes them to parser.
// This is implemented by the NextToken function, which the
// parser will call till it reaches an token.EOF.
package scanner

// TODO: fix issues with column positions, tracebacks etc.

import (
	"dito/src/token"
)

// Scanner implements the methods needed to scan a program.
type Scanner struct {
	input   string // Fixed state from new.
	char    byte   // current char under examination
	pos     int    // current char in input
	peekPos int    // current next position (index of peek char)
	linePos int    // index of the start of the current line.
	lineno  int    // current line under examination
	column  int    // current column position.
	newline bool   // if or not to collect newline tokens.
}

// Init return an intialised lexical scanner. Values not
// initalised are set to go's default type values :
// 0 for all integers, false for bool.
func Init(input string) *Scanner {
	input += " " // add a buffer space at the end.
	s := &Scanner{input: input}
	s.advance()
	return s
}

// NextToken : Returns the next token encountered by the lexical scanner.
func (s *Scanner) NextToken() (tok token.Token, literal string, line int) {

	// first see if we find a newline token return that and set the newline
	// collection to false. This is so that at we only collect one newline token
	// per space between non-whitespace chars. The rest of the newlines are not
	// tokenised but the line postion is incremented.
	if s.newline && (s.char == '\n' || s.char == '\r') {
		s.newline = false
		s.advanceLine()
		return token.NEWLINE, token.NEWLINE.String(), s.lineno - 1
	}

	// reset newline collection so we collect the first trailing
	// newline the next time this function is called.
	s.newline = true

	// Make sure all comments and spaces are skipped.
	s.skipWhitespace()
	for s.char == '#' {
		s.skipComment()
		s.skipWhitespace()
	}
	// Run through all the possible operators and delimeters that are
	// included in dito's grammar, if not default to check for
	// identifers and numbers. If we still havent found anything
	// set tok to token.ILLEGAL.
	switch s.char {
	case '=': // = ==
		tok = s.switch2(token.ASSIGN, '=', token.EQUALS)
	case '*': // * **
		tok = s.switch3(token.MUL, '*', token.POW, '=', token.MULEQUAL)
	case '!': // ! !=
		tok = s.switch2(token.NOT, '=', token.NEQUALS)
	case '>': // > >= >>
		tok = s.switch3(token.GTHAN, '=', token.LEQUALS, '>', token.SHIFTR)
	case '<': // < <= <<
		tok = s.switch3(token.LTHAN, '=', token.LEQUALS, '<', token.SHIFTL)
	case '/': // / //
		tok = s.switch3(token.DIV, '=', token.DIVEQUAL, '/', token.IDIV)
	case '-': // - -= ->
		tok = s.switch3(token.SUB, '=', token.SUBEQUAL, '>', token.RARROW)
	case '+': // + +=
		tok = s.switch2(token.ADD, '=', token.ADDEQUAL)
	case '%': // % %=
		tok = s.switch2(token.MOD, '=', token.MODEQUAL)
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
	case '.':
		if isDigit(s.peek()) {
			return s.readNumber()
		}
		tok = token.ILLEGAL
	case 0:
		// token.EOF represents end of input.
		// the scanners caller should check for
		// this to find out when to stop iterating.
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

// switch 2 checks between 2 possible alternatives give a current token and a
// the peek token returning the correct combination of chars.
func (s *Scanner) switch2(curr token.Token, expected byte, alt token.Token) token.Token {
	if s.peek() == expected {
		s.advance()
		return alt
	}
	return curr
}

// switch 3 checks between 3 possible alternatives give a
// current token and a the peek token.
func (s *Scanner) switch3(
	curr token.Token,
	expected1 byte, alt1 token.Token,
	expected2 byte, alt2 token.Token,
) token.Token {
	switch s.peek() {
	case expected1:
		s.advance()
		return alt1
	case expected2:
		s.advance()
		return alt2
	default:
		return curr
	}
}

// readString reads until it sees a double quote or 0 (EOF).
// Strings can only be created with double quotes.
// They can be multi-line but this should probally be
// treated as an error in the future.
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

// readIdentifer reads while alphanumeric chars are under inspection.
// it then determines if the collecting token literal is a keyword or
// a user defined identifier, returning the type accordingly.
func (s *Scanner) readIdentifer() (token.Token, string, int) {
	start := s.pos
	for isLetter(s.char) || isDigit(s.char) {
		s.advance()
	}
	literal := s.input[start:s.pos]
	return token.LookUpIDVal(literal), literal, s.lineno
}

// readNumber returns either an integer or a float type token
// with support for hex, e notation, and decimal.
// All hex digits are integers and all exponated digits are floats.
func (s *Scanner) readNumber() (token.Token, string, int) {
	start := s.pos
	// loop though digits until we read the end of 0-9.
	for isDigit(s.char) {
		s.advance()
	}
	// Once we have the significand, now find the numbers type.
	switch {
	// all hexadecimals start strictly with a 0x or a 0X.
	case (s.char == 'x' || s.char == 'X') && s.input[start:s.pos] == "0":
		goto Hexadecimal
	case s.char == '.':
		goto Mantissa
	case s.char == 'e' || s.char == 'E':
		goto Exponent
	default:
		return token.INT, s.input[start:s.pos], s.lineno
	}

Hexadecimal:
	// 0xffaf, 0X0032f, ...
	s.advance()
	for isHex(s.char) {
		s.advance()
	}
	return token.INT, s.input[start:s.pos], s.lineno

Mantissa:
	// 0.321, 312.123, ...
	s.advance()
	for isDigit(s.char) {
		s.advance()
	}
	if s.char == 'e' || s.char == 'E' {
		goto Exponent
	}
	return token.FLOAT, s.input[start:s.pos], s.lineno

Exponent:
	// 10e2, 8E-2, 8.23e10, ...
	s.advance()
	if s.char == '+' || s.char == '-' {
		s.advance()
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
	s.char = s.input[s.peekPos]
	s.pos = s.peekPos
	s.peekPos++
	s.column++
}

func (s *Scanner) advanceLine() {
	s.advance()
	s.linePos = s.pos
	s.column = 0
	s.lineno++
}

func (s *Scanner) peek() byte {
	if s.peekPos >= len(s.input) {
		return 0
	}
	return s.input[s.peekPos]
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
	if s.char != 0 {
		s.advanceLine()
	}
}

func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}

func isLetter(char byte) bool {
	return ('a' <= char && char <= 'z' ||
		'A' <= char && char <= 'Z' || char == '_')
}

func isSpace(char byte) bool {
	return (char == ' ' || char == '\t' ||
		char == '\n' || char == '\r')
}

func isHex(char byte) bool {
	return (isDigit(char) || 'a' <= char && char <= 'z' ||
		'A' <= char && char <= 'Z')
}
