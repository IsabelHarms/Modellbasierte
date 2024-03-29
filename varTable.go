package main

// Types

type IMPtype int

const (
	Boolean   IMPtype = 1
	Integer   IMPtype = 2
	Undefined IMPtype = 3 // when parsing illegal expressions
)

type Value struct { // leaves of expression trees
	Type   IMPtype
	bValue bool
	iValue int
}

func showType(t IMPtype) string {
	var s string
	switch {
	case t == Integer:
		s = "Int"
	case t == Boolean:
		s = "Bool"
	case t == Undefined:
		s = "Illtyped"
	}
	return s
}

const MaxNesting = 50

type VarTable struct {
	nesting int
	names   [MaxNesting]map[string]*Value // array of maps, each is still = nil
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
		vt.names[vt.nesting] = map[string]*Value{}
	} // new map for current block
	vt.names[vt.nesting][name] = &Value{Type: Integer, iValue: i}
}

func (vt *VarTable) declareBool(name string, b bool) {
	if vt.names[vt.nesting] == nil {
		vt.names[vt.nesting] = map[string]*Value{}
	} // new map for current block
	vt.names[vt.nesting][name] = &Value{Type: Boolean, bValue: b}
}

func (vt *VarTable) declareUndefined(name string) {
	if vt.names[vt.nesting] == nil {
		vt.names[vt.nesting] = map[string]*Value{}
	} // new map for current block
	vt.names[vt.nesting][name] = &Value{Type: Undefined}
}

func (vt *VarTable) Get(name string) *Value {
	for i := vt.nesting; i >= 0; i-- {
		if vt.names[i] == nil {
			continue
		}
		v := vt.names[i][name]
		if v != nil {
			return v
		}
	}
	return nil
}
