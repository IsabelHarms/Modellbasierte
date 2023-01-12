package main

import "fmt"

func pretty(exp *Exp) {

}

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

// operator, variable, value switch
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
		return &Value{Type: Boolean, bValue: e.left.eval().iValue <= e.right.eval().iValue}
	case AND:
		return &Value{Type: Boolean, bValue: e.left.eval().bValue && e.right.eval().bValue}
	case OR:
		return &Value{Type: Boolean, bValue: e.left.eval().bValue || e.right.eval().bValue}
	case NOT:
		return &Value{Type: Boolean, bValue: !e.left.eval().bValue}
	}
	return nil //todo this should never happen
}

func (v Value) eval() *Value {
	return &v
}
