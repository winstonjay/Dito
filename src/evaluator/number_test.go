package evaluator

import (
	"dito/src/evaluator"
	"dito/src/lexer"
	"dito/src/object"
	"dito/src/parser"
	"fmt"
	"testing"
)

func TestNumberParsing(t *testing.T) {
	tests := []struct {
		input  string
		output interface{}
	}{
		{"x := 123231;", 123231},
		{"x := 14e+4;", 140000},
		{"x := 5e-05;", 0.00005},
		{"x := 5e-05;", 0.00005},
		{"x := 0xfc23;", 64547},
		{"x := 0xffff;", 65535},
		{"x := -0x020;", -32},
	}
	env := object.NewEnvironment()
	for _, tt := range tests {
		s := lexer.Init(tt.input)
		p := parser.New(s)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		evaluator.Eval(program, env)
		if val, ok := env.Get("x"); ok {
			fmt.Printf("%v, ", val.Inspect())
		}
	}
}
