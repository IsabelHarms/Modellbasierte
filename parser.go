package main

import "fmt"

func parse() {

	vartable := &VarTable{nesting: -1} // ready for first start of block
	// work with a pointer, otherwise you will get multiple structs!

	vartable.blockStart() // dereferencing is implied! :-(

	vartable.declareInt("Isabel", 22)
	vartable.declareInt("David", 25)
	vartable.blockStart()
	vartable.declareInt("Isabel", 2000)
	vartable.Get("Isabel").printValue()
	vartable.blockEnd()
	vartable.Get("Isabel").printValue()
	vartable.blockStart()
	vartable.Get("Isabel").printValue()
}

type IMPtype byte

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
}

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
		vt.names[vt.nesting] = map[string]*Variable{}
	} // new map for current block
	vt.names[vt.nesting][name] = &Variable{varType: INTEGER, iValue: i}
}

func (vt *VarTable) Get(name string) *Variable {
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
