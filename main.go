package main

import "fmt"

var tokens Tokens
var varTable VarTable

func main() {
	tokens = Tokens{} // get tokenizer data
	tokens.setSourceCode("{print(5+3)}")
	varTable = VarTable{nesting: -1} // get vartable data

	p := program() // pointer to program  block (root node)
	fmt.Printf("", p)
	if tokens.errorCount == 0 {
		p.exec()
		fmt.Println("IMP program executed")
	} else {
		fmt.Print(tokens.errorCount)
		fmt.Println(" errors")
	}
}
