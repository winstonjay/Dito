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
	FLOAT  // Generic Floats.
	STRING // Strings starting and ending with double quotes.
	BOOL   // Generic bool
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
	CAT     // ++

	SHIFTL // <<
	SHIFTR // >>
	BITAND // &
	BITOR  // \
	BITXOR // ^

	LARROW // <- n.a.
	RARROW // ->

	LPAREN   // (
	RPAREN   // )
	RBRACE   // {
	LBRACE   // }
	LBRACKET // [
	RBRACKET // ]

	beginAssignementOp
	ASSIGN // =

	ADDEQUAL // +=
	SUBEQUAL // -=
	MULEQUAL // *=
	DIVEQUAL // /=
	MODEQUAL // %=
	endAssignementOp

	SEMI  // ;
	COLON // :
	COMMA // ,
	endOperator

	NEWLINE // \n

	QUOTE // "

	HASH // # start comments

	beginKeyword
	TRUE  // true
	FALSE // false
	IF    // if
	ELSE  // else
	FOR   // for
	IN    // in
	FUNC  // func
	AND   // and
	OR    // or
	LET   // let
	MUT   // mut
	// RHO    // rho n.a
	// IOTA   // iota n.a
	RETURN // return
	IMPORT // import
	endKeyword
)

var tokensLiterals = [...]string{

	ILLEGAL: "Illegal Token!",
	EOF:     "EOF",
	NEWLINE: "newline",

	IDVAL:  "ID",
	INT:    "Int",
	FLOAT:  "Float",
	STRING: "String",
	BOOL:   "Bool",

	SUB:  "-",
	ADD:  "+",
	MUL:  "*",
	DIV:  "/",
	IDIV: "//",
	MOD:  "%",
	POW:  "**",
	NOT:  "!",
	CAT:  "++",

	EQUALS:  "==",
	NEQUALS: "!=",
	LEQUALS: "<=",
	GEQUALS: ">=",
	LTHAN:   "<",
	GTHAN:   ">",

	SHIFTL: "<<",
	SHIFTR: ">>",
	BITAND: "&",
	BITOR:  "|",
	BITXOR: "^",

	ASSIGN: "=",

	ADDEQUAL: "+=",
	SUBEQUAL: "-=",
	MULEQUAL: "*=",
	DIVEQUAL: "/=",
	MODEQUAL: "%=",

	LARROW: "<-", // n.a.
	RARROW: "->",

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

	IF:   "if",
	ELSE: "else",
	FOR:  "for",
	IN:   "in",
	FUNC: "func",
	AND:  "and",
	OR:   "or",

	LET: "let",
	MUT: "mut",

	// RHO:  "rho", n.a
	// IOTA:  "iota", n.a
	IMPORT: "import",
	RETURN: "return",
}

func (t Token) String() string {

	if 0 <= t && t < Token(len(tokensLiterals)) {
		return tokensLiterals[t]
	}
	return ILLEGAL.String()
}

// Keywords is a map connecting language keywords to their token values.
var Keywords = make(map[string]Token)

// init : creates fills the map when the program is inited.
func init() {
	for i := beginKeyword + 1; i < endKeyword; i++ {
		Keywords[tokensLiterals[i]] = Token(i)
	}
}

// LookUpIDVal : Query if a alphanumeric token is a keyword.
// If it is, return the keyword if it isn't return the generic
// IDVAL token.
func LookUpIDVal(IDString string) Token {
	if tok, isKeyword := Keywords[IDString]; isKeyword {
		return tok
	}
	return IDVAL
}

// Operator Precdence values declared for implementing
// Pratt-like parsing method. They are called in this package
// by the .Precedence() method return.
const (
	_           uint = iota
	LOWEST           // non operators / default.
	CONDITONS        // if, for
	EQUALITY         // == !=
	LESSGREATER      // <= >= < >
	ADDSUB           // + -
	TERM             // * / % //
	EXPONENT         // **
	PREFIX           // unary operators; eg. + - !
	CALL             // Bracketed expressions, function calls; eg. foobar()
	HIGHEST          // Is extranous... pretty much there just in case. n.a.
)

// TODO Adjust operator precedence for bitwise operations.
// Implement all the bitwise operators.

// Precedence Returns the parsing precedence of a given token.
// values range from token.LOWEST to token.HIGHEST, constants
// defined in this package.
func (t Token) Precedence() uint {
	switch t {
	case IF:
		return CONDITONS
	case EQUALS, NEQUALS, IN:
		return EQUALITY
	case LEQUALS, GEQUALS, LTHAN, GTHAN:
		return LESSGREATER
	case ADD, SUB, CAT, SHIFTL, SHIFTR:
		return ADDSUB
	case MOD, DIV, MUL, IDIV:
		return TERM
	case POW:
		return EXPONENT
	case LPAREN:
		return CALL
	case LBRACKET:
		return HIGHEST
	}
	return LOWEST
}

// IsAssignmentOp replies true or false depending whether a token is
// an assignmen operator. An assignment operator is described token that
// ends with '='.
func (t Token) IsAssignmentOp() bool {
	return beginAssignementOp < t && t < endAssignementOp
}
