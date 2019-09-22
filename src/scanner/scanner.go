/*Package scanner implements Dito's Lexical Scanner.
Scanner scans file for tokens and passes them to parser.
This is implemented by the NextToken function, which the
parser will call till it reaches an token.EOF.
*/
package scanner

// TODO: fix issues with column positions, tracebacks etc.
// TODO: the line handling is pretty crap. trailing comments don't work
// 		 and the general setup just doesn't feel right.

import (
	"bytes"
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

// Init return an intialized lexical scanner. Values not
// initalized are set to go's default type values :
// 0 for all integers, false for bool.
func Init(input string) *Scanner {
	input += " " // add a buffer space at the end.
	s := &Scanner{input: input}
	s.advance()
	return s
}

// NextToken : Returns the next token encountered by the lexical scanner.
func (s *Scanner) NextToken() (tok token.Token, literal string, line int) {

	// ------------------------------------------------------------------------
	// * * * TODO * * *
	// The way new lines are handled is kinda rubbish. Appart from having to
	// the odd awkward check in the parser, the way it handles comments stops
	// trailing comments.

	if s.newline && s.char == '\n' {
		s.newline = false
		s.advanceLine()
		return token.NEWLINE, token.NEWLINE.String(), s.lineno - 1
	}

	s.newline = true

	s.skipWhitespace()

	// ------------------------------------------------------------------------

	// Run through all the possible operators and delimeters that are
	// included in dito's grammar, if not default to check for
	// identifers and numbers. If we still haven't found anything
	// set tok to token.ILLEGAL.
	switch s.char {
	case '=': // = ==
		tok = s.switch2(token.ASSIGN, '=', token.EQUALS)
	case '%': // % %=
		tok = s.switch2(token.MOD, '=', token.MODEQUAL)
	case '!': // ! !=
		tok = s.switch2(token.NOT, '=', token.NEQUALS)
	case '*': // * **
		tok = s.switch3(token.MUL, '*', token.POW, '=', token.MULEQUAL)
	case '>': // > >= >>
		tok = s.switch3(token.GTHAN, '=', token.GEQUALS, '>', token.RSHIFT)
	case '<': // < <= <<
		tok = s.switch3(token.LTHAN, '=', token.LEQUALS, '<', token.LSHIFT)
	case '/': // / /= //
		tok = s.switch3(token.DIV, '=', token.DIVEQUAL, '/', token.IDIV)
	case '-': // - -= ->
		tok = s.switch3(token.SUB, '=', token.SUBEQUAL, '>', token.RARROW)
	case '+': // + ++ +=
		tok = s.switch3(token.ADD, '+', token.CAT, '=', token.ADDEQUAL)
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
	case '|':
		tok = token.BITOR
	case '&':
		tok = token.BITAND
	case '^':
		tok = token.BITXOR
	case ':':
		tok = token.COLON
	case '"':
		return s.readString()
	case '.':
		if isDigit(s.peek()) {
			return s.readNumber()
		}
		tok = token.ILLEGAL
	case 0:
		// token.EOF represents end of input. the scanners caller should
		// check for this to find out when to stop iterating.
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
	var b bytes.Buffer
	for {
		s.advance()
		if s.char == '"' {
			break
		}
		if s.char == '\n' || s.char == 0 {
			return token.UNEXPECTEDEOF, b.String(), s.lineno
		}
		if s.char == '\\' {
			b.WriteByte(s.escapeString())
			continue
		}
		b.WriteByte(s.char)
	}
	s.advance()
	return token.STRING, b.String(), s.lineno
}

func (s *Scanner) escapeString() byte {
	s.advance()
	var val byte
	switch s.char {
	case 'n':
		val = '\n'
	case '\\':
		val = '\\'
	case 'r':
		val = '\r'
	case 't':
		val = '\t'
	case 'f':
		val = '\f'
	case 'b':
		val = '\b'
	case '"':
		val = '"'
	case '\'':
		val = '\''
	default:
		val = '\\'
	}
	return val
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
// with support for hex, scientific notation, and decimal.
// All hex digits are integers and all exponated digits are floats.
func (s *Scanner) readNumber() (token.Token, string, int) {
	start := s.pos
	// step 1: loop though digits until we read the end of 0-9.
	for isDigit(s.char) {
		s.advance()
	}
	// Step 2: Once we have the significand, now find the numbers type.
	// The rule is if it has a decimal point or uses scientific notation
	// its a float. If it is hex or just a plain integer return a int.
	switch {
	// all hexadecimals start strictly with a 0x or a 0X.
	case (s.char == 'x' || s.char == 'X') && s.input[start:s.pos] == "0":
		// Hexadecimal e.g. 0xffaf, 0X0032f, ...
		s.advance()
		for isHex(s.char) {
			s.advance()
		}
		return token.INT, s.input[start:s.pos], s.lineno
	case s.char == '.':
		// Mantissa e.g. 0.321, 312.123, ...
		s.advance()
		for isDigit(s.char) {
			s.advance()
		}
		if s.char != 'e' && s.char != 'E' {
			return token.FLOAT, s.input[start:s.pos], s.lineno
		}
		// if we reach here we go to the next case.
		fallthrough
	case s.char == 'e' || s.char == 'E':
		// Exponent e.g. 10e2, 8E-2, 8.23e10, ...
		s.advance()
		if s.char == '+' || s.char == '-' {
			s.advance()
		}
		s.advance()
		for isDigit(s.char) {
			s.advance()
		}
		return token.FLOAT, s.input[start:s.pos], s.lineno
	default:
		return token.INT, s.input[start:s.pos], s.lineno
	}
}

// skipWhitespace all skips whitespace and comments whilst keeping track
// of line numbers.
func (s *Scanner) skipWhitespace() {
	for {
		switch s.char {
		case ' ', '\r', '\t':
			s.advance()
		case '\n':
			s.advanceLine()
		case '#':
			for s.char != '\n' && s.char != 0 {
				s.advance()
			}
		default:
			return
		}
	}
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
