package parser

import (
	"dito/src/ast"
	"dito/src/scanner"
	"dito/src/token"
	"fmt"
	"testing"
)

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected *ast.IfElseExpression
	}{
		{
			"true if true else false",
			&ast.IfElseExpression{
				Token:       token.IF,
				Initial:     &ast.BooleanLiteral{Token: token.TRUE, Value: true},
				Condition:   &ast.BooleanLiteral{Token: token.TRUE, Value: true},
				Alternative: &ast.BooleanLiteral{Token: token.FALSE, Value: false},
			},
		}, {
			"x if y else z",
			&ast.IfElseExpression{
				Token:       token.IF,
				Initial:     &ast.Identifier{Token: token.IDVAL, Value: "x"},
				Condition:   &ast.Identifier{Token: token.IDVAL, Value: "y"},
				Alternative: &ast.Identifier{Token: token.IDVAL, Value: "z"},
			},
		},

		// TODO multiline expressions are broken.
	}
	for i, tt := range tests {
		program := parseTestProgram(t, tt.input)
		testStatementsLen(t, program, i, 1)
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}
		expr, ok := stmt.Expression.(*ast.IfElseExpression)
		if !ok {
			t.Fatalf("stmt is not %T. got=%T", tt.expected, stmt.Expression)
		}
		testAtomLiteral(t, expr.Initial, tt.expected.Initial)
		testAtomLiteral(t, expr.Condition, tt.expected.Condition)
		testAtomLiteral(t, expr.Alternative, tt.expected.Alternative)
	}
}

func testAtomLiteral(t *testing.T, in ast.Expression, expected ast.Expression) {
	switch v := expected.(type) {
	case *ast.BooleanLiteral:
		testBooleanLiteral(t, in, v.Value)
	case *ast.Identifier:
		testIdentifier(t, in, v.Value)
	case *ast.IntegerLiteral:
		testIntegerLiteral(t, in, v.Value)
	case *ast.FloatLiteral:
		testFloatLiteral(t, in, v.Value)
	default:
		t.Fatalf("Unkown type. got=%T. want=%T", in, expected)
	}
}

func TestAssignmentStatement(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
		expectedToken      token.Token
	}{
		{"let alpha = 100", "alpha", 100, token.LET},
		{"let mut y = 99999999999;", "y", 99999999999, token.MUT},
		{"let beta = 2020202021", "beta", 2020202021, token.LET},
		{"let mut floater = 202.900000", "floater", 202.900000, token.MUT},
		{"let half = 0.500000", "half", 0.500000, token.LET},
		{"let do = true", "do", true, token.LET},
		{"let maybe = false", "maybe", false, token.LET},
	}

	for i, tt := range tests {
		program := parseTestProgram(t, tt.input)
		testStatementsLen(t, program, i, 1)
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

func parseTestProgram(t *testing.T, input string) *ast.Program {
	scanner := scanner.Init(input)
	parser := New(scanner)
	program := parser.ParseProgram()
	// throw if error
	if errors := parser.Errors(); len(errors) != 0 {
		t.Errorf("parser has %d errors", len(errors))
		for _, err := range errors {
			t.Errorf("parser error: %q", err.message)
		}
		t.FailNow()
	}
	return program
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
		program := parseTestProgram(t, tt.input)
		testStatementsLen(t, program, i, 1)

		// TODO
	}
}

func TestPrefixExpression(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		right    interface{}
	}{
		{"-100", "-", 100},
		{"not true", "not", true},
		{"not false", "not", false},
		{"-2.5", "-", 2.50000},
		{"-foo", "-", "foo"},
		{"not foo", "not", "foo"},
		{"+foo", "+", "foo"},
		{"+5", "+", 5},
	}
	for i, tt := range tests {
		program := parseTestProgram(t, tt.input)
		testStatementsLen(t, program, i, 1)

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

func testStatementsLen(t *testing.T, prog *ast.Program, i, testval int) {
	if len(prog.Statements) != testval {
		t.Fatalf("test[%d] Wrong number of program.Statments. expected=%d. got=%d",
			i, 1, len(prog.Statements))
	}
}

func TestInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
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

	for i, tt := range infixTests {
		program := parseTestProgram(t, tt.input)
		testStatementsLen(t, program, i, 1)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		if !testInfixExpression(
			t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue) {
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
		return testIntegerLiteral(t, exp, v)
	case float64:
		return testFloatLiteral(t, exp, float64(v))
	case bool:
		return testBooleanLiteral(t, exp, bool(v))
	}
	t.Errorf("Unable to handle Expression. got=%T", exp)
	return false
}

func testIdentifier(t *testing.T, idv ast.Expression, testval string) bool {
	idVal, ok := idv.(*ast.Identifier)
	if !ok {
		t.Errorf("idv not *ast.Identifier. got=%T", idv)
		return false
	}
	if idVal.Value != testval {
		t.Errorf("Identifier.Value not %s. got=%s", testval, idVal.Value)
		return false
	}
	return true
}

func testIntegerLiteral(t *testing.T, il ast.Expression, testval int) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}
	if integ.Value != testval {
		t.Errorf("integerLiteral.Value not %d. got=%d", testval, integ.Value)
		return false
	}
	if integ.Literal != fmt.Sprintf("%d", testval) {
		t.Errorf("integerLiteral.Literal not %d. got=%s", testval, integ.Literal)
		return false
	}
	return true
}

func testFloatLiteral(t *testing.T, fl ast.Expression, testval float64) bool {
	flo, ok := fl.(*ast.FloatLiteral)
	if !ok {
		t.Errorf("flo not *ast.FloatLiteral. got=%T", flo)
		return false
	}
	if flo.Value != testval {
		t.Errorf("FloatLiteral.Value wrong. expected=%f. got=%f", testval, flo.Value)
		return false
	}
	// need to check the difference is bellow a certain error threshold.
	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, testval bool) bool {
	bo, ok := exp.(*ast.BooleanLiteral)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}
	if bo.Value != testval {
		t.Errorf("bo.Value not %t. got=%t", testval, bo.Value)
		return false
	}
	if bo.String() != fmt.Sprintf("%t", testval) {
		t.Errorf("bo.TokenLiteral not %t. got=%s",
			testval, bo.String())
		return false
	}
	return true
}

func TestDictLiteralsStringKeys(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`

	scanner := scanner.Init(input)
	parser := New(scanner)
	program := parser.ParseProgram()
	if errors := parser.Errors(); len(errors) != 0 {
		t.Errorf("parser has %d errors", len(errors))
		for _, err := range errors {
			t.Errorf("parser error: %q", err.message)
		}
		t.FailNow()
	}
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	dict, ok := stmt.Expression.(*ast.DictLiteral)
	if !ok {
		t.Fatalf("expression is not ast.DictLiteral. got=%T", stmt.Expression)
	}

	if len(dict.Items) != 3 {
		t.Errorf("dict.pairs has wrong length, got=%d. want=%d", len(dict.Items), 3)
	}

	expected := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}
	for key, value := range dict.Items {
		lit, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
		}
		testIntegerLiteral(t, value, expected[lit.String()])
	}
}

func TestDictLiteralsStringKeysWithNewlines(t *testing.T) {
	input := `
{
	"one": 1,
	"two": 2,
	"three": 3,
}`

	scanner := scanner.Init(input)
	parser := New(scanner)
	program := parser.ParseProgram()
	if errors := parser.Errors(); len(errors) != 0 {
		t.Errorf("parser has %d errors", len(errors))
		for _, err := range errors {
			t.Errorf("parser error: %q", err.message)
		}
		t.FailNow()
	}
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	dict, ok := stmt.Expression.(*ast.DictLiteral)
	if !ok {
		t.Fatalf("expression is not ast.DictLiteral. got=%T", stmt.Expression)
	}

	if len(dict.Items) != 3 {
		t.Errorf("dict.pairs has wrong length, got=%d. want=%d", len(dict.Items), 3)
	}

	expected := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}
	for key, value := range dict.Items {
		lit, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
		}
		testIntegerLiteral(t, value, expected[lit.String()])
	}
}

func TestDictLiteralsEmpty(t *testing.T) {
	input := `{}`

	scanner := scanner.Init(input)
	parser := New(scanner)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	dict, ok := stmt.Expression.(*ast.DictLiteral)
	if !ok {
		t.Fatalf("expression is not ast.DictLiteral. got=%T", stmt.Expression)
	}
	if len(dict.Items) != 0 {
		t.Errorf("dict.pairs has wrong length, got=%d. want=%d", len(dict.Items), 3)
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	if errors := p.Errors(); len(errors) != 0 {
		t.Errorf("parser has %d errors", len(errors))
		for _, err := range errors {
			t.Errorf("parser error: %q", err.message)
		}
		t.FailNow()
	}
}
