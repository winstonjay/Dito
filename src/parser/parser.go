// Package parser implements Dito's Parser.
// Its is implemented blah blah blah...
// Top Down Operator Precedence parsing. https://goo.gl/uoH6Ta
package parser

import (
	"dito/src/ast"
	"dito/src/lexer"
	"dito/src/token"
	"strconv"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// Parser : structure who's methods implement a pratt parser.
type Parser struct {
	scanner        *lexer.Scanner
	errors         []*ParseError
	currentToken   token.Token
	currentLiteral string
	peekToken      token.Token
	peekLiteral    string
	currentLine    int
	peekTokenLine  int

	prefixParseFns map[token.Token]prefixParseFn
	infixParseFns  map[token.Token]infixParseFn
}

// New : Initalise a new parser.
func New(s *lexer.Scanner) *Parser {
	p := &Parser{
		scanner: s,
		errors:  []*ParseError{},
	}
	p.prefixParseFns = map[token.Token]prefixParseFn{
		// prefix / unary expressions
		token.SUB: p.prefixExpression,
		token.ADD: p.prefixExpression,
		token.NOT: p.prefixExpression,
		//
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
	}
	// twice to fill current and peek token.
	p.nextToken()
	p.nextToken()
	return p
}

// Refresh :
func (p *Parser) Refresh(s *lexer.Scanner) {
	p.scanner = s
	p.errors = []*ParseError{}
	p.nextToken()
	p.nextToken()
}

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

// // expressions are ended by a semicolon or a newline.
// // there is no newline token, but we can see if the line number
// // has changed from the scanners positon.
func (p *Parser) endExpression(lineno int) bool {
	if p.currentLine > lineno {
		return true
	}
	if p.peekTokenIs(token.SEMI) {
		p.nextToken()
		return true
	}
	return false
}

// ParseProgram : creates ast of the inputed text incrementally
// working with the scanner.
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for p.currentToken != token.EOF {
		stmt := p.statement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

/*
######## Statements.
*/

func (p *Parser) statement() ast.Statement {
	switch p.currentToken {
	case token.IMPORT:
		return p.importStatement()
	case token.IDVAL:
		if p.peekTokenIs(token.NEWASSIGN) || p.peekTokenIs(token.REASSIGN) {
			return p.assignmentStatement()
		}
		return p.expressionStatement()
	case token.IF:
		return p.ifElseStatement()
	case token.FOR:
		return p.forStatement()
	default:
		return p.expressionStatement()
	}
}

func (p *Parser) assignmentStatement() *ast.AssignmentStatement {
	stmt := &ast.AssignmentStatement{}
	stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentLiteral}
	p.nextToken() // we already saw what was ahead before we came here.
	stmt.Token = p.currentToken
	p.nextToken() // this could be anything tbh.
	stmt.Value = p.expression(token.LOWEST)

	// inforce semicolons till we find sort a newline strategy.
	if !p.currentTokenIs(token.SEMI) {
		p.peekError(token.SEMI)
		return nil
	}
	return stmt
}

func (p *Parser) expressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.currentToken}
	stmt.Expression = p.expression(token.LOWEST)
	// inforce semicolons till we find sort a newline strategy.
	if !p.currentTokenIs(token.SEMI) {
		p.peekError(token.SEMI)
		return nil
	}
	return stmt
}

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

func (p *Parser) forStatement() *ast.ForStatement {
	stmt := &ast.ForStatement{Token: p.currentToken}
	p.nextToken()
	stmt.Condition = p.expression(token.LOWEST)
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	stmt.LoopBody = p.blockStatement()
	return stmt
}

func (p *Parser) blockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.currentToken}
	block.Statements = []ast.Statement{}
	p.nextToken()
	for !p.currentTokenIs(token.RBRACE) && !p.currentTokenIs(token.EOF) {
		stmt := p.statement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

func (p *Parser) importStatement() *ast.ImportStatement {
	is := &ast.ImportStatement{Token: p.currentToken}
	if !p.expectPeek(token.IDVAL) {
		return nil
	}
	is.Value = p.identifier().(*ast.Identifier).Value
	return is
}

/*
	Functions.
*/

// lambdaFunction: fn(<parameters>) -> <expr>
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

func (p *Parser) functionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}
	p.nextToken()
	idVal := &ast.Identifier{Token: p.currentToken, Value: p.currentLiteral}
	identifiers = append(identifiers, idVal)
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		idVal := &ast.Identifier{Token: p.currentToken, Value: p.currentLiteral}
		identifiers = append(identifiers, idVal)
	}
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return identifiers
}

func (p *Parser) callExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.currentToken, Function: function}
	exp.Arguments = p.expressionList(token.RPAREN)
	return exp
}

/*
	Expressions.
*/

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

func (p *Parser) expression(precedence uint) ast.Expression {
	lineno := p.currentLine
	prefix := p.prefixParseFns[p.currentToken]
	if prefix == nil {
		p.noParseFnError(p.currentToken)
		return nil
	}
	expr := prefix()
	for !p.endExpression(lineno) && precedence < p.peekToken.Precedence() {
		infix := p.infixParseFns[p.peekToken]
		if infix == nil {
			return expr
		}
		p.nextToken()
		expr = infix(expr)
	}
	// this causing bugs right here.
	if lineno == p.currentLine {
		if p.peekTokenIs(token.IF) {
			p.nextToken()
			return p.ifElseExpression(expr)
		}
	}
	return expr
}

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
	expr.Alternative = p.expression(token.LOWEST)
	return expr
}

func (p *Parser) prefixExpression() ast.Expression {
	expr := &ast.PrefixExpression{
		Token:    p.currentToken,
		Operator: p.currentLiteral,
	}
	p.nextToken()
	expr.Right = p.expression(token.PREFIX)
	return expr
}

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

func (p *Parser) groupedExpression() ast.Expression {
	p.nextToken()
	expr := p.expression(token.LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return expr
}

func (p *Parser) indexExpression(item ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.currentToken, Left: item}
	p.nextToken()
	exp.Index = p.expression(token.LOWEST)
	if !p.expectPeek(token.RBRACKET) {
		return nil
	}
	return exp
}

/*
######## Atoms / Type Literals
*/

func (p *Parser) arrayLiteral() ast.Expression {
	return &ast.ArrayLiteral{
		Token:    p.currentToken,
		Elements: p.expressionList(token.RBRACKET),
	}
}

func (p *Parser) stringLiteral() ast.Expression {
	return &ast.StringLiteral{
		Token: p.currentToken,
		Value: p.currentLiteral,
	}
}

func (p *Parser) identifier() ast.Expression {
	return &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentLiteral,
	}
}

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

func (p *Parser) booleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{
		Token: p.currentToken,
		Value: p.currentTokenIs(token.TRUE),
	}
}
