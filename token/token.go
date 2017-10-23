package token

// Token : type that stores dito tokens its string method
// converts there int forms into human readable output.
type Token uint

// Token constants defined by the dito programming language.
// Tokens that are included but not currently implemented are
// marked with a 'n.a.'
const (
	EOF     Token = iota // end of file 0
	ILLEGAL              // non recognised tokens.

	beginLiteral
	IDVAL  // Alphanumeric idenifiers (varible names).
	INT    // Generic Integers.
	FLOAT  // Generic float token.
	STRING // Strings.
	BOOL   // generic bool
	endLiteral

	beginOperator
	ADD  // +
	SUB  // -
	MUL  // *
	DIV  // /
	IDIV // //
	MOD  // %
	POW  // **

	EQUALS  // ==
	NEQUALS // !=
	LEQUALS // <=
	GEQUALS // >=
	LTHAN   // <
	GTHAN   // >
	NOT     // !

	LARROW // <- n.a.
	RARROW // -> n.a.

	LPAREN   // (
	RPAREN   // )
	RBRACE   // {
	LBRACE   // }
	LBRACKET // [
	RBRACKET // ]

	NEWASSIGN // :=
	REASSIGN  // =

	SEMI  // ;
	COLON // :
	COMMA // ,
	endOperator

	QUOTE // "

	HASH // #

	beginKeyword
	TRUE   // true
	FALSE  // false
	IF     // if
	ELSE   // else
	FOR    // for
	FUNC   // func
	RETURN // return n.a.
	IMPORT // import n.a.
	endKeyword
)

var tokensLiterals = [...]string{

	ILLEGAL: "Illegal Token!",
	EOF:     "EOF",

	IDVAL:  "ID",
	INT:    "Int",
	FLOAT:  "Float",
	STRING: "String",
	BOOL:   "Bool",

	ADD:  "+",
	SUB:  "-",
	MUL:  "*",
	DIV:  "/",
	IDIV: "//",
	MOD:  "%",
	POW:  "**",
	NOT:  "!",

	EQUALS:  "==",
	NEQUALS: "!=",
	LEQUALS: "<=",
	GEQUALS: ">=",
	LTHAN:   "<",
	GTHAN:   ">",

	NEWASSIGN: ":=",
	REASSIGN:  "=",

	LARROW: "<-", // n.a.
	RARROW: "->", // n.a.

	LPAREN:   "(",
	RPAREN:   ")",
	LBRACE:   "{",
	RBRACE:   "}",
	LBRACKET: "[",
	RBRACKET: "]",

	SEMI:  ";",
	COLON: ":",
	COMMA: ",",

	QUOTE: "\"",

	TRUE:  "true",
	FALSE: "false",

	IF:     "if",
	ELSE:   "else",
	FOR:    "for",
	FUNC:   "func",
	IMPORT: "import",
	RETURN: "return",
}

func (t Token) String() string {

	if 0 <= t && t < Token(len(tokensLiterals)) {
		return tokensLiterals[t]
	}
	return ILLEGAL.String()
}

var keywords = make(map[string]Token)

// Keywords : just expose sor this out later.
var Keywords = keywords

// init : creates fills the map when the program is inited.
func init() {
	for i := beginKeyword + 1; i < endKeyword; i++ {
		keywords[tokensLiterals[i]] = Token(i)
	}
}

// LookUpIDVal : Query if a alphanumeric token is a keyword.
// If it is, return the keyword if it isn't return the generic
// IDVAL token.
func LookUpIDVal(IDString string) Token {
	if tok, isKeyword := keywords[IDString]; isKeyword {
		return tok
	}
	return IDVAL
}

// Set of constants denoting operator precedence values starting
// from 1 to 9.
const (
	_           uint = iota
	LOWEST           // non operators / default.
	EQUALITY         // == !=
	LESSGREATER      // <= >= < >
	ADDSUB           // + -
	TERM             // * / %
	EXPONENT         // **
	PREFIX           // unary operators; eg. + - !
	CALL             // Bracketed expressions, function calls; eg. foobar()
	HIGHEST          // Is extranous... pretty much there just in case. n.a.
)

// Precedence : get tokens precedence to help parse it later.
func (t Token) Precedence() uint {
	switch t {
	case EQUALS, NEQUALS:
		return EQUALITY
	case LEQUALS, GEQUALS, LTHAN, GTHAN:
		return LESSGREATER
	case ADD, SUB:
		return ADDSUB
	case MOD, DIV, MUL:
		return TERM
	case POW:
		return EXPONENT
	case LPAREN:
		return CALL
	}
	return LOWEST
}
