package vm

import "gscript/vm/types"

type symbolTable struct {
	values []*types.GsValue
}

func newSymbolTable() *symbolTable {
	return &symbolTable{}
}

func (st *symbolTable) getValue(idx uint32) interface{} {
	if idx >= uint32(len(st.values)) {
		panic("index out of symbol table")
	}
	return st.values[idx].Value
}

func (st *symbolTable) setValue(idx uint32, val interface{}) {
	if idx >= uint32(len(st.values)) {
		panic("index out of symbol table")
	}
	st.values[idx].Value = val
}

func (st *symbolTable) pushSymbol(val interface{}) {
	st.values = append(st.values, &types.GsValue{val})
}

func (st *symbolTable) top() (val interface{}) {
	return st.values[len(st.values)-1].Value
}
