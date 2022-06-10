package vm

import (
	"gscript/vm/types"
)

type symbolTable struct {
	values []*types.GsValue
}

func newSymbolTable() *symbolTable {
	return &symbolTable{}
}

func (st *symbolTable) getValue(idx uint32) interface{} {
	if idx >= uint32(len(st.values)) {
		exit("index(%d) out of variables table length(%d)", idx, len(st.values))
	}
	return st.values[idx].Value
}

func (st *symbolTable) setValue(idx uint32, val interface{}) {
	if idx >= uint32(len(st.values)) {
		exit("index(%d) out of variables table length(%d)", idx, len(st.values))
	}
	st.values[idx].Value = val
}

func (st *symbolTable) pushSymbol(val interface{}) {
	st.values = append(st.values, &types.GsValue{Value: val})
}

func (st *symbolTable) top() (val interface{}) {
	return st.values[len(st.values)-1].Value
}

func (st *symbolTable) resizeTo(size int) {
	if size >= len(st.values) {
		return
	}
	for i := size; i < len(st.values); i++ {
		st.values[i] = nil
	}
	st.values = st.values[:size]
}
