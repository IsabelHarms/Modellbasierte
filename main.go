package main

import "fmt"

var tokens Tokens
var varTable VarTable

func main() {
	var testcases [16][2]string
	testcases[0][0] = "{x := -1; x := x<0; print(x)}"
	testcases[0][1] = "true"
	testcases[1][0] = "{x := 0; while x < 5{print(2*x); x = x+1}}"
	testcases[1][1] = "02468"
	testcases[2][0] = "{x := 0; y := true; if ((x == 0)&& y) {print(x)} else {print(y)}}"
	testcases[2][1] = "0"
	testcases[3][0] = "{var := 1; while (var < 4){var = var +1; print(var)}}"
	testcases[3][1] = "2 3 4"
	testcases[4][0] = "{x := !false; if(x && true){print(x)} else{print(0)}}"
	testcases[4][1] = "true"
	testcases[5][0] = "{x := !true; if(x && true){print(x)} else{print(test)}}"
	testcases[5][1] = "undefined variable"
	testcases[6][0] = "{x := -01; print(x)}"
	testcases[6][1] = "-1"
	testcases[7][0] = "{x := 1; y:= 2; x=x+y; print(x)}"
	testcases[7][1] = "3"
	testcases[8][0] = "{x := true; y := false; if(y || true) {if(x && true){print(1)}else{print(2)}} else {print(3)}}"
	testcases[8][1] = "1"
	testcases[9][0] = "{x := true ; y := 2 ; print(x*y)}"
	testcases[9][1] = "error"
	testcases[10][0] = "{x := true ; y := 2 ; print(x+y)}"
	testcases[10][1] = "error"
	testcases[11][0] = "{x := true ; y := 2 ; print(x<y)}"
	testcases[11][1] = "error"
	testcases[12][0] = "{x := true ; y := 2 ; print(x && y)}"
	testcases[12][1] = "error"
	testcases[13][0] = "{x := true ; y := 2 ; print(x || y)}"
	testcases[13][1] = "error"
	testcases[14][0] = "{x := 1 ; while(x){print(2)}}"
	testcases[14][1] = "error"
	testcases[15][0] = "{x := 3 ; y := 2 ; if ( 2 == (x || y )){print(10)}else{print(15)}}"
	testcases[15][1] = "error" //hoffe ich
	for i, v := range testcases {
		tokens = Tokens{} // get tokenizer data
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
