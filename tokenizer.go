package main

import (
	"fmt"
	"strings"
	"unicode"
)

type Token byte

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
	SEPARATOR   Token = 20
	NAME        Token = 21
	END         Token = 22
)

func (t *Tokens) setSourceCode(source string) {
	t.lines = strings.Split(source, "\n")
}

func (t *Tokens) unGet() {
	t.again = true
}

func (t *Tokens) Get() Token {

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
	case '&':
		if t.followingRune('&') {
			t.position += 2
			t.lastToken = AND
		} else {
			t.lastToken = INVALID
		}
	case '|':
		if t.followingRune('|') {
			t.position += 2
			t.lastToken = OR
		} else {
			t.lastToken = INVALID
		}
	case '!':
		t.position++
		t.lastToken = NOT
	case '<':
		t.position++
		t.lastToken = LESS
	case ';':
		t.position++
		t.lastToken = SEPARATOR
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
	default:
		if isLetter(t.currentLine[t.position]) { //is letter is wrong
			stop := t.position + 1
			for stop < len(t.currentLine) && (isLetter(t.currentLine[stop]) || unicode.IsDigit(t.currentLine[stop])) {
				stop++
			}
			word := t.currentLine[t.position:stop]
			t.position = stop
			t.lastString = string(word)
			//toLower?
			switch t.lastString {
			case "while":
				t.lastToken = WHILE
			case "if":
				t.lastToken = IF
			case "else":
				t.lastToken = ELSE
			case "true":
				t.lastToken = BOOLLITERAL
			case "false":
				t.lastToken = BOOLLITERAL
			case "print":
				t.lastToken = PRINT
			default:
				t.lastToken = NAME //if none of the cases are true it has to be a NAME
			}
		}

		if unicode.IsDigit(t.currentLine[t.position]) {
			t.getTokenNumber()
		}
	}

	fmt.Printf("%v\n", t.lastToken)
	return t.lastToken
}

func (t *Tokens) getTokenNumber() {
	stop := t.position
	for stop < len(t.currentLine) && unicode.IsDigit(t.currentLine[stop]) {
		stop++
	}
	number := t.currentLine[t.position:stop]
	t.lastString = string(number)
	t.lastToken = INTLITERAL
	t.position = stop
}

func (t *Tokens) nextLine() []rune {
	t.lineNr++
	if t.lineNr > len(t.lines) {
		return []rune(" ")
	}
	return []rune(t.lines[t.lineNr-1]) //todo
}

func (t *Tokens) followingRune(r rune) bool {
	return t.position+1 < len(t.currentLine) && t.currentLine[t.position+1] == r
}

func skipInvalid() {

}

func (t *Tokens) error(s string) {
	t.errorCount++
	fmt.Println(s)
}

func (t *Tokens) warning(s string) {
	fmt.Println(s)
}

func isLetter(c rune) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}
