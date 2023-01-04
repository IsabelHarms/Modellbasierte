package main

var tokens Tokens
var vartable *VarTable
var invalidExp Exp

func parse() {

	tokens = Tokens{position: 0, sourceCode: "A+B*C", currentLine: []rune(""), errorCount: 0, again: false}

	vartable = &VarTable{nesting: -1} // ready for first start of block
	invalidExp = (Var)
	// work with a pointer, otherwise you will get multiple structs!

	vartable.blockStart() // dereferencing is implied! :-(

	vartable.declareInt("Isabel", 22)
	vartable.declareInt("David", 25)
	vartable.blockStart()
	vartable.declareInt("Isabel", 2000)
	vartable.Get("Isabel").showType()
	vartable.blockEnd()
	vartable.Get("Isabel").printValue()
	vartable.blockStart()
	vartable.Get("Isabel").printValue()
}

/*type IMPtype byte

const (
	BOOLEAN   IMPtype = 0
	INTEGER   IMPtype = 1
	UNDEFINED IMPtype = 2 // when parsing illegal expressions
)

type Variable struct {
	varType IMPtype
	iValue  int
	bValue  bool
}

func (v *Variable) printValue() {
	switch {
	case v == nil:
		fmt.Println("undeclared Variable")
	case v.varType == BOOLEAN:
		fmt.Println(v.bValue)
	case v.varType == INTEGER:
		fmt.Println(v.iValue)
	default:
		fmt.Println("undefined Type")
	}
}

const MaxNesting = 50

type VarTable struct {
	nesting int
	names   [MaxNesting]map[string]*Variable // array of maps, each is still = nil
}*/

func (vt *VarTable) blockStart() {
	vt.nesting++

	if vt.nesting == MaxNesting {
	} // toDo: runtime error
}

func (vt *VarTable) blockEnd() {
	if vt.nesting < 0 {
	} // toDo: runtime error
	vt.names[vt.nesting] = nil // clear all declarations
	vt.nesting--
}

func (vt *VarTable) declareInt(name string, i int) {
	if vt.names[vt.nesting] == nil {
		vt.names[vt.nesting] = map[string]*Val{}
	} // new map for current block
	vt.names[vt.nesting][name] = &Val{flag: ValueInt, valI: i}
}

func (vt *VarTable) Get(name string) *Val {
	for i := vt.nesting; i >= 0; i-- {
		if vt.names[vt.nesting] == nil {
			continue
		}
		v := vt.names[vt.nesting][name]
		if v != nil {
			return v
		}
	}
	return nil
}

/*
Operand := Variable | Literal | "!" Operand |"(" Expression ")"
PlusOp := Operand | Operand "*" PlusOp
LessOp := PlusOp | PlusOp "+" LessOp
EqualOp := LessOp | LessOp "<" LessOp
AndOp := EqualOp | EqualOp "==" EqualOp
OrOp := AndOp | AndOp "&&" OrOp
Expression := OrOp | OrOp "| |" Expression
*/

func operand() Exp {
	switch tokens.getToken() {
	case NAME:
		varPtr := vartable.Get(tokens.lastString)
		if varPtr == nil {
			tokens.error("undefined variable")
			//return
		}
		//return (Exp)(*varPtr)

	case BOOLLITERAL:
		return Bool(tokens.lastString == "true")
	case INTLITERAL:
		//todo return Num()
	case NOT:
		//op := operand()
		//if op.infer() typ checken und negieren
	case OPEN:
		exp := expression()
		if tokens.getToken() != CLOSE {
			tokens.error("lacking ')'")
			//todo
		}
		return exp
	}
}

func plusOp() Exp {
	lhs := operand()
	if tokens.getToken() == MULT {
		return (Mult)([2]Exp{lhs, plusOp()})
	}
	tokens.unGetToken()
	return lhs
}

func lessOp() Exp {
	lhs := plusOp()
	if tokens.getToken() == PLUS {
		return (Plus)([2]Exp{lhs, lessOp()})
	}
	tokens.unGetToken()
	return lhs
}

func equalOp() Exp {
	lhs := lessOp()
	if tokens.getToken() == LESS {
		return (Less)([2]Exp{lhs, lessOp()})
	}
	tokens.unGetToken()
	return lhs
}

func andOp() Exp {
	lhs := equalOp()
	if tokens.getToken() == EQUAL {
		return (Equal)([2]Exp{lhs, equalOp()})
	}
	tokens.unGetToken()
	return lhs
}

func orOp() Exp {
	lhs := andOp()
	if tokens.getToken() == PLUS {
		return (And)([2]Exp{lhs, orOp()})
	}
	tokens.unGetToken()
	return lhs
}

func expression() Exp {
	lhs := orOp()
	if tokens.getToken() == OR {
		return (Or)([2]Exp{lhs, expression()})
	}
	tokens.unGetToken()
	return lhs
}
