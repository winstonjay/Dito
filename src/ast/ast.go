package ast

import (
	"bytes"
	"dito/src/token"
	"strings"
)

// Node :
type Node interface {
	tokenLiteral() string
	String() string
}

// Statement : Expression |
type Statement interface {
	Node
	statementNode()
}

// Expression :
type Expression interface {
	Node
	expressionNode()
}

/* -----------------------------------------------------------------------
Program


*/

// Program : list of statements
type Program struct {
	Statements []Statement
}

// tokenLiteral :
func (p *Program) tokenLiteral() string {
	if len(p.Statements) > 0 {
		// return p.statements[0].tokenLiteral()
		return "<DitoProgram>"
	}
	return ""
}
func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

/* -----------------------------------------------------------------------
Statements


*/

// BlockStatement : a group of one or more statments inside curly brackets.
type BlockStatement struct {
	Token      token.Token // "{"
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) tokenLiteral() string { return bs.Token.String() }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// AssignmentStatement : Identifier = expr | Identifier := expr
type AssignmentStatement struct {
	Token token.Token // either := or =
	Name  *Identifier
	Value Expression
}

func (as *AssignmentStatement) statementNode()       {}
func (as *AssignmentStatement) tokenLiteral() string { return as.Token.String() }
func (as *AssignmentStatement) String() string {
	var out bytes.Buffer
	out.WriteString(as.Name.String())
	out.WriteString(" " + as.tokenLiteral() + " ")
	out.WriteString(as.Value.String())
	return out.String()
}

// ExpressionStatement :
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) tokenLiteral() string { return es.Token.String() }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// IfStatement :
type IfStatement struct {
	Token       token.Token // 'if'
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfStatement) statementNode()       {}
func (ie *IfStatement) tokenLiteral() string { return ie.Token.String() }
func (ie *IfStatement) String() string {
	var out bytes.Buffer
	out.WriteString("if ")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())
	if ie.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}

// ForStatement :
type ForStatement struct {
	Token     token.Token // 'for'
	Condition Expression
	LoopBody  *BlockStatement
}

func (fs *ForStatement) statementNode()       {}
func (fs *ForStatement) tokenLiteral() string { return fs.Token.String() }
func (fs *ForStatement) String() string {
	return "for statement"
}

// ImportStatement :
type ImportStatement struct {
	Token token.Token // 'import'
	Value string
}

func (is *ImportStatement) statementNode()       {}
func (is *ImportStatement) tokenLiteral() string { return is.Token.String() }
func (is *ImportStatement) String() string {
	return "Import statement"
}

/* -----------------------------------------------------------------------
Function Types


*/

// LambdaFunction : single expression function.
// fn(args) -> expr
type LambdaFunction struct {
	Token      token.Token // token.FUNC
	Parameters []*Identifier
	Expr       Expression
}

func (lf *LambdaFunction) expressionNode()      {}
func (lf *LambdaFunction) tokenLiteral() string { return lf.Token.String() }
func (lf *LambdaFunction) String() string       { return "<LambdaFunction>" }

// CallExpression :
type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) tokenLiteral() string { return ce.Token.String() }

func (ce *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

/* -----------------------------------------------------------------------
Expressions


*/

// PrefixExpression :
type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) tokenLiteral() string { return pe.Token.String() }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

// InfixExpression :
type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) tokenLiteral() string { return ie.Token.String() }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

// IfElseExpression :
type IfElseExpression struct {
	Token       token.Token // 'if'
	Initial     Expression
	Condition   Expression
	Alternative Expression
}

func (ie *IfElseExpression) expressionNode()      {}
func (ie *IfElseExpression) tokenLiteral() string { return ie.Token.String() }
func (ie *IfElseExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if ")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Initial.String())
	if ie.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}

/* -----------------------------------------------------------------------
Atoms


*/

// ArrayLiteral : Arrays can contain an assorted range of elements.
type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) tokenLiteral() string { return al.Token.String() }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer
	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

// Identifier : alphanumeric varible name.
type Identifier struct {
	Token token.Token // token.IDVAL
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) tokenLiteral() string { return i.Token.String() }
func (i *Identifier) String() string       { return i.Value }

// StringLiteral :
type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) tokenLiteral() string { return sl.Token.String() }
func (sl *StringLiteral) String() string       { return sl.Value }

// IntegerLiteral :  any non decimal numeric constant between
type IntegerLiteral struct {
	Token   token.Token // token.INT
	Literal string      // int as a string repr
	Value   int64       // int as a int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) tokenLiteral() string { return il.Literal }
func (il *IntegerLiteral) String() string       { return il.Literal }

// FloatLiteral :
type FloatLiteral struct {
	Token   token.Token
	Literal string
	Value   float64
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) tokenLiteral() string { return fl.Literal }
func (fl *FloatLiteral) String() string       { return fl.Literal }

// BooleanLiteral :
type BooleanLiteral struct {
	Token token.Token
	// no literal because Token.String() will be true or false
	Value bool
}

func (bl *BooleanLiteral) expressionNode()      {}
func (bl *BooleanLiteral) tokenLiteral() string { return bl.Token.String() }
func (bl *BooleanLiteral) String() string       { return bl.Token.String() }
