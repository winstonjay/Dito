package eval

// TODO:
// Need to review these tests. At a glance some of them seem a bit odd
// or incorect.

import (
	"dito/src/object"
	"dito/src/parser"
	"dito/src/scanner"
	"testing"
)

// TestIntExpr :
func TestIntExpr(t *testing.T) {
	tests := []struct {
		input    string
		expected int
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
		{"0x14 + 0x00000002 * -0X000A", 0x000}, // what is going on here.
		{"0X32 // 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37}, {"((3 * (3 * 3)) + 10)", 37},
		{"(5 + 10 * 2 + 15 // 3) * 2 + -10", 50},
		{"(0XF-0xa)**(0X2c-0x29) // 5", 25},
		{"0xFFFFCDDE21 * 2", 2199016684610},
		{"0x0001 << 0x003f - 1", 9223372036854775807},
		{"1 << 63 - 1", 9223372036854775807},
		{"1 << 63", -9223372036854775808},
		{"9223372036854775807 >> 63", 0},
		{"9223372036854775807 >> 32", 2147483647},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testInt(t, evaluated, tt.expected)
	}
}

func TestEvalIfElseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"35 if 35 <= 1 else 35 - 2", 33},
		{"(35 / 2) if 35 % 2 == 0 else 35 * 3 + 1", 106},
		{"35 / 2 if 35 % 2 == 0 else 35 * 3 + 1", 106},
	}
	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testInt(t, evaluated, tt.expected)
	}
}

func testInt(t *testing.T, obj object.Object, expected int) bool {
	result, ok := obj.(*object.Int)
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

func testEval(t *testing.T, input string) object.Object {
	l := scanner.Init(input)
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

func TestDictEvaluation(t *testing.T) {
	input := `
let a = 2
{
	"one": 10 - 9,
	"two": a,
	"thr" ++ "ee": 3,
	4: 4,
	true: 5,
	false: 6,
}
`
	evaluated := testEval(t, input)
	result, ok := evaluated.(*object.Dict)
	if !ok {
		t.Fatalf("eval didn't produce a dict. got=%T", evaluated)
	}

	expected := map[object.HashKey]int{
		(&object.String{Value: "one"}).Hash():   1,
		(&object.String{Value: "two"}).Hash():   2,
		(&object.String{Value: "three"}).Hash(): 3,
		(&object.Int{Value: 4}).Hash():          4,
		(object.TRUE).Hash():                    5,
		(object.FALSE).Hash():                   6,
	}

	if len(result.Items) != len(expected) {
		t.Fatalf("dict has the wrong number of items. got=%d, want=%d",
			len(result.Items), len(expected))
	}
	if result.Len != len(expected) {
		t.Fatalf("dict has the wrong recorded number of items. got=%d, want=%d",
			len(result.Items), len(expected))
	}

	for expectedKey, expectedValue := range expected {
		item, ok := result.Items[expectedKey]
		if !ok {
			t.Errorf("no item for given key")
		}
		testInt(t, item.Value, expectedValue)
	}

}
