package main

import "fmt"

var tokens Tokens
var varTable VarTable

func main() {
	var testcases [3][2]string
	testcases[0][0] = "{x := -1; x := x<0; print(x)}"
	testcases[0][1] = "true"
	testcases[1][0] = "{x := 0; while x < 5{print(2*x); x = x+1}}"
	testcases[1][1] = "02468"
	testcases[2][0] = "{x := 0; y := true; if ((x == 0)&& y) {print(x)} else {print(y)}}"
	testcases[2][1] = "0"
	tokens = Tokens{} // get tokenizer data
	for i, v := range testcases {
		fmt.Printf("Test No %v: \n", i)
		tokens.setSourceCode(v[0])
		varTable = VarTable{} // get vartable data
		p := program()        // pointer to program  block (root node)
		fmt.Printf("Expected Result: %v \n", v[1])
		fmt.Printf("Produced Result: ")
		if tokens.errorCount == 0 {
			p.exec()
			//fmt.Println("IMP program executed")
		} else {
			fmt.Print(tokens.errorCount)
			fmt.Println(" errors")
		}
	}
}
