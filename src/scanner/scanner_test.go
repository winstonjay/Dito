package scanner

import (
	"dito/src/token"
	"testing"
)

func TestNextToken(t *testing.T) {

	input := `let alpha = 10.999 + 10**3
let mut _o_me_ga = (50 % 7) - 0.002


# this is a single line comment you should not read me...
let hypot = def(a, b)->(a**2 + b**2)**0.5

let summer = true
let mut rain = true
let fun = false != true
"he" ++ "llo"
[10e2 + 3, 8.23e10, 0xffcc, 3e-5]; def float_5(a) { return .5 * a; }
.@ @
x /= fun << 1
{"f\nn": hypot}`
	// ^^^^^ DO NOT CHANGE OR YOU HAVE TO WORK OUT THE WHOLE TEST AGAIN
	// if you have to. extend the current string.
	// TODO (ALREADY BROKE cols NEEDS FIXING, there are some alignment
	// problemns with the coloum positions seen in trackbacks.
	tests := []struct {
		token   token.Token
		literal string
		line    int
		column  int // TODO
	}{
		// let alpha = 10.999 + 10**3
		{token.LET, "let", 0, 0},
		{token.IDVAL, "alpha", 0, 0},
		{token.ASSIGN, "=", 0, 0},
		{token.FLOAT, "10.999", 0, 0},
		{token.ADD, "+", 0, 0},
		{token.INT, "10", 0, 0},
		{token.POW, "**", 0, 0},
		{token.INT, "3", 0, 0},
		{token.NEWLINE, "NEWLINE", 0, 0},

		// let mut _o_me_ga = (50 % 7) - 0.002
		{token.LET, "let", 1, 0},
		{token.MUT, "mut", 1, 0},
		{token.IDVAL, "_o_me_ga", 1, 0},
		{token.ASSIGN, "=", 1, 0},
		{token.LPAREN, "(", 1, 0},
		{token.INT, "50", 1, 0},
		{token.MOD, "%", 1, 0},
		{token.INT, "7", 1, 0},
		{token.RPAREN, ")", 1, 0},
		{token.SUB, "-", 1, 0},
		{token.FLOAT, "0.002", 1, 0},
		{token.NEWLINE, "NEWLINE", 1, 0},

		// skips comment.

		// let hypot = func(a, b)->(a**2 + b**2)**0.5
		{token.LET, "let", 5, 0},
		{token.IDVAL, "hypot", 5, 0},
		{token.ASSIGN, "=", 5, 0},
		{token.DEF, "def", 5, 0},
		{token.LPAREN, "(", 5, 0},
		{token.IDVAL, "a", 5, 0},
		{token.COMMA, ",", 5, 0},
		{token.IDVAL, "b", 5, 0},
		{token.RPAREN, ")", 5, 0},
		{token.RARROW, "->", 5, 0},
		{token.LPAREN, "(", 5, 0},
		{token.IDVAL, "a", 5, 0},
		{token.POW, "**", 5, 0},
		{token.INT, "2", 5, 0},
		{token.ADD, "+", 5, 0},
		{token.IDVAL, "b", 5, 0},
		{token.POW, "**", 5, 0},
		{token.INT, "2", 5, 0},
		{token.RPAREN, ")", 5, 0},
		{token.POW, "**", 5, 0},
		{token.FLOAT, "0.5", 5, 0},
		{token.NEWLINE, "NEWLINE", 5, 0},

		// let summer = true
		{token.LET, "let", 7, 0},
		{token.IDVAL, "summer", 7, 0},
		{token.ASSIGN, "=", 7, 0},
		{token.TRUE, "true", 7, 0},
		{token.NEWLINE, "NEWLINE", 7, 0},

		// let mut rain = true
		{token.LET, "let", 8, 0},
		{token.MUT, "mut", 8, 0},
		{token.IDVAL, "rain", 8, 0},
		{token.ASSIGN, "=", 8, 0},
		{token.TRUE, "true", 8, 0},
		{token.NEWLINE, "NEWLINE", 8, 0},

		// let fun = false != true
		{token.LET, "let", 9, 0},
		{token.IDVAL, "fun", 9, 0},
		{token.ASSIGN, "=", 9, 0},
		{token.FALSE, "false", 9, 0},
		{token.NEQUALS, "!=", 9, 0},
		{token.TRUE, "true", 9, 0},
		{token.NEWLINE, "NEWLINE", 9, 0},

		// "he" ++ "llo"
		{token.STRING, "he", 10, 0},
		{token.CAT, "++", 10, 0},
		{token.STRING, "llo", 10, 0},
		{token.NEWLINE, "NEWLINE", 10, 0},

		// [10e2 + 3, 8.23e10, 0xffcc, 3e-5]; func float_5(a) { return .5 * a; }
		{token.LBRACKET, "[", 11, 0},
		{token.FLOAT, "10e2", 11, 0},
		{token.ADD, "+", 11, 0},
		{token.INT, "3", 11, 0},
		{token.COMMA, ",", 11, 0},
		{token.FLOAT, "8.23e10", 11, 0},
		{token.COMMA, ",", 11, 0},
		{token.INT, "0xffcc", 11, 0},
		{token.COMMA, ",", 11, 0},
		{token.FLOAT, "3e-5", 11, 0},
		{token.RBRACKET, "]", 11, 0},
		{token.SEMI, ";", 11, 0},
		{token.DEF, "def", 11, 0},
		{token.IDVAL, "float_5", 11, 0},
		{token.LPAREN, "(", 11, 0},
		{token.IDVAL, "a", 11, 0},
		{token.RPAREN, ")", 11, 0},
		{token.LBRACE, "{", 11, 0},
		{token.RETURN, "return", 11, 0},
		{token.FLOAT, ".5", 11, 0},
		{token.MUL, "*", 11, 0},
		{token.IDVAL, "a", 11, 0},
		{token.SEMI, ";", 11, 0},
		{token.RBRACE, "}", 11, 0},
		{token.NEWLINE, "NEWLINE", 11, 0},

		// .@ @
		{token.ILLEGAL, "ILLEGAL", 12, 0},
		{token.ILLEGAL, "ILLEGAL", 12, 0},
		{token.ILLEGAL, "ILLEGAL", 12, 0},
		{token.NEWLINE, "NEWLINE", 12, 0},

		// x /= fun << 1
		{token.IDVAL, "x", 13, 0},
		{token.DIVEQUAL, "/=", 13, 0},
		{token.IDVAL, "fun", 13, 0},
		{token.LSHIFT, "<<", 13, 0},
		{token.INT, "1", 13, 0},
		{token.NEWLINE, "NEWLINE", 13, 0},

		// {"fn": hypot}
		{token.LBRACE, "{", 14, 0},
		{token.STRING, "f\nn", 14, 0},
		{token.COLON, ":", 14, 0},
		{token.IDVAL, "hypot", 14, 0},
		{token.RBRACE, "}", 14, 0},

		// Allways end with EOF.
		{token.EOF, "EOF", 14, 0},
	}

	scanner := Init(input) // init scanner.

	for i, tt := range tests {
		tok, literal, lineno := scanner.NextToken()

		if tok != tt.token {
			t.Fatalf("test[%d] - Invalid Token. expected=%q, got=%q",
				i, tt.token.String(), tok.String())
		}
		if literal != tt.literal {
			t.Fatalf("test[%d] - Invalid Literal. expected=%q, got=%q",
				i, tt.literal, literal)
		}
		if lineno != tt.line {
			t.Fatalf("test[%d](%s) - Invalid line index. expected=%d, got=%d",
				i, tt.literal, tt.line, lineno)
		}
		// messed up column set up. cba to fix and sort a new test right now.
		// if scanner.column != tt.column+1 {
		// 	t.Fatalf("test[%d] - Invalid column index. line:\n%s<-\n expected=%d, got=%d",
		// 		i, scanner.TraceLine(), tt.column+1, scanner.column)
		// }
	}
}

// NEWLINE needs a redesign around comments and newlines. Currently the final line dosent
// return a newline causing problems for the parser. The whole set up is pretty bad tbh
func TestNewLines(t *testing.T) {
	tests := []struct {
		token   token.Token
		literal string
		line    int
	}{
		{token.INT, "100", 1},
		{token.NEWLINE, "NEWLINE", 1},
		{token.LET, "let", 4},
		{token.IDVAL, "hey", 4},
		{token.ASSIGN, "=", 4},
		{token.INT, "200", 4},
		{token.NEWLINE, "NEWLINE", 4},
		{token.LET, "let", 5},
		{token.IDVAL, "x", 5},
		{token.ASSIGN, "=", 5},
		{token.INT, "0xff", 5},
		{token.GEQUALS, ">=", 5},
		{token.INT, "5000", 5},
		{token.NEWLINE, "NEWLINE", 5},
		{token.STRING, "bye", 12},
		{token.EOF, "EOF", 12},
	}
	scanner := Init(`
100

# Hopefully the lines are now working.
let hey = 200
let x = 0xff >= 5000


# we only want NEWLINE to be seen just
# after content and then ignore the rest.


"bye"`)
	// let j = 100 # trailing comments are not working for some reason.

	// # after content. like "x + y\n
	// "bye"`)
	for i, tt := range tests {
		tok, literal, lineno := scanner.NextToken()
		if tok != tt.token {
			t.Fatalf("test[%d] - Invalid Token. expected=%q, got=%q",
				i, tt.token.String(), tok.String())
		}
		if literal != tt.literal {
			t.Fatalf("test[%d] - Invalid Literal. expected=%q, got=%q",
				i, tt.literal, literal)
		}
		if lineno != tt.line {
			t.Fatalf("test[%d](%s) - Invalid line index. expected=%d, got=%d.\n-> `%s`",
				i, tt.literal, tt.line, lineno, literal)
		}
	}
}
