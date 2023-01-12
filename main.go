package main

import "fmt"

var tokens Tokens
var varTable VarTable

func main() {
	tokens = Tokens{position: 0, currentLine: []rune(""), errorCount: 0, again: false} // get tokenizer data
	varTable = VarTable{nesting: -1}                                                   // get vartable data
	tokens.setSourceCode("{while true{x := 5;print(4+3*7)}}")

	p := program() // pointer to program  block (root node)
	fmt.Printf("", p)
	/*if tokens.errorCount == 0 {
		p.exec()
		fmt.Println("IMP program executed")
	} else {
		fmt.Print(tokens.errorCount)
		fmt.Println(" errors")
	}*/
}
