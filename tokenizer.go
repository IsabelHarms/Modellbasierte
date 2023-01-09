package main

import (
	"fmt"
	"strings"
	"unicode"
)

type Token int

type Tokens struct { //infos on the current code
	position    int
	lines       []string
	lineNr      int
	currentLine []rune
	errorCount  int
	lastString  string
	lastToken   Token
	again       bool
}

const (
	INVALID     Token = 0
	WHILE       Token = 1
	IF          Token = 2
	ELSE        Token = 3
	PRINT       Token = 4
	BOOLLITERAL Token = 5
	INTLITERAL  Token = 6
	OPEN        Token = 7
	CLOSE       Token = 8
	BLOCKSTART  Token = 9
	BLOCKSTOP   Token = 10
	DECLARE     Token = 11
	ASSIGN      Token = 12
	PLUS        Token = 13
	MULT        Token = 14
	AND         Token = 15
	OR          Token = 16
	NOT         Token = 17
	LESS        Token = 18
	EQUAL       Token = 19
	SEPERATOR   Token = 20
	NAME        Token = 21
	END         Token = 22
)

func (t Tokens) setSourceCode(source string) {
	t.lines = strings.Split(source, "\n")
}

func (t Tokens) unGetToken() {
	t.again = true
}

func (t Tokens) getToken() Token {

	if t.again {
		t.again = false
		return t.lastToken
	}

	for t.position < len(t.currentLine) && (t.currentLine[t.position] == ' ' || t.currentLine[t.position] == '\t') {
		t.position++
	}

	if t.position == len(t.currentLine) { //get next line
		t.currentLine = t.nextLine()
		t.position = 0
	}
	if len(t.currentLine) == 0 { //end of code
		t.lastToken = END
	}

	//todo nicht return DECLARE, sondern lastToken = DECLARE und am ende return lastToken
	switch t.currentLine[t.position] {
	case ':':
		if t.followingRune('=') {
			t.position += 2
			t.lastToken = DECLARE
		} else {
			t.position++
			t.error("invalid operand")
			skipInvalid()
			t.lastToken = INVALID
		}
	case '+':
		t.position++
		t.lastToken = PLUS
	case '*':
		t.position++
		t.lastToken = MULT
	case '!':
		t.position++
		t.lastToken = NOT
	case '<':
		t.position++
		t.lastToken = LESS
	case ',':
		t.position++
		t.lastToken = SEPERATOR
	case '=':
		if t.followingRune('=') {
			t.position += 2
			t.lastToken = EQUAL
		} else {
			t.position++
			t.lastToken = ASSIGN
		}
	case '(':
		t.position++
		t.lastToken = OPEN
	case ')':
		t.position++
		t.lastToken = CLOSE
	case '{':
		t.position++
		t.lastToken = BLOCKSTART
	case '}':
		t.position++
		t.lastToken = BLOCKSTOP
	case '-':
		t.position++
		if unicode.IsDigit(t.currentLine[t.position+1]) {
			t.getTokenNumber()
		}
		t.lastToken = INVALID
	}

	if unicode.IsLetter(t.currentLine[t.position]) {
		var stop int
		for stop := t.position + 1; stop < len(t.currentLine) && (unicode.IsLetter(t.currentLine[stop]) || unicode.IsDigit(t.currentLine[stop])); stop++ {

		}
		word := t.currentLine[t.position:stop]
		t.lastString = string(word)
		switch string(t.lastString) {
		case "WHILE":
			t.lastToken = WHILE
		case "IF":
			t.lastToken = IF
		case "ELSE":
			t.lastToken = ELSE
		case "TRUE":
			t.lastToken = BOOLLITERAL
		case "FALSE":
			t.lastToken = BOOLLITERAL
		case "AND":
			t.lastToken = AND
		case "OR":
			t.lastToken = OR
		case "PRINT":
			t.lastToken = PRINT
		default:
			t.lastToken = NAME //if none of the cases are true it has to be a NAME
		}
	}

	if unicode.IsDigit(t.currentLine[t.position]) {
		t.getTokenNumber()
	}
	//fmt.Println("%v\n", t.lastToken)
	return t.lastToken
}

func (t Tokens) getTokenNumber() {
	stop := t.position + 1
	for unicode.IsDigit(t.currentLine[stop]) {
		stop++
	}
	number := t.currentLine[t.position:stop]
	t.lastString = string(number)
	t.lastToken = INTLITERAL
}

func (t Tokens) nextLine() []rune {
	t.lineNr++
	//fmt.Println(t.lines[t.lineNr])
	return []rune(t.lines[t.lineNr])
}

func (t Tokens) followingRune(r rune) bool {
	return (t.position+1 < len(t.currentLine) && t.currentLine[t.position+1] == r)
}

func skipInvalid() {

}

func (t Tokens) error(s string) {
	t.errorCount++
	fmt.Println(s)
}

func (t Tokens) warning(s string) {
	fmt.Println(s)
}

func main() {
	tokens = Tokens{position: 0, currentLine: []rune(""), errorCount: 0, again: false}
	tokens.setSourceCode("3+4*5")
	tokens.getToken()
	println(tokens.lastString)
}
