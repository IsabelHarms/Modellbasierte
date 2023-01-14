package main

import "fmt"

var tokens Tokens
var varTable VarTable

func main() {
	var testcases [21][2]string
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
	testcases[16][0] = "{x := 5; y := 4; z := x*y; if (z == 20) {print(z)} else {print(y)}}"
	testcases[16][1] = "20"
	testcases[17][0] = "{a := 0; b := 15; z := true; while !(a == b) {if(z) {a = a+1} else {b = b + -1}; z = !z}; print(a)}"
	testcases[17][1] = "8"
	testcases[18][0] = "{a := 0; b := 6;  c := 5;while a < b {a = a+1; c = 5; while c < 20 {c= c+5; if (50< c*a) {print(c*a)} else {c=c}}}}"
	testcases[18][1] = "60 60 80 75 100 60 90 120"
	testcases[19][0] = "{a := false; b := false;  c := false; while (((a == false)|| (b==false)) || c==false) {if(c){c = false; if(b){b = false; if(a){a = false} else {a = true; print(a)}} else  {b = true}} else {c = true}}}"
	testcases[19][1] = "true"
	testcases[20][0] = "{a := 1; b := 1; c := 0; while a < 5 {b=1; while b < 5 {if(a == b) {c = c + (-1* (a*b))} else {c = c + (a*b)}; b = b+1}; a = a+1}; print(c)}"
	testcases[20][1] = "40"
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
