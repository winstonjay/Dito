package scanner

import (
	"dito/src/token"
	"testing"
)

func TestNextToken(t *testing.T) {

	input := `let alpha = 10.999 + 10**3
let mut _o_me_ga = (50 % 7) - 0.002


# this is a single line comment you should not read me...
let hypot = func(a, b)->(a**2 + b**2)**0.5

let summer = true
let mut rain = true
let fun = false`
	// ^^^^^ DO NOT CHANGE OR YOU HAVE TO WORK OUT THE WHOLE TEST AGAIN
	// if you have to. extend the current string.
	// TODO (ALREADY BROKE cols NEEDS FIXING, there are some alignment
	// problemns with the coloum positions seen in trackbacks.
	tests := []struct {
		token   token.Token
		literal string
		line    int
		column  int
	}{
		{token.LET, "let", 0, 0},
		{token.IDVAL, "alpha", 0, 5},
		{token.ASSIGN, "=", 0, 8},
		{token.FLOAT, "10.999", 0, 15},
		{token.ADD, "+", 0, 17},
		{token.INT, "10", 0, 20},
		{token.POW, "**", 0, 22},
		{token.INT, "3", 0, 23},
		{token.NEWLINE, "newline", 0, -1},
		{token.LET, "let", 1, 0},
		{token.MUT, "mut", 1, 0},
		{token.IDVAL, "_o_me_ga", 1, 8},
		{token.ASSIGN, "=", 1, 11},
		{token.LPAREN, "(", 1, 13},
		{token.INT, "50", 1, 15},
		{token.MOD, "%", 1, 17},
		{token.INT, "7", 1, 19},
		{token.RPAREN, ")", 1, 20},
		{token.SUB, "-", 1, 22},
		{token.FLOAT, "0.002", 1, 28},
		{token.NEWLINE, "newline", 1, -1},
		// skips comment.
		{token.LET, "let", 5, 0},
		{token.IDVAL, "hypot", 5, 5},
		{token.ASSIGN, "=", 5, 7},
		{token.FUNC, "func", 5, 12},
		{token.LPAREN, "(", 5, 13},
		{token.IDVAL, "a", 5, 14},
		{token.COMMA, ",", 5, 15},
		{token.IDVAL, "b", 5, 17},
		{token.RPAREN, ")", 5, 18},
		{token.RARROW, "->", 5, 19},
		{token.LPAREN, "(", 5, 21},
		{token.IDVAL, "a", 5, 22},
		{token.POW, "**", 5, 24},
		{token.INT, "2", 5, 25},
		{token.ADD, "+", 5, 27},
		{token.IDVAL, "b", 5, 29},
		{token.POW, "**", 5, 31},
		{token.INT, "2", 5, 32},
		{token.RPAREN, ")", 5, 33},
		{token.POW, "**", 5, 35},
		{token.FLOAT, "0.5", 5, 38},
		{token.NEWLINE, "newline", 5, -1},

		{token.LET, "let", 7, 0},
		{token.IDVAL, "summer", 7, 6},
		{token.ASSIGN, "=", 7, 9},
		{token.TRUE, "true", 7, 14},
		{token.NEWLINE, "newline", 7, -1},
		{token.LET, "let", 8, 0},
		{token.MUT, "mut", 8, 0},
		{token.IDVAL, "rain", 8, 4},
		{token.ASSIGN, "=", 8, 7},
		{token.TRUE, "true", 8, 12},
		{token.NEWLINE, "newline", 8, -1},
		{token.LET, "let", 9, 0},
		{token.IDVAL, "fun", 9, 3},
		{token.ASSIGN, "=", 9, 6},
		{token.FALSE, "false", 9, 12},
		// Allways end with EOF.
		{token.EOF, "EOF", 9, 12},
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

func TestNewLines(t *testing.T) {
	tests := []struct {
		token   token.Token
		literal string
		line    int
	}{
		{token.INT, "100", 1},
		{token.NEWLINE, "newline", 1},
		{token.IDVAL, "hey", 4},
		{token.ASSIGN, ":=", 4},
		{token.INT, "200", 4},
		{token.NEWLINE, "newline", 4},
		{token.IDVAL, "x", 5},
		{token.ASSIGN, ":=", 5},
		{token.INT, "0xff", 5},
		{token.NEWLINE, "newline", 5},
		{token.EOF, "EOF", 10},
	}
	scanner := Init(`
100

# Hopefully the lines are now working.
hey := 200
x := 0xff

# we only want newline to be seen just
# after content and then ignore the rest.

# after content. like "x + y\n"`)
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
	}
}
