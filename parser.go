package main

var tokens Tokens
var vartable *VarTable
var invalidExp Exp

func parse() {

	tokens = Tokens{position: 0, currentLine: []rune(""), errorCount: 0, again: false}
	tokens.setSourceCode("print(3+4*5)")
	tokens.getToken()
	vartable = &VarTable{nesting: -1} // ready for first start of block
	//invalidExp = (Var)
	// work with a pointer, otherwise you will get multiple structs!

	/*vartable.blockStart() // dereferencing is implied! :-(

	vartable.declareInt("Isabel", 22)
	vartable.declareInt("David", 25)
	vartable.blockStart()
	vartable.declareInt("Isabel", 2000)
	vartable.Get("Isabel").showType()
	vartable.blockEnd()
	vartable.Get("Isabel").printValue()
	vartable.blockStart()
	vartable.Get("Isabel").printValue()*/
}

type Exp interface {
	//pretty() string
	//eval() Value
	GetType() IMPtype
}

// a value is an Expression, it must have:
func (v Value) GetType() IMPtype {
	return v.Type
}

type ExpNode struct {
	op    Token
	Type  IMPtype
	left  Exp
	right Exp // unused for Not-operator
}

// a ExpNode is an Expression, it must have:
func (r ExpNode) GetType() IMPtype {
	return r.Type
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

func fsth() {
	parse()
	/*r := ExpNode{op: PLUS, left: Value{Type: Integer, iValue: 3}, right: Value{Type: Integer, iValue: 2}}
	SetType(&r)
	fmt.Println(r.GetType())

	r = ExpNode{op: LESS, left: Value{Type: Integer, iValue: 3}, right: Value{Type: Integer, iValue: 2}}
	SetType(&r)
	fmt.Println(r.GetType())

	r = ExpNode{op: EQUAL, left: Value{Type: Integer, iValue: 3}, right: Value{Type: Boolean, bValue: false}}
	SetType(&r)
	fmt.Println(r.GetType())*/
}

/*
Grammatik
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
		valPtr := vartable.Get(tokens.lastString)
		if valPtr == nil {
			tokens.error("undefined variable")
			return &Value{Type: Undefined}
		}
		return valPtr

	case BOOLLITERAL:
		return &Value{Type: Boolean, bValue: tokens.lastString == "true"}
	case INTLITERAL:
		//todo return Num()
	case NOT:
		node := ExpNode{op: NOT, left: operand()}
		SetType(&node)
		return &node
	case OPEN:
		exp := expression()
		if tokens.getToken() != CLOSE {
			tokens.error("lacking ')'")
			//todo
		}
		return exp
	}
	tokens.error("missing operand")
	tokens.unGetToken()
	return &Value{Type: Undefined}
}

func plusOp() Exp {
	lhs := operand()
	if tokens.getToken() == MULT {
		node := ExpNode{op: MULT, left: lhs, right: plusOp()}
		SetType(&node)
		return &node
	}
	tokens.unGetToken()
	return lhs
}

func lessOp() Exp {
	lhs := plusOp()
	if tokens.getToken() == PLUS {
		node := ExpNode{op: PLUS, left: lhs, right: lessOp()}
		SetType(&node)
		return &node
	}
	tokens.unGetToken()
	return lhs
}

func equalOp() Exp {
	lhs := lessOp()
	if tokens.getToken() == LESS {
		node := ExpNode{op: LESS, left: lhs, right: lessOp()}
		SetType(&node)
		return &node
	}
	tokens.unGetToken()
	return lhs
}

func andOp() Exp {
	lhs := equalOp()
	if tokens.getToken() == EQUAL {
		node := ExpNode{op: EQUAL, left: lhs, right: equalOp()}
		SetType(&node)
		return &node
	}
	tokens.unGetToken()
	return lhs
}

func orOp() Exp {
	lhs := andOp()
	if tokens.getToken() == PLUS {
		node := ExpNode{op: AND, left: lhs, right: orOp()}
		SetType(&node)
		return &node
	}
	tokens.unGetToken()
	return lhs
}

func expression() Exp {
	lhs := orOp()
	if tokens.getToken() == OR {
		node := ExpNode{op: OR, left: lhs, right: expression()}
		SetType(&node)
		return &node
	}
	tokens.unGetToken()
	return lhs
}
