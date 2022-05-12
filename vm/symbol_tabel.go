package vm

type symbolTable struct {
	values []*GsValue
}

type GsValue struct {
	value interface{}
}

func newSymbolTable() *symbolTable {
	return &symbolTable{}
}

func (st *symbolTable) getValue(idx uint32) interface{} {
	if idx >= uint32(len(st.values)) {
		panic("index out of symbol table")
	}
	return st.values[idx].value
}

func (st *symbolTable) setValue(idx uint32, val interface{}) {
	if idx >= uint32(len(st.values)) {
		panic("index out of symbol table")
	}
	st.values[idx].value = val
}

func (st *symbolTable) pushSymbol(val interface{}) {
	st.values = append(st.values, &GsValue{val})
}

func (st *symbolTable) top() (val interface{}) {
	return st.values[len(st.values)-1].value
}
