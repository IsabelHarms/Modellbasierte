package main

import "fmt"

// executes
func (b Block) exec() {
	for _, stmt := range b {
		stmt.exec()
	}
}
func (w While) exec() {
	for w.cond.eval().bValue {
		w.block.exec()
	}
}

func (i IfThenElse) exec() {
	if i.cond.eval().bValue {
		i.thenBlock.exec()
	} else {
		i.elseBlock.exec()
	}
}

func (p Print) exec() {
	v := p.exp.eval()
	if v.Type == Integer {
		fmt.Println(v.iValue)
	} else {
		fmt.Println(v.bValue)
	}
}

func (d Decl) exec() {
	name := d.name
	value := d.expression
	if value.GetType() == Integer {
		varTable.declareInt(name, value.eval().iValue)
	} else {
		varTable.declareBool(name, value.eval().bValue)
	}
}
func (a Assign) exec() {
	name := a.name
	value := a.expression
	if value.GetType() == Integer {
		varTable.declareInt(name, value.eval().iValue)
	} else {
		varTable.declareBool(name, value.eval().bValue)
	}
}

// evals
func (e ExpNode) eval() *Value {
	//switch plus, mal
	switch e.op {
	case PLUS:
		return &Value{Type: Integer, iValue: e.left.eval().iValue + e.right.eval().iValue}
	case MULT:
		return &Value{Type: Integer, iValue: e.left.eval().iValue * e.right.eval().iValue}
	case EQUAL:
		if e.left.GetType() == Integer {
			return &Value{Type: Boolean, bValue: e.left.eval().iValue == e.right.eval().iValue}
		}
		return &Value{Type: Boolean, bValue: e.left.eval().bValue == e.right.eval().bValue}
	case LESS:
		return &Value{Type: Boolean, bValue: e.left.eval().iValue < e.right.eval().iValue}
	case AND:
		return &Value{Type: Boolean, bValue: e.left.eval().bValue && e.right.eval().bValue}
	case OR:
		return &Value{Type: Boolean, bValue: e.left.eval().bValue || e.right.eval().bValue}
	case NOT:
		return &Value{Type: Boolean, bValue: !e.left.eval().bValue}
	case NAME:
		//return
	}
	return nil //todo this should never happen
}

func (v Variable) eval() *Value {
	return varTable.Get(v.name)
}

func (v Value) eval() *Value {
	return &v
}

// prettys
func (b Block) pretty() string {
	string := "{"
	for _, stmt := range b {
		string += stmt.pretty() + ";"
	}
	return string + "}"
}

func (w While) pretty() string {
	return "while" + w.cond.pretty() + "{" + w.block.pretty() + "}"
}

func (i IfThenElse) pretty() string {
	return "if" + i.cond.pretty() + "{" + i.thenBlock.pretty() + "} else {" + i.elseBlock.pretty() + "}"
}

func (p Print) pretty() string {
	return "print(" + p.exp.pretty() + ")"
}

func (d Decl) pretty() string {
	return d.name + ":=" + d.expression.pretty()
}

func (a Assign) pretty() string {
	return a.name + "=" + a.expression.pretty()
}

func (e ExpNode) pretty() string {
	switch e.op {
	case PLUS:
		return e.left.pretty() + "+" + e.right.pretty()
	case MULT:
		return e.left.pretty() + "*" + e.right.pretty()
	case AND:
		return e.left.pretty() + "&&" + e.right.pretty()
	case OR:
		return e.left.pretty() + "||" + e.right.pretty()
	case NOT:
		return "!" + e.left.pretty()
	case EQUAL:
		return e.left.pretty() + "==" + e.right.pretty()
	case LESS:
		return e.left.pretty() + "<" + e.right.pretty()
	}
	return ""
}

func (v Variable) pretty() string {
	return v.name
}

func (v Value) pretty() string {
	if v.Type == Integer {
		return string(rune(v.iValue))
	}
	if v.bValue == true {
		return "true"
	}
	return "false"
}
