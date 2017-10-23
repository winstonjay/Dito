package main

import "dito/token"
import "fmt"

/*
This file collects all easily assible information about
the programing language and formats it into a readable
documentation group.
*/

func readTokens() {

	for key := range token.Keywords {
		fmt.Println(key)
	}
}

func main() {
	readTokens()
}
