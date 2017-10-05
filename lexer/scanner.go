package lexer

import (
	"dito/token"
)

/*
==================================
	Dito's Lexical Scanner.
==================================




*/

// Scanner : Lexical scanner analises a given program string.
// It must be created/initalised by using the 'New' function
//
type Scanner struct {
	// Fixed state from new.
	input string

	// Scanning state.
	char    byte // current char under examination
	pos     int  // current char in input
	peekPos int  // current next position (index of peek char)
	linePos int  // index of the start of the current line.
	Lineno  int  // current line under examination
	Column  int  // current Column position.
}

// New : return an intialised lexical scanner. values not
// initalised are set to go's default type values. so 0 for
// all integers.
func New(input string) *Scanner {
	input += " " // add a buffer space at the end.
	s := &Scanner{input: input}
	s.advance()
	return s
}

// NextToken : Returns the next token encountered by the lexical scanner.
func (s *Scanner) NextToken() (tok token.Token, literal string) {

	// First make sure all comments and spaces are skipped.
	s.skipWhitespace()
	for s.char == '#' {
		s.skipComment()
		s.skipWhitespace()
	}
	// Run through all the possible operators and delimeters that are
	// included in dito's grammar.
	switch s.char {
	case ':': // :, :=
		tok = s.switch2(token.COLON, token.REASSIGN, token.NEWASSIGN)
	case '=': // =, ==
		tok = s.switch2(token.REASSIGN, token.REASSIGN, token.EQUALS)
	case '*': // *, **
		tok = s.switch2(token.MUL, token.MUL, token.POW)
	case '!': // !, !=
		tok = s.switch2(token.NOT, token.REASSIGN, token.NEQUALS)
	case '>': // >, >=
		tok = s.switch2(token.GTHAN, token.REASSIGN, token.GEQUALS)
	case '<': // <, <=
		tok = s.switch2(token.LTHAN, token.REASSIGN, token.LEQUALS)
	case '/': // /, //
		tok = s.switch2(token.DIV, token.DIV, token.IDIV)
	case '+':
		tok = token.ADD
	case '-':
		tok = token.SUB
	case '/':
		tok = token.DIV
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
	case '[':
		tok = token.LBRACKET
	case ']':
		tok = token.RBRACKET
	case '"', '\'':
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
		// handed on to the parser and will show up as a no parse function
		// found error.
		tok = token.ILLEGAL
	}
	s.advance() // Always advance.
	return tok, tok.String()
}

func (s *Scanner) switch2(current, expected, alt token.Token) token.Token {
	if string(s.peek()) == expected.String() {
		s.advance()
		return alt
	}
	return current
}

func (s *Scanner) readString() (token.Token, string) {
	start := s.pos + 1
	for {
		s.advance()
		if s.char == '"' || s.char == 0 {
			break
		}
	}
	literal := s.input[start:s.pos]
	s.advance()
	return token.STRING, literal
}

func (s *Scanner) readIdentifer() (token.Token, string) {
	start := s.pos
	for isLetter(s.char) || isDigit(s.char) {
		s.advance()
	}
	literal := s.input[start:s.pos]
	return token.LookUpIDVal(literal), literal
}

func (s *Scanner) readNumber() (token.Token, string) {
	start := s.pos
	for isDigit(s.char) {
		s.advance()
	}
	if s.char != '.' {
		return token.INT, s.input[start:s.pos]
	}
	s.advance()
	for isDigit(s.char) {
		s.advance()
	}
	return token.FLOAT, s.input[start:s.pos]
}

// advance : progresses the scanners current position, peek position,
// current char being examined and the column position. If we reach the
// end of the input file this method will return 0 which is designated
// as our EOF indicator. This interacts with the NextToken so that once
// this end of input is reached it forever more just return EOF tokens.
// at which point you should cut it off.
func (s *Scanner) advance() {
	if s.peekPos >= len(s.input) {
		s.char = 0
		return
	}
	s.char = s.input[s.peekPos]
	s.pos = s.peekPos
	s.peekPos++
	s.Column++
}

func (s *Scanner) advanceLine() {
	s.advance()
	s.linePos = s.pos
	s.Column = 0
	s.Lineno++
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
	s.advanceLine()
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
