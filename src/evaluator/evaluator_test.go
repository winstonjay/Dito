package evaluator

import (
	"dito/src/lexer"
	"dito/src/object"
	"dito/src/parser"
	"testing"
)

func testEval(t *testing.T, input string) object.Object {
	l := lexer.Init(input + ";")
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	env := object.NewEnvironment()
	return Eval(program, env)
}

func checkParserErrors(t *testing.T, p *parser.Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"0", 0}, {"-0", 0}, {"0x0", 0}, {"-0x0", 0},
		{"1", 1}, {"-2", -2}, {"-(3)", -3}, {"-(-(-(4)))", -4},
		{"-6", -6}, {"----7", 7}, {"-0xA", -10}, {"--0xA", 10},
		{"0x7fffffffffffffff", 9223372036854775807},
		{"((((((1000000000000001))))))", 1000000000000001},
		{"---9223372036854775807", -9223372036854775807},
		{"-0xA**(0xB % 3)", 100},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"0x002 * 2 * 0x002 * 0x002 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"0x20 * 2 + 10", 74},
		{"5 + 2 * 10", 25}, {"(5 + (2 * 10))", 25},
		{"0x5a + (0xa + 0xb + 0xc)", 123},
		{"0x14 + 0x00000002 * -0X000A", 0x000},
		{"0X32 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37}, {"((3 * (3 * 3)) + 10)", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
		{"(0XF-0xa)**(0X2c-0x29) / 5", 25},
		{"0xFFFFCDDE21 * 2", 2199016684610},
		{"0x0001 << 0x003f - 1", 9223372036854775807},
		{"1 << 63 - 1", 9223372036854775807},
		{"1 << 63", -9223372036854775808},
		{"9223372036854775807 >> 63", 0},
		{"9223372036854775807 >> 32", 2147483647},
		// {"2 ** 63", gives bad overflow error }
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
		return false
	}
	return true
}
