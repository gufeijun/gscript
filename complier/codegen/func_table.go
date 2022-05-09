package codegen

import (
	"gscript/complier/ast"
	"gscript/proto"
)

type FuncTable struct {
	funcTable []proto.Func
	funcMap   map[string]uint32 // funcname -> index
}

func newFuncTable(funcs []*ast.FuncDefStmt) *FuncTable {
	ft := &FuncTable{
		funcMap:   make(map[string]uint32),
		funcTable: make([]proto.Func, len(funcs)),
	}
	for i, f := range funcs {
		ft.funcMap[f.Name] = uint32(i)
		ft.funcTable[i].Parameters = f.Parameters
		ft.funcTable[i].VaArgs = f.VaArgs != ""
		i++
	}
	return ft
}
