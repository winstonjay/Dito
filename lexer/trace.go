package lexer

import (
	"dito/token"
	"fmt"
)

// PrintScan : print out the entire lexical analysis of an input
// in one go.
func (s *Scanner) PrintScan() {

	tok, literal := s.NextToken()
	tokenCount := 0
	fmt.Printf("Input:\n\n%s\n\n", s.input)
	fmt.Printf("| line | col  | Token        | Literal     |\n")
	fmt.Printf("-----------------------------------------\n")
	for tok != token.EOF {
		fmt.Printf("| %4d | %4d | %12s | %12s |\n",
			s.Lineno+1, s.Column-len(literal), tok.String(), literal)
		tokenCount++
		tok, literal = s.NextToken()
	}
	fmt.Printf("\nTotal Tokens: %d, \n", tokenCount)
	fmt.Printf("Total Chars: %d, \n", s.pos)
	fmt.Printf("Total Lines: %d, \n", s.Lineno+1)
}

// TraceLine : Returns last line up to current column.
// eg. at index 8 of "alpha := 100" we would get: 'alpha :=' <-.
func (s *Scanner) TraceLine() string {
	return s.input[s.linePos : s.linePos+s.Column]
}
