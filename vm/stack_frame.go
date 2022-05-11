package vm

import "gscript/proto"

type stackFrame struct {
	prev        *stackFrame
	symbolTable *symbolTable
	wantRetCnt  int
	returnAddr  uint32
	upValues    []*interface{}
}

type Closure struct {
	Info     *proto.BasicInfo
	UpValues []*interface{}
}

func newFuncFrame() *stackFrame {
	return &stackFrame{
		symbolTable: newSymbolTable(),
	}
}
