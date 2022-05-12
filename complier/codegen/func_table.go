package codegen

import (
	"gscript/complier/ast"
	"gscript/proto"
)

type FuncTable struct {
	funcTable []proto.FuncProto // index -> proto
	funcMap   map[string]uint32 // funcname -> index

	anonymousFuncs []proto.AnonymousFuncProto
}

func newFuncTable(funcs []*ast.FuncDefStmt) *FuncTable {
	ft := &FuncTable{
		funcMap:   make(map[string]uint32),
		funcTable: make([]proto.FuncProto, len(funcs)),
	}
	for i, f := range funcs {
		ft.funcMap[f.Name] = uint32(i)
		info := new(proto.BasicInfo)
		info.Parameters = f.Parameters
		info.VaArgs = f.VaArgs != ""
		ft.funcTable[i].Info = info
	}
	return ft
}
