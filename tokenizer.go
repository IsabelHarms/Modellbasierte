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
	t.lineNr = 0
	t.again = false
	if len(t.lines) > 0 { // get first line ready
		t.currentLine = []rune(t.lines[t.lineNr])
		t.position = 0
	}
	t.printLine()
}

func (t *Tokens) unGet() {
	t.again = true
}

func (t *Tokens) Get() Token {

	if t.again {
		t.again = false
		return t.lastToken
	}
	if t.lineNr == len(t.lines) { // behind last line, EOF
		return END
	}
	for { // find an input line with non-white chars
		// skip white space:
		for t.position < len(t.currentLine) && (t.currentLine[t.position] == ' ' || t.currentLine[t.position] == '\t') {
			t.position++
		}
		if t.position < len(t.currentLine) { // found something
			break
		}
		// try get next line:
		t.lineNr++
		if t.lineNr == len(t.lines) {
			return END
		}
		t.printLine()
		t.currentLine = []rune(t.lines[t.lineNr])
		t.position = 0
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
			return t.lastToken
		}

		if unicode.IsDigit(t.currentLine[t.position]) || t.currentLine[t.position] == '-' {
			negative := t.currentLine[t.position] == '-'
			stop := t.position
			if negative {
				stop++ //skip -
				if !unicode.IsDigit(t.currentLine[stop]) {
					t.error("Expected digit after -")
					return INVALID
				}
			}
			for stop < len(t.currentLine) && unicode.IsDigit(t.currentLine[stop]) {
				stop++
			}
			number := t.currentLine[t.position:stop]
			t.lastString = string(number)
			t.lastToken = INTLITERAL
			t.position = stop
			return t.lastToken
		}
		t.error("invalid character")
		t.lastToken = INVALID
	}
	return t.lastToken
}

func (t *Tokens) followingRune(r rune) bool {
	return t.position+1 < len(t.currentLine) && t.currentLine[t.position+1] == r
}

func (t *Tokens) printLine() {
	fmt.Print(t.lineNr + 1)
	fmt.Print(": ")
	fmt.Println(t.lines[t.lineNr])
}

func skipInvalid() {

}

func (t *Tokens) error(s string) {
	t.errorCount++
	fmt.Print(strings.Repeat("-", t.position+2))
	fmt.Println("^")
	fmt.Println(s)
}

func (t *Tokens) warning(s string) {
	fmt.Println(s)
}

func isLetter(c rune) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}
