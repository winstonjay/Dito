// Package parser implements Dito's Parser.
package parser

// TODO more package docs so i dont forget who i am.

import (
	"dito/src/ast"
	"dito/src/lexer"
	"dito/src/token"
	"fmt"
	"strconv"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// Parser : structure who's methods implement a top-down
// operator precedence parser.
// (Based roughly on pratt parsing model: https://goo.gl/uoH6Ta)
// expression prefix and infix parse functions are stored in a
// table and mapped by particular tokens.
// It should be initialised with a pointer to a lexical scanner.
// Calling the method ParseProgram will return a fully formed AST.
// For notes on the structure of the parser  most parse functions
// are annotated with simple grammar notes. These are in the form
// nodetype: one or more alternative rules of what ast structure
// the function can return.
type Parser struct {
	scanner        *lexer.Scanner
	errors         []*ParseError
	currentToken   token.Token
	currentLiteral string
	peekToken      token.Token
	peekLiteral    string
	currentLine    int
	peekTokenLine  int
	openParen      bool
	prefixParseFns map[token.Token]prefixParseFn
	infixParseFns  map[token.Token]infixParseFn
}

// New : Initalise a new parser with an newly initialsed Scanner.
func New(s *lexer.Scanner) *Parser {
	p := &Parser{
		scanner: s,
		errors:  []*ParseError{},
	}
	// Define a table of methods for parsing expressions given a token.
	p.prefixParseFns = map[token.Token]prefixParseFn{
		// prefix / unary expressions
		token.SUB: p.prefixExpression,
		token.ADD: p.prefixExpression,
		token.NOT: p.prefixExpression,
		// function.
		token.FUNC: p.lambdaFunction,
		// Token Literals.
		token.IDVAL:    p.identifier,
		token.INT:      p.integerLiteral,
		token.FLOAT:    p.floatLiteral,
		token.TRUE:     p.booleanLiteral,
		token.FALSE:    p.booleanLiteral,
		token.LBRACKET: p.arrayLiteral,
		token.LPAREN:   p.groupedExpression,
		token.STRING:   p.stringLiteral,
	}

	p.infixParseFns = map[token.Token]infixParseFn{
		// infix / binary Expressions.
		token.SUB:      p.infixExpression,
		token.ADD:      p.infixExpression,
		token.MUL:      p.infixExpression,
		token.DIV:      p.infixExpression,
		token.IDIV:     p.infixExpression,
		token.MOD:      p.infixExpression,
		token.POW:      p.infixExpression,
		token.EQUALS:   p.infixExpression,
		token.NEQUALS:  p.infixExpression,
		token.LEQUALS:  p.infixExpression,
		token.GEQUALS:  p.infixExpression,
		token.LTHAN:    p.infixExpression,
		token.GTHAN:    p.infixExpression,
		token.SHIFTL:   p.infixExpression,
		token.SHIFTR:   p.infixExpression,
		token.LPAREN:   p.callExpression,
		token.LBRACKET: p.indexExpression,
		token.IF:       p.ifElseExpression,
	}
	// twice to fill current and peek token.
	p.nextToken()
	p.nextToken()
	return p
}

// Refresh : reset all values to defualt
// for parsing a fresh input stream.
func (p *Parser) Refresh(s *lexer.Scanner) {
	p.scanner = s
	p.errors = []*ParseError{}
	p.nextToken()
	p.nextToken()
}

// is the next token what we want if not create an error.
func (p *Parser) expectPeek(t token.Token) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekTokenIs(t token.Token) bool {
	return p.peekToken == t
}

func (p *Parser) currentTokenIs(t token.Token) bool {
	return p.currentToken == t
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.currentLiteral = p.peekLiteral
	p.peekToken, p.peekLiteral, p.currentLine = p.scanner.NextToken()
}

// check to see if there is a newline, EOF or Semicolon token or complain.
func (p *Parser) stmtEnd() bool {
	if p.peekTokenIs(token.SEMI) ||
		p.peekTokenIs(token.NEWLINE) ||
		p.peekTokenIs(token.EOF) {
		p.nextToken()
		return true
	}
	p.peekError(token.NEWLINE)
	return false
}

// ParseProgram creates ast of the inputed text incrementally
// working with the scanner.
// Program: list of statements
// stmtend: newline | semicolon | EOF
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for p.currentToken != token.EOF {
		stmt := p.statement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		if !p.stmtEnd() {
			return nil
		}
		p.nextToken()
	}
	return program
}

/*
#### Statements.
*/

// statement:
// 	   assignmentStatement
// 	   indexAssignmentStatement
// 	   expressionStatement
// 	   functionStatement
// 	   returnStatement
// 	   forStatement
// 	   importStatement
func (p *Parser) statement() ast.Statement {
	switch p.currentToken {
	case token.IDVAL:
		if p.peekToken.IsAssignmentOp() {
			return p.assignmentStatement()
		}
		if p.peekTokenIs(token.LBRACKET) {
			idxExp := p.expression(token.LOWEST)
			if !p.peekToken.IsAssignmentOp() {
				return &ast.ExpressionStatement{
					Token:      token.LBRACE,
					Expression: idxExp,
				}
			}
			return p.indexAssignmentStatement(idxExp.(*ast.IndexExpression))
		}
		return p.expressionStatement()
	case token.FUNC:
		return p.functionStatement()
	case token.RETURN:
		return p.returnStatement()
	case token.IF:
		return p.ifElseStatement()
	case token.FOR:
		return p.forStatement()
	case token.IMPORT:
		return p.importStatement()
	default:
		return p.expressionStatement()
	}
}

// returnStatement:
//	   'return' expression
func (p *Parser) returnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currentToken}
	p.nextToken()
	stmt.Value = p.expression(token.LOWEST)
	return stmt
}

// assignmentStatement:
//	   identifier assignmentOperator expression
func (p *Parser) assignmentStatement() *ast.AssignmentStatement {
	stmt := &ast.AssignmentStatement{}
	stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentLiteral}
	p.nextToken() // we already saw what was ahead before we came here.
	stmt.Token = p.currentToken
	p.nextToken() // this could be anything tbh.
	stmt.Value = p.expression(token.LOWEST)
	return stmt
}

// indexAssignmentStatement:
//	   identifier '[' expression ']' assignmentOperator expression
func (p *Parser) indexAssignmentStatement(idxExp *ast.IndexExpression) *ast.IndexAssignmentStatement {
	stmt := &ast.IndexAssignmentStatement{Token: token.LBRACE, IdxExp: idxExp}
	p.nextToken()
	p.nextToken()
	stmt.Value = p.expression(token.LOWEST)
	return stmt
}

// expressionStatement:
// 	   expression
func (p *Parser) expressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.currentToken}
	stmt.Expression = p.expression(token.LOWEST)
	return stmt
}

// ifElseStatement:
// 	   'if' expression '{' blockStatement '}'
// 	   'if' expression '{' blockStatement '}' 'else' '{' blockStatement '}'
func (p *Parser) ifElseStatement() *ast.IfStatement {
	expression := &ast.IfStatement{Token: p.currentToken}
	p.nextToken()
	expression.Condition = p.expression(token.LOWEST)
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	expression.Consequence = p.blockStatement()
	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		expression.Alternative = p.blockStatement()
	}
	return expression
}

// functionStatement
// 	   'func' identifier '(' functionParameters ')' '{' blockStatement '}'
func (p *Parser) functionStatement() *ast.Function {
	fn := &ast.Function{Token: p.currentToken}
	if !p.expectPeek(token.IDVAL) {
		return nil
	}
	fn.Name = p.identifier().(*ast.Identifier)
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	fn.Parameters = p.functionParameters()
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	fn.Body = p.blockStatement()
	return fn
}

// forStatement:
//     'for' identifier 'in' identifier '{' blockStatement '}'
// 	   'for' expression '{' blockStatement '}'
func (p *Parser) forStatement() *ast.ForStatement {
	stmt := &ast.ForStatement{Token: p.currentToken}
	p.nextToken()
	if p.currentTokenIs(token.IDVAL) && p.peekTokenIs(token.IN) {
		stmt.ID = p.identifier().(*ast.Identifier)
		p.nextToken()
		p.nextToken()
		stmt.Iter = p.expression(token.LOWEST)
	} else {
		stmt.ID = nil
		stmt.Condition = p.expression(token.LOWEST)
	}
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	stmt.LoopBody = p.blockStatement()
	return stmt
}

// blockStatement:
//     statement stmtend blockstatement
// 	   statement stmtend
func (p *Parser) blockStatement() *ast.BlockStatement {
	if p.peekTokenIs(token.NEWLINE) {
		p.nextToken()
	}
	block := &ast.BlockStatement{Token: p.currentToken}
	block.Statements = []ast.Statement{}
	p.nextToken()
	for !p.currentTokenIs(token.RBRACE) && !p.currentTokenIs(token.EOF) {
		stmt := p.statement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		if !p.stmtEnd() {
			return nil
		}
		p.nextToken()
	}
	return block
}

// importStatement
//     'import' identifier
func (p *Parser) importStatement() *ast.ImportStatement {
	is := &ast.ImportStatement{Token: p.currentToken}
	if !p.expectPeek(token.IDVAL) {
		return nil
	}
	is.Value = p.identifier().(*ast.Identifier).Value
	return is
}

/*
#### Expressions.
*/

// expression:
//     lambdaFunction
// 	   callExpression
// 	   ifElseExpression
// 	   prefixExpression
// 	   infixExpression
// 	   groupedExpression
// 	   indexExpression
// 	   identifier
// 	   atom
func (p *Parser) expression(precedence uint) ast.Expression {
	prefix := p.prefixParseFns[p.currentToken]
	// we want to be able to do multiline expr inside parenthesis.
	if prefix == nil {
		p.noParseFnError(p.currentToken)
		return nil
	}
	expr := prefix()
	for precedence < p.peekToken.Precedence() {
		infix := p.infixParseFns[p.peekToken]
		if infix == nil {
			return expr
		}
		p.nextToken()
		expr = infix(expr)
		if p.openParen && p.currentTokenIs(token.NEWLINE) {
			p.nextToken()
		}
	}
	return expr
}

// lambdaFunction:
// 	   'func' '(' functionParameters ')' '->' expression
func (p *Parser) lambdaFunction() ast.Expression {
	lambda := &ast.LambdaFunction{Token: p.currentToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	lambda.Parameters = p.functionParameters()
	if !p.expectPeek(token.RARROW) {
		return nil
	}
	p.nextToken()
	lambda.Expr = p.expression(token.LOWEST)
	return lambda
}

// functionParameters:
// 	   identifer ',' functionParameters
// 	   identifier
func (p *Parser) functionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}
	p.nextToken()
	idVal := p.identifier().(*ast.Identifier)
	identifiers = append(identifiers, idVal)
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		idVal := p.identifier().(*ast.Identifier)
		identifiers = append(identifiers, idVal)
	}
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return identifiers
}

// callExpression:
// 	   identifier '(' expressionList ')'
func (p *Parser) callExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.currentToken, Function: function}
	exp.Arguments = p.expressionList(token.RPAREN)
	return exp
}

// expressionList
// 	   expression ',' expressionList
// 	   expression
func (p *Parser) expressionList(delimiter token.Token) []ast.Expression {
	list := []ast.Expression{}
	if p.peekTokenIs(delimiter) {
		p.nextToken()
		return list
	}
	p.nextToken()
	list = append(list, p.expression(token.LOWEST))
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.expression(token.LOWEST))
	}
	if !p.expectPeek(delimiter) {
		return nil
	}
	return list
}

// ifElseExpression:
// 		expression 'if' expression 'else' expression
func (p *Parser) ifElseExpression(inital ast.Expression) ast.Expression {
	expr := &ast.IfElseExpression{
		Initial: inital,
		Token:   p.currentToken,
	}
	p.nextToken()
	expr.Condition = p.expression(token.LOWEST)
	if !p.expectPeek(token.ELSE) {
		return nil
	}
	p.nextToken()
	if p.openParen && p.currentTokenIs(token.NEWLINE) {
		p.nextToken()
	}
	expr.Alternative = p.expression(token.LOWEST)
	return expr
}

// prefixExpression:
//     prefixOperator expression
func (p *Parser) prefixExpression() ast.Expression {
	expr := &ast.PrefixExpression{
		Token:    p.currentToken,
		Operator: p.currentLiteral,
	}
	p.nextToken()
	expr.Right = p.expression(token.PREFIX)
	return expr
}

// infixExpression:
//     expression binaryOperator expression
func (p *Parser) infixExpression(left ast.Expression) ast.Expression {
	expr := &ast.InfixExpression{
		Token:    p.currentToken,
		Operator: p.currentLiteral,
		Left:     left,
	}
	precedence := p.currentToken.Precedence()
	p.nextToken()
	expr.Right = p.expression(precedence)
	return expr
}

// groupedExpression:
// 	   '(' expression ')'
func (p *Parser) groupedExpression() ast.Expression {
	p.nextToken()
	p.openParen = true
	if p.currentTokenIs(token.NEWLINE) {
		p.nextToken()
	}
	expr := p.expression(token.LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	p.openParen = false
	return expr
}

// indexExpression:
// 	   identifier '[' expression ']'
func (p *Parser) indexExpression(item ast.Expression) ast.Expression {
	id, ok := item.(*ast.Identifier)
	if !ok {
		fmt.Printf("ERROR ASSERTION")
		return nil
	}
	exp := &ast.IndexExpression{Token: p.currentToken, Left: id}
	p.nextToken()
	exp.Index = p.expression(token.LOWEST)
	if !p.expectPeek(token.RBRACKET) {
		return nil
	}
	return exp
}

/*
#### Atoms / Type Literals
*/

// arrayLiteral: '[' expressionList ']'
func (p *Parser) arrayLiteral() ast.Expression {
	return &ast.ArrayLiteral{
		Token:    p.currentToken,
		Elements: p.expressionList(token.RBRACKET),
	}
}

// stringLiteral: "[^"]*"
func (p *Parser) stringLiteral() ast.Expression {
	return &ast.StringLiteral{
		Token: p.currentToken,
		Value: p.currentLiteral,
	}
}

// identifier: [A-Za-z_][A-Za-z_0-9]*
func (p *Parser) identifier() ast.Expression {
	return &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentLiteral,
	}
}

// integer: base10Int | hexInt
func (p *Parser) integerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.currentToken, Literal: p.currentLiteral}
	value, err := strconv.ParseInt(p.currentLiteral, 0, 64)
	if err != nil {
		p.genericError("integer", err)
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) floatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.currentToken, Literal: p.currentLiteral}
	value, err := strconv.ParseFloat(p.currentLiteral, 64)
	if err != nil {
		p.genericError("float", err)
		return nil
	}
	lit.Value = value
	return lit
}

// boolean: true | false
func (p *Parser) booleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{
		Token: p.currentToken,
		Value: p.currentTokenIs(token.TRUE),
	}
}
