package vm

import (
	"gscript/proto"
	"gscript/vm/types"
)

type stackFrame struct {
	prev        *stackFrame
	symbolTable *symbolTable
	wantRetCnt  int
	pc          uint32
	upValues    []*types.GsValue
	text        []proto.Instruction
}

func newFuncFrame() *stackFrame {
	return &stackFrame{
		symbolTable: newSymbolTable(),
	}
}
