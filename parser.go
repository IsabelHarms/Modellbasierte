package main

import (
	"strconv"
)

type Exp interface {
	//pretty() string
	eval() *Value
	GetType() IMPtype
}

// a value is an Expression, it must have:
func (v Value) GetType() IMPtype {
	return v.Type
}

type While struct {
	cond  Exp
	block *Block
}

type IfThenElse struct {
	cond      Exp
	thenBlock *Block
	elseBlock *Block
}

type Print struct {
	exp Exp
}

type Decl struct {
	name       string
	expression Exp
}

type Assign struct {
	name       string
	expression Exp
}

type ExpNode struct {
	op    Token
	Type  IMPtype
	left  Exp
	right Exp // unused for Not-operator
}

type Executable interface { // blocks and statements
	exec()
	// pretty() string
}

// Block and Statements

type Block []Executable // slice of statements

// a ExpNode is an Expression, it must have:
func (e ExpNode) GetType() IMPtype {
	return e.Type
}

// Typabgleich
func SetType(r *ExpNode) {

	r.Type = Undefined // unless we find a legal combination
	if r.left.GetType() == Undefined || r.op != NOT && r.right.GetType() == Undefined {
		return // no additional Error
	}
	switch {
	case r.op == PLUS || r.op == MULT: // arithmetic
		if r.left.GetType() == Integer && r.right.GetType() == Integer {
			r.Type = Integer
		}
	case r.op == AND || r.op == OR: // boolean binary
		if r.left.GetType() == Boolean && r.right.GetType() == Boolean {
			r.Type = Boolean
		}
	case r.op == NOT: // boolean unary
		if r.left.GetType() == Boolean {
			r.Type = Boolean
		}
	case r.op == EQUAL:
		if r.left.GetType() == r.right.GetType() {
			r.Type = Boolean
		} else {
			tokens.warning("Operands of '==' must have same type")
		}
	case r.op == LESS:
		if r.left.GetType() == Integer && r.right.GetType() == Integer {
			r.Type = Boolean
		}
	}
	if r.Type == Undefined {
		tokens.error(" boolean / arithmetic type mismatch")
	}
}

func program() Executable {
	b := block()
	if tokens.Get() != END {
		tokens.error("trailing garbage")
	}
	return b
}

func block() Block {
	if tokens.Get() != BLOCKSTART {
		tokens.error("expected { block }")
	}
	varTable.blockStart()
	b := make(Block, 0)
	for { // loop while finding statements
		s := statement()
		b = append(b, s)
		switch tokens.Get() {
		case SEPARATOR:
			continue
		case BLOCKSTOP:
			{
				varTable.blockEnd()
				/*if tokens.Get() != END {
					tokens.error("missing '{' or garbage after program block")
				}*/
				return b
			}
		default:
			tokens.unGet()
			tokens.error("missing ';' or '}' after statement")
			// todo: skip garbage until next statement ?
			return b
		}
	}
}

func statement() Executable {
	switch tokens.Get() {
	case IF:
		return ifStatement()
	case WHILE:
		return whileStatement()
	case PRINT:
		return printStatement()
	case NAME:
		name := tokens.lastString
		switch tokens.Get() {
		case ASSIGN:
			return assignment(name)
		case DECLARE:
			return declaration(name)
		default:
			tokens.unGet()
			tokens.error("assignment or declaration expected")
			return nil
		}
	}
	tokens.unGet()
	tokens.error("statement expected")
	return nil
}

func ifStatement() Executable {
	e := expression()
	if e.GetType() == Integer {
		tokens.error("if condition must be boolean")
	}
	thenPart := block()
	if tokens.Get() != ELSE {
		tokens.unGet()
		tokens.error("if without 'else' (not allowed in IMP)")
	}
	elsePart := block()
	return IfThenElse{cond: e, thenBlock: &thenPart, elseBlock: &elsePart}
}

func whileStatement() Executable {
	e := expression()
	if e.GetType() == Integer {
		tokens.error("while condition must be boolean")
	}
	b := block()
	return While{cond: e, block: &b}
}

func printStatement() Executable {
	return Print{expression()}
}

func assignment(name string) Executable {
	v := varTable.Get(name)
	if v == nil {
		tokens.error("undefined variable")
		return nil
	}
	exp := expression()
	if v.Type == Undefined || exp.GetType() == Undefined {
		return nil
	}
	if v.Type == exp.GetType() {
		return Assign{name, exp}
	}
	tokens.error("type mismatch")
	return nil
}

func declaration(name string) Executable {
	exp := expression()
	switch exp.GetType() {
	case Integer:
		varTable.declareInt(name, 0) //Pseudo value
	case Boolean:
		varTable.declareBool(name, false) //Pseudo value
	default:
		varTable.declareUndefined(name)
	}
	return Decl{name, exp}
}

/*
Grammatik
Expression := OrOp | OrOp "| |" Expression
OrOp := AndOp | AndOp "&&" OrOp
AndOp := EqualOp | EqualOp "==" EqualOp
EqualOp := LessOp | LessOp "<" LessOp
LessOp := PlusOp | PlusOp "+" LessOp
PlusOp := Operand | Operand "*" PlusOp
Operand := Variable | Literal | "!" Operand |"(" Expression ")"
*/
func expression() Exp {
	lhs := orOp()
	if tokens.Get() == OR {
		node := ExpNode{op: OR, left: lhs, right: expression()}
		SetType(&node)
		return &node
	}
	tokens.unGet()
	return lhs
}

func orOp() Exp {
	lhs := andOp()
	if tokens.Get() == PLUS {
		node := ExpNode{op: AND, left: lhs, right: orOp()}
		SetType(&node)
		return &node
	}
	tokens.unGet()
	return lhs
}

func andOp() Exp {
	lhs := equalOp()
	if tokens.Get() == EQUAL {
		node := ExpNode{op: EQUAL, left: lhs, right: equalOp()}
		SetType(&node)
		return &node
	}
	tokens.unGet()
	return lhs
}

func equalOp() Exp {
	lhs := lessOp()
	if tokens.Get() == LESS {
		node := ExpNode{op: LESS, left: lhs, right: lessOp()}
		SetType(&node)
		return &node
	}
	tokens.unGet()
	return lhs
}

func lessOp() Exp {
	lhs := plusOp()
	if tokens.Get() == PLUS {
		node := ExpNode{op: PLUS, left: lhs, right: lessOp()}
		SetType(&node)
		return &node
	}
	tokens.unGet()
	return lhs
}

func plusOp() Exp {
	lhs := operand()
	if tokens.Get() == MULT {
		node := ExpNode{op: MULT, left: lhs, right: plusOp()}
		SetType(&node)
		return &node
	}
	tokens.unGet()
	return lhs
}

func operand() Exp {
	switch tokens.Get() {
	case NAME:
		valPtr := varTable.Get(tokens.lastString)
		if valPtr == nil {
			tokens.error("undefined variable")
			return &Value{Type: Undefined}
		}
		return valPtr
	case BOOLLITERAL:
		return &Value{Type: Boolean, bValue: tokens.lastString == "true"}
	case INTLITERAL:
		num, _ := strconv.Atoi(tokens.lastString)
		return &Value{Type: Integer, iValue: num}
	case NOT:
		node := ExpNode{op: NOT, left: operand()}
		SetType(&node)
		return &node
	case OPEN:
		exp := expression()
		if tokens.Get() != CLOSE {
			tokens.error("lacking ')'")
			//todo
		}
		return exp
	}
	tokens.error("missing operand")
	tokens.unGet()
	return &Value{Type: Undefined}
}
