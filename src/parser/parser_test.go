package parser

import (
	"dito/src/ast"
	"dito/src/lexer"
	"fmt"
	"testing"
)

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, err := range errors {
		t.Errorf("parser error: %q", err.message)
	}
	t.FailNow()
}

func TestAssignmentStatement(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"alpha := 100", "alpha", 100},
		{"y = 99999999999;", "y", 99999999999},
		{"beta = 2020202021", "beta", 2020202021},
		{"floater := 202.900000", "floater", 202.900000},
		{"half := 0.500000", "half", 0.500000},
		{"do = true", "do", true},
		{"maybe = false", "maybe", false},
	}

	for i, tt := range tests {

		scanner := lexer.Init(tt.input)
		parser := New(scanner)
		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("test[%d] Wrong number of program.Statments. expected=%d. got=%d",
				i, 1, len(program.Statements))
		}
		s := program.Statements[0]
		stmt, ok := s.(*ast.AssignmentStatement)
		if !ok {
			t.Fatalf("test[%d] not *ast.AssignmentStatement. got=%T", i, s)
		}
		if stmt.Name.Value != tt.expectedIdentifier {
			t.Fatalf("test[%d] unexpected identifier name. Expected=%s got=%s.",
				i, tt.expectedIdentifier, stmt.Name.Value)
		}

		if !testLiteralExpression(t, stmt.Value, tt.expectedValue) {
			return
		}
	}
}

func TestExpressionStatement(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"100", 100},
		{"true", true},
		{"false", false},
		{"2.5", 2.50000},
	}

	for i, tt := range tests {

		scanner := lexer.Init(tt.input)
		parser := New(scanner)
		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("test[%d] Wrong number of program.Statments. expected=%d. got=%d",
				i, 1, len(program.Statements))
		}
	}
}

func TestInfixExpression(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		right    interface{}
	}{
		{"-100", "-", 100},
		{"!true", "!", true},
		{"!false", "!", false},
		{"-2.5", "-", 2.50000},
		{"-foo", "-", "foo"},
		{"!foo", "!", "foo"},
		{"+foo", "+", "foo"},
		{"+5", "+", 5},
	}
	for i, tt := range tests {
		scanner := lexer.Init(tt.input)
		parser := New(scanner)
		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("test[%d] Wrong number of program.Statments. expected=%d. got=%d",
				i, 1, len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}
		if !testLiteralExpression(t, exp.Right, tt.right) {
			return
		}
	}
}

func TestInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"foobar + barfoo;", "foobar", "+", "barfoo"},
		{"foobar - barfoo;", "foobar", "-", "barfoo"},
		{"foobar * barfoo;", "foobar", "*", "barfoo"},
		{"foobar / barfoo;", "foobar", "/", "barfoo"},
		{"foobar > barfoo;", "foobar", ">", "barfoo"},
		{"foobar < barfoo;", "foobar", "<", "barfoo"},
		{"foobar == barfoo;", "foobar", "==", "barfoo"},
		{"foobar != barfoo;", "foobar", "!=", "barfoo"},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		l := lexer.Init(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

func testInfixExpression(
	t *testing.T,
	exp ast.Expression,
	left interface{},
	operator string,
	right interface{},
) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.OperatorExpression. got=%T(%s)", exp, exp)
		return false
	}
	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}
	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}
	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}
	return true
}

/*
	Literal Helpers.
*/

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case string:
		return testIdentifier(t, exp, string(v))
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case float64:
		return testFloatLiteral(t, exp, float64(v))
	case bool:
		return testBooleanLiteral(t, exp, bool(v))
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testIdentifier(t *testing.T, idv ast.Expression, value string) bool {
	idVal, ok := idv.(*ast.Identifier)
	if !ok {
		t.Errorf("idv not *ast.Identifier. got=%T", idv)
		return false
	}
	if idVal.Value != value {
		t.Errorf("Identifier.Value not %s. got=%s", value, idVal.Value)
		return false
	}
	return true
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}
	if integ.Value != value {
		t.Errorf("integerLiteral.Value not %d. got=%d", value, integ.Value)
		return false
	}
	if integ.Literal != fmt.Sprintf("%d", value) {
		t.Errorf("integerLiteral.Literal not %d. got=%s", value, integ.Literal)
		return false
	}
	return true
}

func testFloatLiteral(t *testing.T, fl ast.Expression, value float64) bool {
	flo, ok := fl.(*ast.FloatLiteral)
	if !ok {
		t.Errorf("flo not *ast.FloatLiteral. got=%T", flo)
		return false
	}
	if flo.Value != value {
		t.Errorf("FloatLiteral.Value wrong. expected=%f. got=%f", value, flo.Value)
		return false
	}
	// literals would require us to alter the structure of the tests as floats in go
	// when converted numerically contain trailing zeros.
	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.BooleanLiteral)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}
	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}
	if bo.String() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s",
			value, bo.String())
		return false
	}
	return true
}
