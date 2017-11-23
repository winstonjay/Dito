package evaluator

import (
	"dito/src/evaluator"
	"dito/src/lexer"
	"dito/src/object"
	"dito/src/parser"
	"testing"
)

func TestNumberParsing(t *testing.T) {
	tests := []struct {
		input   string
		output  interface{}
		outtype string
	}{
		// cases where we expect a Float
		{"-.0", 0.0, object.FloatObj},
		{"0.0", 0.0, object.FloatObj},
		{"0e0", 0.0, object.FloatObj},
		{"0E000", 0.0, object.FloatObj},
		{"14e+4", 140000.0, object.FloatObj},
		{"5e-05", 0.00005, object.FloatObj},
		{"5e-05", 5e-05, object.FloatObj},
		{"10 / 4", 2.5, object.FloatObj},
		{"100 / 3", 100.0 / 3.0, object.FloatObj},
		{"3.0 + 0x7", 10.0, object.FloatObj},
		// cases where we expect a integer
		{"10 / 2", 5, object.IntergerObj},
		{"123231", 123231, object.IntergerObj},
		{"0xfc23", 64547, object.IntergerObj},
		{"0xffff", 65535, object.IntergerObj},
		{"-0x020", -32, object.IntergerObj},
	}
	env := object.NewEnvironment()
	for i, tt := range tests {
		s := lexer.Init("x := " + tt.input)
		p := parser.New(s)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		evaluator.Eval(program, env)

		val, ok := env.Get("x")
		if !ok {
			t.Fatalf("Value error")
		}
		if val.Type() != tt.outtype {
			t.Fatalf("test %d. Wrong type returned. Expected=%s. Got=%s",
				i, tt.outtype, val.Type())
		}
		switch v := val.(type) {
		case *object.Integer:
			if out := int64(tt.output.(int)); v.Value != out {
				t.Fatalf("test %d. Wrong type returned. Expected=%d. Got=%d",
					i, out, v.Value)
			}
		case *object.Float:
			if out := tt.output.(float64); v.Value != out {
				t.Fatalf("test %d. Wrong type returned. Expected=%f. Got=%f",
					i, out, v.Value)
			}
		}

	}
}

func BenchmarkNumberParsing(b *testing.B) {
	env := object.NewEnvironment()
	for i := 0; i < b.N; i++ {
		s := lexer.Init("x := 123231;")
		p := parser.New(s)
		program := p.ParseProgram()
		// checkParserErrors(t, p)
		_ = evaluator.Eval(program, env)
	}
}
