package vm

import "gscript/proto"

type stackFrame struct {
	prev        *stackFrame
	symbolTable *symbolTable
	wantRetCnt  int
	returnAddr  uint32
	upValues    []*GsValue
}

type Closure struct {
	Info     *proto.BasicInfo
	UpValues []*GsValue
}

func newFuncFrame() *stackFrame {
	return &stackFrame{
		symbolTable: newSymbolTable(),
	}
}
