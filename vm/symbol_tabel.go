package vm

type symbolTable struct {
	values []interface{}
}

func newSymbolTable() *symbolTable {
	return &symbolTable{}
}

func (st *symbolTable) getValue(idx uint32) interface{} {
	if idx >= uint32(len(st.values)) {
		panic("index out of symbol table")
	}
	return st.values[idx]
}

func (st *symbolTable) setValue(idx uint32, val interface{}) {
	if idx >= uint32(len(st.values)) {
		panic("index out of symbol table")
	}
	st.values[idx] = val
}

func (st *symbolTable) pushSymbol(val interface{}) {
	st.values = append(st.values, val)
}
