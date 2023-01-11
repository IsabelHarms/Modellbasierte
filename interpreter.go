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

/*func (d Decl) exec() {
	for _, stmt := range b {
		stmt.exec()
	}
}
func (b Assign) exec() {
	for _, stmt := range b {
		stmt.exec()
	}
}*/

func (e ExpNode) eval() Value {
	return Value{}
}

func (v Value) eval() Value {
	return Value{}
}
