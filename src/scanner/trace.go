package scanner

// TraceLine : Returns last line up to current column.
// eg. at index 8 of "alpha := 100" we would get: 'alpha :=' <-.
func (s *Scanner) TraceLine() string {
	return s.input[s.linePos : s.linePos+s.column]
}

// // PrintScan method prints out a table of the entire lexical
// // analysis of an input in one go.
// func (s *Scanner) PrintScan() {
// 	tok, literal, _ := s.NextToken()
// 	tokenCount := 0
// 	fmt.Printf("input:\n\n%s\n\n", s.input)
// 	fmt.Printf("| line | col  | Token        | Literal     |\n")
// 	fmt.Printf("-----------------------------------------\n")
// 	for tok != token.EOF {
// 		fmt.Printf("| %4d | %4d | %12s | %12s |\n",
// 			s.lineno+1, s.column-len(literal), tok.String(), literal)
// 		tokenCount++
// 		tok, literal, _ = s.NextToken()
// 	}
// 	fmt.Printf("\nTotal Tokens: %d, \n", tokenCount)
// 	fmt.Printf("Total Chars: %d, \n", s.pos)
// 	fmt.Printf("Total Lines: %d, \n", s.lineno+1)
// }
